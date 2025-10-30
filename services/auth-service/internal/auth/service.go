package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func NewAuthService(db *gorm.DB) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	return &AuthService{DB: db, JWTSecret: []byte(secret)}
}

type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Register(input RegisterReq) (*User, error) {
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	var count int64
	s.DB.Model(&User{}).Where("email = ? OR username = ?", input.Email, input.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("email/username already used")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hash),
	}

	if err := s.DB.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(input LoginReq) (*TokenResp, error) {
	var u User
	if err := s.DB.Where("email = ?", input.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.GenerateJWT(u.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, err
	}

	return &TokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GenerateJWT(userID uint) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "unbound",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JWTSecret)
}

func (s *AuthService) GenerateRefreshToken(userID uint) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	expiry := time.Now().Add(7 * 24 * time.Hour)

	rt := RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiry,
	}
	if err := s.DB.Create(&rt).Error; err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) ParseToken(tokenStr string) (uint, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.JWTSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := tok.Claims.(*jwt.RegisteredClaims); ok && tok.Valid {
		id64, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(id64), nil
	}
	return 0, errors.New("invalid token")
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.DB.Where("token = ?", refreshToken).Delete(&RefreshToken{}).Error
}

func (s *AuthService) RefreshAccess(refreshToken string) (*TokenResp, error) {
	var rt RefreshToken
	if err := s.DB.Where("token = ?", refreshToken).First(&rt).Error; err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		s.DB.Delete(&rt)
		return nil, errors.New("refresh token expired")
	}

	newAccess, err := s.GenerateJWT(rt.UserID)
	if err != nil {
		return nil, err
	}

	return &TokenResp{
		AccessToken:  newAccess,
		RefreshToken: refreshToken,
	}, nil
}
