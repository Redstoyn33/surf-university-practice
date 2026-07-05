package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glini/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	clientRepo domain.ClientRepository
	authSecret string
}

func NewAuthService(clientRepo domain.ClientRepository, authSecret string) *AuthService {
	return &AuthService{clientRepo: clientRepo, authSecret: authSecret}
}

type RegisterResult struct {
	Client domain.Client
	Token  string
}

func (s *AuthService) Register(ctx context.Context, login, password string) (RegisterResult, error) {
	if login == "" {
		return RegisterResult{}, fmt.Errorf("login: %w", domain.ErrValidation)
	}
	if len(password) < 6 {
		return RegisterResult{}, fmt.Errorf("password must be at least 6 characters: %w", domain.ErrValidation)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResult{}, fmt.Errorf("hash password: %w", err)
	}

	client, err := s.clientRepo.InsertClient(ctx, login, string(hash))
	if err != nil {
		if errors.Is(err, domain.ErrDuplicate) {
			return RegisterResult{}, fmt.Errorf("login already taken: %w", domain.ErrConflict)
		}
		return RegisterResult{}, fmt.Errorf("insert client: %w", err)
	}

	token, err := s.generateToken(client.ID)
	if err != nil {
		return RegisterResult{}, fmt.Errorf("generate token: %w", err)
	}

	return RegisterResult{Client: client, Token: token}, nil
}

type LoginResult struct {
	Client domain.Client
	Token  string
}

func (s *AuthService) Login(ctx context.Context, login, password string) (LoginResult, error) {
	client, err := s.clientRepo.GetClientByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return LoginResult{}, fmt.Errorf("invalid credentials: %w", domain.ErrForbidden)
		}
		return LoginResult{}, fmt.Errorf("get client: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.PasswordHash), []byte(password)); err != nil {
		return LoginResult{}, fmt.Errorf("invalid credentials: %w", domain.ErrForbidden)
	}

	token, err := s.generateToken(client.ID)
	if err != nil {
		return LoginResult{}, fmt.Errorf("generate token: %w", err)
	}

	return LoginResult{Client: client, Token: token}, nil
}

func (s *AuthService) generateToken(clientID int64) (string, error) {
	claims := jwt.MapClaims{
		"client_id": clientID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.authSecret))
}
