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
			name: "CreatePost_Success",
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
			name: "CreatePost_ErrorEmptyBody",
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
			name: "CreatePost_InternalError",
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
			name: "OK",
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
			name: "Error_NotFoundToken",
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
			name: "Error_InternalError",
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
			name: "OK",
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
			name: "Error_NotFoundToken",
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
			name: "Error_InternalError",
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

func TestUpdatePostsById(t *testing.T) {
	randomPost := generateRandomPost()
	postResponse := mapToPostResponse(randomPost)

	// + TODO: update both title and content
	// + TODO: update title
	// + TODO: update content
	// + TODO: request when no title and content
	// + TODO: not found
	// TODO: internal error when looking for post
	// TODO: internal error when updating
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "UpdatePost_Success",
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
			name: "UpdatePost_UpdateOnlyTitleSuccess",
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
			name: "UpdatePost_UpdateOnlyContentSuccess",
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
			name: "UpdatePost_UpdateOnlyUpdateAtAttributeSuccess",
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
			name: "UpdatePost_PostNotFoundError",
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
			name: "Error_FetchPostByIdInternalError",
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
			name: "Error_UpdatePostByIdInternalError",
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
			name: "OK",
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
			name: "Error_NotFoundToken",
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
			name: "Error_InternalError",
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
