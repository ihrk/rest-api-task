package handlers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/suite"

	"github.com/ihrk/rest-api-task/internal/app"
	"github.com/ihrk/rest-api-task/internal/app/database"
	"github.com/ihrk/rest-api-task/internal/app/database/migrations"
	"github.com/ihrk/rest-api-task/internal/app/handlers"
	"github.com/ihrk/rest-api-task/internal/app/logic"
	"github.com/ihrk/rest-api-task/internal/app/models"
	"github.com/ihrk/rest-api-task/internal/app/requests"
)

const (
	hosturl    = "localhost"
	port       = ":8080"
	pathPrefix = "api/v1"
)

type HandlersTestSuite struct {
	suite.Suite
	m      *migrate.Migrate
	srv    *http.Server
	dbConn *sql.DB
}

func (s *HandlersTestSuite) SetupSuite() {
	err := s.Prepare("../../../configs/config.yml")
	s.Require().NoError(err)
}

func (s *HandlersTestSuite) TearDownSuite() {
	err := s.m.Down()
	s.Require().NoError(err)

	err = s.srv.Shutdown(context.Background())
	s.Require().NoError(err)

	err = s.dbConn.Close()
	s.Require().NoError(err)
}

func (s *HandlersTestSuite) Prepare(cfgPath string) error {
	cfg := app.LoadConfig(cfgPath)

	databaseURL := cfg.DB.URL()

	var err error

	s.m, err = migrations.NewMigrate(databaseURL)
	if err != nil {
		return err
	}

	err = s.m.Up()
	if err != nil {
		return err
	}

	s.dbConn, err = createDBConn(databaseURL)
	if err != nil {
		return err
	}

	s.srv = createTestServer(s.dbConn)

	startErrCh := make(chan error)

	go func() {
		startErrCh <- s.srv.ListenAndServe()
	}()

	time.Sleep(3 * time.Second)

	select {
	case startErr := <-startErrCh:
		return fmt.Errorf("start server: %w", startErr)
	default:
		go func() {
			<-startErrCh
		}()
	}

	return nil
}

func (s *HandlersTestSuite) signIn() *models.Token {
	var token models.Token

	s.Run("sign in", func() {
		resp, err := doRequest(
			http.MethodPost,
			"/signin",
			&requests.SignIn{
				Login:    "admin",
				Password: "qwerty",
			},
			"",
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		var result handlers.TokenResponse

		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)

		token.AccessToken = result.AccessToken
		token.RefreshToken = result.RefreshToken

		s.Require().NoError(resp.Body.Close())
	})

	return &token
}

func createDBConn(databaseURL string) (*sql.DB, error) {
	dbConn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return dbConn, nil
}

func createTestServer(dbConn *sql.DB) *http.Server {
	h := app.MakeRouter(
		logic.NewAuthService(),
		logic.NewCompanyService(
			goqu.New("postgres", dbConn),
			database.CompanyRepository{},
		),
	)

	return &http.Server{
		Addr:    port,
		Handler: h,
	}
}

func fullURL(u string) string {
	return fmt.Sprintf("http://%v%v/%v", hosturl, port, path.Join(pathPrefix, u))
}

func buildRequest(method, u string, reqBody interface{}) (*http.Request, error) {
	if reqBody == nil {
		return http.NewRequest(method, fullURL(u), nil)
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fullURL(u), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func doRequest(method, u string, reqBody interface{}, accessToken string) (*http.Response, error) {
	req, err := buildRequest(method, u, reqBody)
	if err != nil {
		return nil, err
	}

	if accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+accessToken)
	}

	return http.DefaultClient.Do(req)
}
