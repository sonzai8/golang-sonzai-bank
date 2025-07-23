package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"github.com/sonzai8/golang-sonzai-bank/token"
	"github.com/sonzai8/golang-sonzai-bank/utils"
)

type Server struct {
	config     utils.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Fatal().Msg("register validation error:")
		}
	}
	server.SetupRouter()
	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()
	router.POST("/users", server.CreateUser)
	router.POST("users/login", server.LoginUser)
	router.POST("tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/accounts", server.ListAccount)
	authRoutes.GET("/accounts/:id", server.GetAccount)
	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.POST("/accounts:id", server.UpdateAccount)
	authRoutes.DELETE("/accounts/:id", server.DeleteAccount)

	authRoutes.POST("/transfers", server.Transfer)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(":" + address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
