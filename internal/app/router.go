package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ihrk/rest-api-task/internal/app/handlers"
	"github.com/ihrk/rest-api-task/internal/app/logic"
)

func MakeRouter(
	authSvc *logic.AuthService,
	companySvc *logic.CompanyService,
) http.Handler {
	jwtAuthMiddleware := handlers.NewJWTAuthMiddleware(authSvc)

	router := mux.NewRouter()
	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	apiV1Router.Handle("/signin",
		handlers.NewSignInHandler(authSvc)).Methods(http.MethodPost)

	apiV1Router.Handle("/refreshtoken",
		handlers.NewRefreshTokenHandler(authSvc)).Methods(http.MethodPost)

	apiV1Router.Handle("/companies",
		jwtAuthMiddleware.Wrap(
			handlers.NewCreateCompanyHandler(companySvc))).Methods(http.MethodPost)

	apiV1Router.Handle("/companies/{id:[0-9]+}",
		handlers.NewGetCompanyHandler(companySvc)).Methods(http.MethodGet)

	apiV1Router.Handle("/companies/{id:[0-9]+}",
		jwtAuthMiddleware.Wrap(
			handlers.NewDeleteCompanyHandler(companySvc))).Methods(http.MethodDelete)

	apiV1Router.Handle("/companies/{id:[0-9]+}",
		handlers.NewUpdateCompanyHandler(companySvc)).Methods(http.MethodPut)

	apiV1Router.Handle("/companies",
		handlers.NewListCompaniesHandler(companySvc)).Methods(http.MethodGet)

	return router
}
