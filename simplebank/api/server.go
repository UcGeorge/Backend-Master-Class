package api

import (
	db "github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Initialize router
	server.router = router

	// Setup account routes
	server.setupAccountRoutes()

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorsResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
