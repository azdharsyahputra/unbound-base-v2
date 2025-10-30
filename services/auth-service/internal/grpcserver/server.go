package grpcserver

import (
	"context"
	"fmt"
	"log"

	"unbound-v2/services/auth-service/internal/auth"
	authpb "unbound-v2/shared/proto/auth"
)

type AuthGRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	authService *auth.AuthService
}

// Konstruktor
func NewAuthGRPCServer(authService *auth.AuthService) *AuthGRPCServer {
	return &AuthGRPCServer{authService: authService}
}

// ✅ Validasi JWT token via gRPC
func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	tokenStr := req.Token
	log.Printf("🔍 Validating token via gRPC...")

	userID, err := s.authService.ParseToken(tokenStr)
	if err != nil {
		log.Printf("❌ Invalid token: %v", err)
		return &authpb.ValidateTokenResponse{Valid: false}, nil
	}

	log.Printf("✅ Token valid untuk user_id=%d", userID)

	return &authpb.ValidateTokenResponse{
		Valid:  true,
		UserId: fmt.Sprintf("%d", userID),
	}, nil
}

// ✅ Dapatkan user berdasarkan username (untuk user-service follow)
func (s *AuthGRPCServer) GetUserByUsername(ctx context.Context, req *authpb.GetUserByUsernameRequest) (*authpb.GetUserByUsernameResponse, error) {
	var user auth.User
	if err := s.authService.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.Printf("❌ user %s not found: %v", req.Username, err)
		return nil, fmt.Errorf("user not found")
	}

	return &authpb.GetUserByUsernameResponse{
		Id:       uint64(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
