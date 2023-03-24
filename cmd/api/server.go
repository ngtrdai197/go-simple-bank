package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/ngtrdai197/simple-bank/db/sqlc"
	"github.com/ngtrdai197/simple-bank/token"
	"github.com/ngtrdai197/simple-bank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer create a new HTTP server, and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter(config)
	return server, nil
}

func (server *Server) setupRouter(config util.Config) {
	router := gin.Default()

	router.Use(util.DefaultStructuredLogger())
	if config.AppEnv != "DEVELOP" {
		gin.SetMode(gin.ReleaseMode)
	}

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/")
	authRoutes.Use(authMiddleware(server.tokenMaker))
	{
		authRoutes.POST("/accounts", server.createAccount)
		authRoutes.GET("/accounts/:id", server.getAccount)
		authRoutes.GET("/accounts", server.listAccounts)

		authRoutes.POST("/transfers", server.createTransfer)
	}
	server.router = router
}

// Start a new HTTP server, and specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
