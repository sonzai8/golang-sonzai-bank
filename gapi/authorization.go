package gapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/sonzai8/golang-sonzai-bank/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	autorizationHeader = "authorization"
)

func (server *Server) AuthorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}
	values := md.Get(autorizationHeader)
	if len(values) == 0 {
		return nil, errors.New("missing authorization header")
	}
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return nil, errors.New("invalid authorization header")
	}
	authType := strings.ToLower(fields[0])
	if authType != "bearer" {
		return nil, errors.New("unsupported authorization type")
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	return payload, nil
}
