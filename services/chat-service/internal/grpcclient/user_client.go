package grpcclient

import (
	"context"
	"fmt"

	userpb "unbound-v2/shared/proto/user"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(c userpb.UserServiceClient) *UserClient {
	return &UserClient{client: c}
}

// VerifyUserExists memanggil user-service gRPC
func (u *UserClient) VerifyUserExists(userID uint) (bool, error) {
	resp, err := u.client.VerifyUserExists(context.Background(), &userpb.VerifyUserExistsRequest{
		UserId: uint64(userID),
	})
	if err != nil {
		return false, fmt.Errorf("user grpc error: %w", err)
	}
	return resp.Exists, nil
}
