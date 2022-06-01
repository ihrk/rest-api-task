package logic

import (
	"context"

	"github.com/doug-martin/goqu/v9"

	"github.com/ihrk/rest-api-task/internal/app/database"
	"github.com/ihrk/rest-api-task/internal/app/models"
	"github.com/ihrk/rest-api-task/internal/app/requests"
)

type CompanyRepo interface {
	Create(
		ctx context.Context,
		runner database.SessionRunner,
		company *models.Company,
	) error

	GetByID(
		ctx context.Context,
		runner database.SessionRunner,
		companyID int64,
	) (models.Company, error)

	Update(
		ctx context.Context,
		runner database.SessionRunner,
		company *models.Company,
	) error

	Delete(
		ctx context.Context,
		runner database.SessionRunner,
		companyID int64,
	) error

	List(
		ctx context.Context,
		runner database.SessionRunner,
	) ([]models.Company, error)
}

type CompanyService struct {
	txRunner database.TxRunner
	repo     CompanyRepo
}

func NewCompanyService(txRunner database.TxRunner, repo CompanyRepo) *CompanyService {
	return &CompanyService{txRunner, repo}
}

func (s *CompanyService) CreateCompany(
	ctx context.Context,
	r *requests.CreateCompany,
) error {
	company := models.Company{
		Name:    r.Name,
		Code:    r.Code,
		Country: r.Country,
		Website: r.Website,
		Phone:   r.Phone,
	}

	err := s.txRunner.WithTx(func(td *goqu.TxDatabase) error {
		txErr := s.repo.Create(ctx, td, &company)

		return txErr
	})

	return err
}

func (s *CompanyService) GetCompany(
	ctx context.Context,
	companyID int64,
) (*models.Company, error) {
	var company models.Company

	err := s.txRunner.WithTx(func(td *goqu.TxDatabase) error {
		var txErr error

		company, txErr = s.repo.GetByID(ctx, td, companyID)

		return txErr
	})

	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (s *CompanyService) UpdateCompany(
	ctx context.Context,
	companyID int64,
	r *requests.UpdateCompany,
) error {
	company := models.Company{
		ID:      companyID,
		Name:    r.Name,
		Code:    r.Code,
		Country: r.Country,
		Website: r.Website,
		Phone:   r.Phone,
	}

	err := s.txRunner.WithTx(func(td *goqu.TxDatabase) error {
		txErr := s.repo.Update(ctx, td, &company)

		return txErr
	})

	return err
}

func (s *CompanyService) DeleteCompany(
	ctx context.Context,
	companyID int64,
) error {
	err := s.txRunner.WithTx(func(td *goqu.TxDatabase) error {
		txErr := s.repo.Delete(ctx, td, companyID)

		return txErr
	})

	return err
}

func (s *CompanyService) ListCompanies(
	ctx context.Context,
) ([]models.Company, error) {
	var companies []models.Company

	err := s.txRunner.WithTx(func(td *goqu.TxDatabase) error {
		var txErr error

		companies, txErr = s.repo.List(ctx, td)

		return txErr
	})

	return companies, err
}
