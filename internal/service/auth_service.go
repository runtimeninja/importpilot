package service

import (
	"context"
	"errors"
	"strings"

	"github.com/runtimeninja/importpilot/internal/repository"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *JWTService
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token string
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)

	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	roles, err := s.userRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, roles)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
	}, nil
}
