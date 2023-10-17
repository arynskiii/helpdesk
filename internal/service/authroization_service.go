package service

import (
	"context"
	"fmt"

	"github.com/arynskiii/help_desk/internal/repository"
	"github.com/arynskiii/help_desk/models"
	"github.com/gofrs/uuid"
)

var emptyStruct models.User

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user models.User) (int64, error) {

	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateToken(name string, password string) (models.User, error) {

	user, err := s.repo.GetUser(name)

	if err != nil {

		return emptyStruct, err
	}

	token, err := uuid.NewV4()
	if err != nil {
		return emptyStruct, fmt.Errorf("failed to generate token: ", err)

	}
	user.Token = token.String()
	if err := s.repo.SaveTokens(name, user.Token); err != nil {
		return emptyStruct, err
	}
	return user, nil
}

func (s *AuthService) GetUserByToken(token string) (models.User, error) {
	return s.repo.GetUserByToken(token)
}
