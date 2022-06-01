package database

import "github.com/ihrk/rest-api-task/internal/app/models"

type Company struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Code    string `db:"code"`
	Country string `db:"country"`
	Website string `db:"website"`
	Phone   string `db:"phone"`
}

func CompanyToModel(src *Company) models.Company {
	return models.Company{
		ID:      src.ID,
		Name:    src.Name,
		Code:    src.Code,
		Country: src.Country,
		Website: src.Website,
		Phone:   src.Phone,
	}
}
