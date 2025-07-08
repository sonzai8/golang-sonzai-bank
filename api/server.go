package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"log"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	router := gin.Default()

	server := &Server{
		store:  store,
		router: router,
	}

	router.POST("/users", server.CreateUser)

	router.GET("/accounts", server.ListAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.POST("/accounts", server.CreateAccount)
	router.POST("/accounts:id", server.UpdateAccount)
	router.DELETE("/accounts/:id", server.DeleteAccount)

	router.POST("/transfer", server.Transfer)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Fatalf("register validation error: %v", err)
		}
	}
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(":" + address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
