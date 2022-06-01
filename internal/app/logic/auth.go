package logic

import (
	"context"

	"github.com/ihrk/rest-api-task/internal/app/apperrors"
	"github.com/ihrk/rest-api-task/internal/app/logic/auth"
	"github.com/ihrk/rest-api-task/internal/app/models"
	"github.com/ihrk/rest-api-task/internal/app/requests"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) SignIn(
	ctx context.Context,
	r *requests.SignIn,
) (*models.Token, error) {
	user := models.User{
		Login:    r.Login,
		Password: r.Password,
	}

	err := s.checkCreds(&user)
	if err != nil {
		return nil, err
	}

	return auth.CreateToken(&user)
}

func (s *AuthService) checkCreds(user *models.User) error {
	if user.Login != "admin" || user.Password != "qwerty" {
		return apperrors.ErrUnauthorized
	}

	return nil
}

func (s *AuthService) ValidateAccessToken(
	ctx context.Context,
	accessToken string,
) error {
	user, err := auth.ValidateAccessToken(accessToken)
	if err != nil {
		return err
	}

	return s.checkCreds(user)
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	r *requests.RefreshToken,
) (*models.Token, error) {
	user, err := auth.ValidateRefreshToken(r.RefreshToken)
	if err != nil {
		return nil, err
	}

	err = s.checkCreds(user)
	if err != nil {
		return nil, err
	}

	return auth.CreateToken(user)
}
