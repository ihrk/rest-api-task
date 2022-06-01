package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ihrk/rest-api-task/internal/app/handlers"
	"github.com/ihrk/rest-api-task/internal/app/logic/auth"
	"github.com/ihrk/rest-api-task/internal/app/requests"
)

type AuthTestSuite struct {
	HandlersTestSuite
}

func TestAuthHandlers(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) TestAuth() {
	token := s.signIn()

	s.Run("refresh token", func() {
		resp, err := doRequest(
			http.MethodPost,
			"/refreshtoken",
			&requests.RefreshToken{
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
			},
			"",
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		var result handlers.TokenResponse

		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)

		u, err := auth.ValidateAccessToken(result.AccessToken)
		s.Require().NoError(err)
		s.Require().Equal("admin", u.Login)
		s.Require().Equal("qwerty", u.Password)

		u, err = auth.ValidateRefreshToken(result.RefreshToken)
		s.Require().NoError(err)
		s.Require().Equal("admin", u.Login)
		s.Require().Equal("qwerty", u.Password)

		s.Require().NoError(resp.Body.Close())
	})
}
