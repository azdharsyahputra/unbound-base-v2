package grpcclient

import (
	"context"
	"fmt"
	"strconv"

	authpb "unbound-v2/shared/proto/auth"
)

type AuthClient struct {
	client authpb.AuthServiceClient
}

func NewAuthClient(c authpb.AuthServiceClient) *AuthClient {
	return &AuthClient{client: c}
}

func (a *AuthClient) ValidateToken(token string) (uint, error) {
	resp, err := a.client.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return 0, fmt.Errorf("auth grpc error: %w", err)
	}

	// Convert string → uint64
	uid64, err := strconv.ParseUint(resp.UserId, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user id returned from auth service: %w", err)
	}

	// Then convert uint64 → uint
	return uint(uid64), nil
}
