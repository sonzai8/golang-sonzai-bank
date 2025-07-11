package gapi

import (
	"fmt"
	"github.com/sonzai8/golang-sonzai-bank/pb"

	db "github.com/sonzai8/golang-sonzai-bank/db/sqlc"
	"github.com/sonzai8/golang-sonzai-bank/token"
	"github.com/sonzai8/golang-sonzai-bank/utils"
)

type Server struct {
	pb.UnimplementedSonZaiBankServer
	config     utils.Config
	tokenMaker token.Maker
	store      db.Store
}

// NewServer creates a new gRPC server
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

	return server, nil
}
