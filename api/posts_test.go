package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "promova-test-task/db/mock"
	db "promova-test-task/db/sqlc"
	"promova-test-task/util"
	"testing"
	"time"
)

func TestCreatePosts(t *testing.T) {
	randomPost := generateRandomPost()
	postResponse := mapToPostResponse(randomPost)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "positive_CreatePost",
			body: gin.H{
				"title":   randomPost.Title,
				"content": randomPost.Content,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := db.CreatePostParams{
					Title:   randomPost.Title,
					Content: randomPost.Content,
				}
				querier.EXPECT().
					CreatePost(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "negative_CreatePost_EmptyRequestBody",
			body: gin.H{},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatchErrResponse(t, recorder.Body, ErrResponse{Error: "Key: 'createPostRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag\nKey: 'createPostRequest.Content' Error:Field validation for 'Content' failed on the 'required' tag"})
			},
		},
		{
			name: "negative_CreatePost_InternalError",
			body: gin.H{
				"title":   randomPost.Title,
				"content": randomPost.Content,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := db.CreatePostParams{
					Title:   randomPost.Title,
					Content: randomPost.Content,
				}
				querier.EXPECT().
					CreatePost(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Post{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockQuerier(controller)
			testCase.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			path := "/posts"
			requestUrl := fmt.Sprintf("%s", path)
			request, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestGetPostById(t *testing.T) {
	randomPost := generateRandomPost()
	postResponse := mapToPostResponse(randomPost)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "positive_GetPostById",
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "negative_GetPostById_PostNotFound",
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Post{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				requireBodyMatchErrResponse(t, recorder.Body, ErrResponse{Error: "sql: no rows in result set"})
			},
		},
		{
			name: "negative_GetPostById_InternalError",
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Post{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockQuerier(controller)
			testCase.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			path := "/posts"
			requestUrl := fmt.Sprintf("%s/%d", path, randomPost.ID)
			request, err := http.NewRequest(http.MethodGet, requestUrl, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestGetPosts(t *testing.T) {
	randomPost := generateRandomPost()
	postsResponse := mapToPostsResponse([]db.Post{randomPost})

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "positive_GetPosts",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetPosts(gomock.Any()).
					Times(1).
					Return([]db.Post{randomPost}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPosts(t, recorder.Body, postsResponse)
			},
		},
		{
			name: "negative_GetPosts_PostsNotFound",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetPosts(gomock.Any()).
					Times(1).
					Return([]db.Post{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				requireBodyMatchErrResponse(t, recorder.Body, ErrResponse{Error: "sql: no rows in result set"})
			},
		},
		{
			name: "negative_GetPosts_InternalError",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetPosts(gomock.Any()).
					Times(1).
					Return([]db.Post{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockQuerier(controller)
			testCase.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			path := "/posts"
			requestUrl := fmt.Sprintf("%s", path)
			request, err := http.NewRequest(http.MethodGet, requestUrl, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestUpdatePostById(t *testing.T) {
	randomPost := generateRandomPost()
	postResponse := mapToPostResponse(randomPost)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "positive_UpdatePost",
			body: gin.H{
				"title":   randomPost.Title,
				"content": randomPost.Content,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)

				updateArg := db.UpdatePostByIdParams{
					ID:      randomPost.ID,
					Title:   randomPost.Title,
					Content: randomPost.Content,
				}

				updatedPost := db.Post{
					ID:        updateArg.ID,
					Title:     updateArg.Title,
					Content:   updateArg.Content,
					CreatedAt: randomPost.CreatedAt,
					UpdatedAt: time.Now(),
				}

				postResponse.Title = updatedPost.Title
				postResponse.Content = updatedPost.Content
				postResponse.UpdatedAt = updatedPost.UpdatedAt.Format("2006-01-02 15:04:05")

				querier.EXPECT().
					UpdatePostById(gomock.Any(), gomock.Eq(updateArg)).
					Times(1).
					Return(updatedPost, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "positive_UpdatePost_UpdateOnlyTitle",
			body: gin.H{
				"title": randomPost.Title,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)

				updateArg := db.UpdatePostByIdParams{
					ID:    randomPost.ID,
					Title: randomPost.Title,
				}

				updatedPost := db.Post{
					ID:        updateArg.ID,
					Title:     updateArg.Title,
					Content:   randomPost.Content,
					CreatedAt: randomPost.CreatedAt,
					UpdatedAt: time.Now(),
				}

				postResponse.Title = updatedPost.Title
				postResponse.UpdatedAt = updatedPost.UpdatedAt.Format("2006-01-02 15:04:05")

				querier.EXPECT().
					UpdatePostById(gomock.Any(), gomock.Eq(updateArg)).
					Times(1).
					Return(updatedPost, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "positive_UpdatePost_UpdateOnlyContent",
			body: gin.H{
				"content": randomPost.Content,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)

				updateArg := db.UpdatePostByIdParams{
					ID:      randomPost.ID,
					Content: randomPost.Content,
				}

				updatedPost := db.Post{
					ID:        updateArg.ID,
					Title:     randomPost.Title,
					Content:   updateArg.Content,
					CreatedAt: randomPost.CreatedAt,
					UpdatedAt: time.Now(),
				}

				postResponse.Content = updatedPost.Content
				postResponse.UpdatedAt = updatedPost.UpdatedAt.Format("2006-01-02 15:04:05")

				querier.EXPECT().
					UpdatePostById(gomock.Any(), gomock.Eq(updateArg)).
					Times(1).
					Return(updatedPost, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "positive_UpdatePost_UpdateOnlyUpdatedAtAttribute",
			body: gin.H{},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)

				updateArg := db.UpdatePostByIdParams{
					ID: randomPost.ID,
				}

				updatedPost := db.Post{
					ID:        updateArg.ID,
					Title:     randomPost.Title,
					Content:   randomPost.Content,
					CreatedAt: randomPost.CreatedAt,
					UpdatedAt: time.Now(),
				}

				postResponse.UpdatedAt = updatedPost.UpdatedAt.Format("2006-01-02 15:04:05")

				querier.EXPECT().
					UpdatePostById(gomock.Any(), gomock.Eq(updateArg)).
					Times(1).
					Return(updatedPost, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, postResponse)
			},
		},
		{
			name: "negative_UpdatePost_PostNotFound",
			body: gin.H{
				"title":   randomPost.Title,
				"content": randomPost.Content,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Post{}, sql.ErrNoRows)

				querier.EXPECT().
					UpdatePostById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				requireBodyMatchErrResponse(t, recorder.Body, ErrResponse{Error: "sql: no rows in result set"})
			},
		},
		{
			name: "negative_UpdatePostById_FetchPostByIdInternalError",
			body: gin.H{
				"title":   randomPost.Title,
				"content": randomPost.Content,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Post{}, sql.ErrConnDone)

				querier.EXPECT().
					UpdatePostById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "negative_UpdatePostById_UpdatePostByIdInternalError",
			body: gin.H{
				"title":   randomPost.Title,
				"content": randomPost.Content,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)

				querier.EXPECT().
					UpdatePostById(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Post{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockQuerier(controller)
			testCase.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			path := "/posts"
			requestUrl := fmt.Sprintf("%s/%d", path, randomPost.ID)
			request, err := http.NewRequest(http.MethodPut, requestUrl, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func TestDeletePostById(t *testing.T) {
	randomPost := generateRandomPost()

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "positive_DeletePost",
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)

				querier.EXPECT().
					DeletePost(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "negative_DeletePost_PostNotFound",
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Post{}, sql.ErrNoRows)

				querier.EXPECT().
					DeletePost(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				requireBodyMatchErrResponse(t, recorder.Body, ErrResponse{Error: "sql: no rows in result set"})
			},
		},
		{
			name: "negative_DeletePost_InternalError",
			buildStubs: func(querier *mockdb.MockQuerier) {
				arg := randomPost.ID

				querier.EXPECT().
					GetPostById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(randomPost, nil)

				querier.EXPECT().
					DeletePost(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockQuerier(controller)
			testCase.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			path := "/posts"
			requestUrl := fmt.Sprintf("%s/%d", path, randomPost.ID)
			request, err := http.NewRequest(http.MethodDelete, requestUrl, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

type PostCollection struct {
	Posts []PostResponse
}

func requireBodyMatchPosts(t *testing.T, body *bytes.Buffer, expected []PostResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var actual PostCollection
	err = json.Unmarshal(data, &actual.Posts)
	for index, actual := range actual.Posts {
		require.NoError(t, err)
		require.Equal(t, expected[index], actual)
	}
}

func requireBodyMatchPost(t *testing.T, body *bytes.Buffer, expected PostResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var actual PostResponse
	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func requireBodyMatchErrResponse(t *testing.T, body *bytes.Buffer, expected ErrResponse) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var actual ErrResponse
	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func generateRandomPost() db.Post {
	randomPost := util.GenerateRandomPost()
	return db.Post{
		ID:        randomPost.ID,
		Title:     randomPost.Title,
		Content:   randomPost.Content,
		CreatedAt: randomPost.CreatedAt,
		UpdatedAt: randomPost.UpdatedAt,
	}
}
