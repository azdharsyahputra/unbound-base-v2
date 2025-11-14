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

// =======================
// VALIDATE TOKEN
// =======================
func (a *AuthClient) ValidateToken(token string) (uint, error) {
	resp, err := a.client.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return 0, fmt.Errorf("auth grpc error: %w", err)
	}

	uid64, err := strconv.ParseUint(resp.UserId, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user id returned from auth service: %w", err)
	}

	return uint(uid64), nil
}

// =======================
// VERIFY USER EXISTS
// =======================
func (a *AuthClient) VerifyUserExists(userID uint) (bool, error) {

	req := &authpb.GetUserByIDRequest{
		Id: uint64(userID),
	}

	resp, err := a.client.GetUserByID(context.Background(), req)
	if err != nil {
		// auth-service return error kalau user tidak ada
		return false, nil
	}

	// user ada kalau ID != 0
	return resp.Id != 0, nil
}
