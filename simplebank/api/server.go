package api

import (
	db "github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	// Register validations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// Setup routes
	server.setupAccountRoutes()
	server.setupTransferRoutes()
	server.setupUserRoutes()

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorsResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
