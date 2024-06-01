package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	db "promova-test-task/db/sqlc"
	_ "promova-test-task/docs"
)

type ErrResponse struct {
	Error string `json:"error"`
}

// Server serves HTTP requests for attendance service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("/posts", server.getPosts)
	router.GET("/posts/:id", server.getPost)
	router.POST("/posts", server.createPost)
	router.PUT("/posts/:id", server.updatePost)
	router.DELETE("/posts/:id", server.deletePost)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) ErrResponse {
	return ErrResponse{Error: err.Error()}
}
