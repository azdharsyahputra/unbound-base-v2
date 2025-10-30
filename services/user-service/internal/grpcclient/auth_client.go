package grpcclient

import (
	"context"
	"log"
	"time"

	authpb "unbound-v2/shared/proto/auth"

	"google.golang.org/grpc"
)

type AuthClient struct {
	client authpb.AuthServiceClient
}

// Konstruktor (buat koneksi ke Auth-service gRPC)
func NewAuthClient(addr string) *AuthClient {
	log.Printf("🔌 Connecting to Auth gRPC at %s...", addr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("❌ Failed to connect to auth-service via gRPC: %v", err)
	}

	log.Println("✅ Connected to Auth gRPC")
	c := authpb.NewAuthServiceClient(conn)
	return &AuthClient{client: c}
}

// ✅ Validasi token JWT via gRPC
func (a *AuthClient) ValidateToken(token string) (*authpb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := a.client.ValidateToken(ctx, &authpb.ValidateTokenRequest{Token: token})
	if err != nil {
		log.Printf("❌ ValidateToken gRPC error: %v", err)
		return nil, err
	}
	return res, nil
}

// ✅ Ambil user by username via gRPC (dipakai di follow_handler)
func (a *AuthClient) GetUserByUsername(username string) (*authpb.GetUserByUsernameResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := a.client.GetUserByUsername(ctx, &authpb.GetUserByUsernameRequest{Username: username})
	if err != nil {
		log.Printf("❌ GetUserByUsername error: %v", err)
		return nil, err
	}
	return res, nil
}
