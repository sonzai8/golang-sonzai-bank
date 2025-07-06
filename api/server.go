package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	router := gin.Default()

	server := &Server{
		store:  store,
		router: router,
	}
	router.GET("/accounts", server.ListAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.POST("/accounts", server.CreateAccount)
	router.POST("/accounts:id", server.UpdateAccount)
	router.DELETE("/accounts/:id", server.DeleteAccount)

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(":" + address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
