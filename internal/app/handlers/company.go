package handlers

import (
	"context"
	"net/http"

	"github.com/ihrk/rest-api-task/internal/app/models"
	"github.com/ihrk/rest-api-task/internal/app/requests"
)

type CompanyCreator interface {
	CreateCompany(ctx context.Context, r *requests.CreateCompany) error
}

type CreateCompanyHandler struct {
	companyCreator CompanyCreator
}

func NewCreateCompanyHandler(companyCreator CompanyCreator) *CreateCompanyHandler {
	return &CreateCompanyHandler{companyCreator}
}

func (h *CreateCompanyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(r)
	if err != nil {
		RespondError(w, err)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CreateCompanyHandler) serveHTTP(r *http.Request) error {
	var req requests.CreateCompany

	err := ParseRequestJSON(r, &req)
	if err != nil {
		return err
	}

	return h.companyCreator.CreateCompany(r.Context(), &req)
}

type CompanyGetter interface {
	GetCompany(ctx context.Context, companyID int64) (*models.Company, error)
}

type GetCompanyHandler struct {
	companyGetter CompanyGetter
}

func NewGetCompanyHandler(companyGetter CompanyGetter) *GetCompanyHandler {
	return &GetCompanyHandler{companyGetter}
}

func (h *GetCompanyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	company, err := h.serveHTTP(r)
	if err != nil {
		RespondError(w, err)

		return
	}

	respBody := CompanyToReponse(company)

	RespondJSON(w, http.StatusOK, &respBody)
}

func (h *GetCompanyHandler) serveHTTP(r *http.Request) (*models.Company, error) {
	companyID, err := GetIDVar(r, "id")
	if err != nil {
		return nil, err
	}

	return h.companyGetter.GetCompany(r.Context(), companyID)
}

type CompanyUpdater interface {
	UpdateCompany(ctx context.Context, companyID int64, r *requests.UpdateCompany) error
}

type UpdateCompanyHandler struct {
	companyUpdater CompanyUpdater
}

func NewUpdateCompanyHandler(companyUpdater CompanyUpdater) *UpdateCompanyHandler {
	return &UpdateCompanyHandler{companyUpdater}
}

func (h *UpdateCompanyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(r)
	if err != nil {
		RespondError(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UpdateCompanyHandler) serveHTTP(r *http.Request) error {
	companyID, err := GetIDVar(r, "id")
	if err != nil {
		return err
	}

	var req requests.UpdateCompany

	err = ParseRequestJSON(r, &req)
	if err != nil {
		return err
	}

	return h.companyUpdater.UpdateCompany(r.Context(), companyID, &req)
}

type CompanyDeleter interface {
	DeleteCompany(ctx context.Context, companyID int64) error
}

type DeleteCompanyHandler struct {
	companyDeleter CompanyDeleter
}

func NewDeleteCompanyHandler(companyDeleter CompanyDeleter) *DeleteCompanyHandler {
	return &DeleteCompanyHandler{companyDeleter}
}

func (h *DeleteCompanyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(r)
	if err != nil {
		RespondError(w, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *DeleteCompanyHandler) serveHTTP(r *http.Request) error {
	companyID, err := GetIDVar(r, "id")
	if err != nil {
		return err
	}

	return h.companyDeleter.DeleteCompany(r.Context(), companyID)
}

type CompanyLister interface {
	ListCompanies(ctx context.Context) ([]models.Company, error)
}

type ListCompaniesHandler struct {
	companyLister CompanyLister
}

func NewListCompaniesHandler(companyLister CompanyLister) *ListCompaniesHandler {
	return &ListCompaniesHandler{companyLister}
}

func (h *ListCompaniesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	companies, err := h.serveHTTP(r)
	if err != nil {
		RespondError(w, err)

		return
	}

	respBody := models.ConvertSlice(CompanyToReponse, companies)
	if respBody == nil {
		respBody = []CompanyResponse{}
	}

	RespondJSON(w, http.StatusOK, respBody)
}

func (h *ListCompaniesHandler) serveHTTP(r *http.Request) ([]models.Company, error) {
	return h.companyLister.ListCompanies(r.Context())
}
