package api

import (
	db "github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/db/sqlc"
	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/middleware"
	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/token"
	"github.com/UcGeorge/Upskill/BackendMasterClass/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	router     *gin.Engine
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: maker,
	}
	router := gin.Default()

	// Initialize router
	server.router = router

	// Register validations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// Setup routes
	server.setupUserRoutes(server.router)

	authRoutes := server.router.Group("/").Use(middleware.AuthMiddleware(server.tokenMaker))

	server.setupAccountRoutes(authRoutes)
	server.setupTransferRoutes(authRoutes)

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorsResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
