package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ihrk/rest-api-task/internal/app/apperrors"
	"github.com/ihrk/rest-api-task/internal/app/models"
	"github.com/ihrk/rest-api-task/internal/app/requests"
)

type SignerIn interface {
	SignIn(ctx context.Context, r *requests.SignIn) (*models.Token, error)
}

type SignInHandler struct {
	signerIn SignerIn
}

func NewSignInHandler(signerIn SignerIn) *SignInHandler {
	return &SignInHandler{signerIn}
}

func (h *SignInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := h.serveHTTP(r)
	if err != nil {
		RespondError(w, err)

		return
	}

	respBody := TokenToResponse(token)

	RespondJSON(w, http.StatusOK, &respBody)
}

func (h *SignInHandler) serveHTTP(r *http.Request) (*models.Token, error) {
	var req requests.SignIn

	err := ParseRequestJSON(r, &req)
	if err != nil {
		return nil, err
	}

	return h.signerIn.SignIn(r.Context(), &req)
}

type TokenRefresher interface {
	RefreshToken(ctx context.Context, r *requests.RefreshToken) (*models.Token, error)
}

type RefreshTokenHandler struct {
	tokenRefresher TokenRefresher
}

func NewRefreshTokenHandler(tokenRefresher TokenRefresher) *RefreshTokenHandler {
	return &RefreshTokenHandler{tokenRefresher}
}

func (h *RefreshTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := h.serveHTTP(r)
	if err != nil {
		RespondError(w, err)

		return
	}

	respBody := TokenToResponse(token)

	RespondJSON(w, http.StatusOK, &respBody)
}

func (h *RefreshTokenHandler) serveHTTP(r *http.Request) (*models.Token, error) {
	var req requests.RefreshToken

	err := ParseRequestJSON(r, &req)
	if err != nil {
		return nil, err
	}

	return h.tokenRefresher.RefreshToken(r.Context(), &req)
}

type AccessTokenValidator interface {
	ValidateAccessToken(ctx context.Context, token string) error
}

type JWTAuthMiddleware struct {
	accessTokenValidator AccessTokenValidator
}

func NewJWTAuthMiddleware(accessTokenValidator AccessTokenValidator) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{accessTokenValidator}
}

func (m *JWTAuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")

		if !strings.HasPrefix(bearerToken, "Bearer ") {
			RespondError(w, fmt.Errorf("malformed token: %w", apperrors.ErrUnauthorized))

			return
		}

		err := m.accessTokenValidator.ValidateAccessToken(
			r.Context(),
			strings.TrimPrefix(bearerToken, "Bearer "),
		)
		if err != nil {
			RespondError(w, fmt.Errorf("invalid token: %w", err))

			return
		}

		next.ServeHTTP(w, r)
	})
}
