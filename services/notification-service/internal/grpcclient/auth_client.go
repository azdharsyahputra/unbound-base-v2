package grpcclient

import (
	"context"
	"log"

	"google.golang.org/grpc"

	authpb "unbound-v2/shared/proto/auth"
)

type AuthClient struct {
	client authpb.AuthServiceClient
}

func NewAuthClient(addr string) *AuthClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Failed to connect to auth gRPC: %v", err)
	}

	return &AuthClient{
		client: authpb.NewAuthServiceClient(conn),
	}
}

func (a *AuthClient) GetUserByID(id uint64) (*authpb.GetUserByIDResponse, error) {
	req := &authpb.GetUserByIDRequest{Id: id}
	return a.client.GetUserByID(context.Background(), req)
}
