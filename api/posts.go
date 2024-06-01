package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	db "promova-test-task/db/sqlc"
)

type createPostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type PostResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type getPostRequest struct {
	ID int `uri:"id" binding:"required"`
}

type updatePostRequestBody struct {
	Title   string `json:"title" binding:"omitempty"`
	Content string `json:"content" binding:"omitempty"`
}

// @Summary Create a post
// @Tags Post
// @Description Create a post
// @ID create-post
// @Accept json
// @Produce json
// @Param input body createPostRequest true "post entity related data"
// @Success 200 {object} PostResponse
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /posts [post]
func (s *Server) createPost(context *gin.Context) {
	var request createPostRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: Validate incoming data
	post, err := s.store.CreatePost(context, db.CreatePostParams{
		Title:   request.Title,
		Content: request.Content,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	postResponse := mapToPostResponse(post)
	context.JSON(http.StatusOK, postResponse)
}

// @Summary Get posts
// @Tags Post
// @Description Get all posts
// @ID get-posts
// @Produce json
// @Success 200 {object} []PostResponse
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /posts [get]
func (s *Server) getPosts(context *gin.Context) {
	posts, err := s.store.GetPosts(context)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	postsResponse := mapToPostsResponse(posts)
	context.JSON(http.StatusOK, postsResponse)
}

// @Summary Get post by id
// @Tags Post
// @Description Get a specific post by the specified id
// @ID get-post-by-id
// @Accept json
// @Produce json
// @Param id path string true "the specific post id"
// @Success 200 {object} PostResponse
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /posts/{id} [get]
func (s *Server) getPost(context *gin.Context) {
	var request getPostRequest

	if err := context.ShouldBindUri(&request); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	post, err := s.store.GetPostById(context, int32(request.ID))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	postResponse := mapToPostResponse(post)
	context.JSON(http.StatusOK, postResponse)
}

// @Summary Update post by id
// @Tags Post
// @Description Update a specific post by the specified id
// @ID update-post-by-id
// @Accept json
// @Produce json
// @Param id path string true "the specific post id"
// @Param input body updatePostRequestBody true "post entity related data"
// @Success 200
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /posts/{id} [put]
func (s *Server) updatePost(context *gin.Context) {
	var request getPostRequest
	var requestBody updatePostRequestBody

	if err := context.ShouldBindUri(&request); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := s.store.GetPostById(context, int32(request.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdatePostByIdParams{
		ID: post.ID,
	}

	if len(requestBody.Title) > 0 {
		arg.Title = requestBody.Title
	}
	if len(requestBody.Content) > 0 {
		arg.Content = requestBody.Content
	}

	post, err = s.store.UpdatePostById(context, arg)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			switch pqError.Code.Name() {
			case "modifying_sql_data_not_permitted":
				context.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	postResponse := mapToPostResponse(post)
	context.JSON(http.StatusOK, postResponse)
}

// @Summary Delete post by id
// @Tags Post
// @Description Delete a specific post by the specified id
// @ID delete-post-by-id
// @Accept json
// @Produce json
// @Param id path string true "the specific post id"
// @Success 200
// @Failure 400 {object} ErrResponse
// @Failure 500 {object} ErrResponse
// @Router /posts/{id} [delete]
func (s *Server) deletePost(context *gin.Context) {
	var request getPostRequest

	if err := context.ShouldBindUri(&request); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	post, err := s.store.GetPostById(context, int32(request.ID))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.store.DeletePost(context, post.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	context.Status(http.StatusOK)
}

func mapToPostsResponse(posts []db.Post) []PostResponse {
	responsePosts := make([]PostResponse, 0, len(posts))

	for _, post := range posts {
		responsePosts = append(responsePosts, mapToPostResponse(post))
	}

	return responsePosts
}

func mapToPostResponse(post db.Post) PostResponse {

	postResponse := PostResponse{
		ID:        int(post.ID),
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return postResponse
}
