package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedFor              = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	metaData := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("%+v \n", md)
		if clientIp := md.Get(xForwardedFor); len(clientIp) > 0 {
			metaData.ClientIP = clientIp[0]
		}

		if agent := md.Get(userAgentHeader); len(agent) > 0 {
			metaData.UserAgent = agent[0]
		}

		if agentHeader := md.Get(grpcGatewayUserAgentHeader); len(agentHeader) > 0 {
			metaData.UserAgent = agentHeader[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		{
			metaData.ClientIP = p.Addr.String()
		}
	}
	return metaData
}
