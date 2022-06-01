package database

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/ihrk/rest-api-task/internal/app/models"
)

type CompanyRepository struct{}

func (CompanyRepository) Create(
	ctx context.Context,
	runner SessionRunner,
	company *models.Company,
) error {
	_, err := runner.Insert("companies").
		Rows(goqu.Record{
			"name":    company.Name,
			"code":    company.Code,
			"country": company.Country,
			"website": company.Website,
			"phone":   company.Phone,
		}).
		Returning("id").
		Executor().
		ScanValContext(ctx, &company.ID)

	if err != nil {
		return fmt.Errorf("create company: %w", err)
	}

	return nil
}

func (CompanyRepository) GetByID(
	ctx context.Context,
	runner SessionRunner,
	companyID int64,
) (models.Company, error) {
	var company Company

	q := runner.Select(
		"id",
		"name",
		"code",
		"country",
		"website",
		"phone",
	).
		From("companies").
		Where(goqu.C("id").Eq(companyID)).
		Executor()

	err := SelectStructQuery(ctx, q, &company)
	if err != nil {
		return models.Company{}, fmt.Errorf("get company by id: %w", err)
	}

	return CompanyToModel(&company), err
}

func (CompanyRepository) Update(
	ctx context.Context,
	runner SessionRunner,
	company *models.Company,
) error {
	_, err := runner.Update("companies").
		Set(goqu.Record{
			"name":    company.Name,
			"code":    company.Code,
			"country": company.Country,
			"website": company.Website,
			"phone":   company.Phone,
		}).
		Where(goqu.C("id").Eq(company.ID)).
		Executor().
		ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("update company: %w", err)
	}

	return nil
}

func (CompanyRepository) Delete(
	ctx context.Context,
	runner SessionRunner,
	companyID int64,
) error {
	_, err := runner.Delete("companies").
		Where(goqu.C("id").Eq(companyID)).
		Executor().
		ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("delete company: %w", err)
	}

	return nil
}

func (CompanyRepository) List(
	ctx context.Context,
	runner SessionRunner,
) ([]models.Company, error) {
	var companies []Company

	err := runner.Select(
		"id",
		"name",
		"code",
		"country",
		"website",
		"phone",
	).
		From("companies").
		Executor().
		ScanStructsContext(ctx, &companies)

	if err != nil {
		return nil, fmt.Errorf("list companies: %w", err)
	}

	return models.ConvertSlice(CompanyToModel, companies), err
}
