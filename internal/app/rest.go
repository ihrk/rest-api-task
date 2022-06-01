package app

import (
	"database/sql"
	"net/http"

	"github.com/doug-martin/goqu/v9"

	"github.com/ihrk/rest-api-task/internal/app/database"
	"github.com/ihrk/rest-api-task/internal/app/logic"
)

func Run(cfgPath string) error {
	cfg := LoadConfig(cfgPath)

	dbConn, err := sql.Open("pgx", cfg.DB.URL())
	if err != nil {
		return err
	}

	goquConn := goqu.New("postgres", dbConn)

	router := MakeRouter(
		logic.NewAuthService(),
		logic.NewCompanyService(goquConn, database.CompanyRepository{}),
	)

	return http.ListenAndServe(":8080", router)
}
