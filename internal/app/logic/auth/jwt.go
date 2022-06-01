package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ihrk/rest-api-task/internal/app/models"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Login    string
	Password string
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Token string
}

func CreateToken(user *models.User) (*models.Token, error) {
	var (
		err   error
		token models.Token
	)

	token.AccessToken, err = getTokenSignedString(
		&UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
			Login:    user.Login,
			Password: user.Password,
		})
	if err != nil {
		return nil, err
	}

	token.RefreshToken, err = getTokenSignedString(
		&TokenClaims{
			Token: token.AccessToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			},
		})
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func getTokenSignedString(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func ValidateAccessToken(accessToken string) (*models.User, error) {
	var claims UserClaims

	err := parseToken(accessToken, &claims)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Login:    claims.Login,
		Password: claims.Password,
	}, nil
}

func ValidateRefreshToken(refreshToken string) (*models.User, error) {
	var claims TokenClaims

	err := parseToken(refreshToken, &claims)
	if err != nil {
		return nil, err
	}

	return ValidateAccessToken(claims.Token)
}

func parseToken(tokenString string, claims jwt.Claims) error {
	p := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	_, err := p.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	return err
}
