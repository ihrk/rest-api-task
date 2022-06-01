package handlers

import "github.com/ihrk/rest-api-task/internal/app/models"

type CompanyResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}

func CompanyToReponse(src *models.Company) CompanyResponse {
	return CompanyResponse{
		ID:      src.ID,
		Name:    src.Name,
		Code:    src.Code,
		Country: src.Country,
		Website: src.Website,
		Phone:   src.Phone,
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func TokenToResponse(src *models.Token) TokenResponse {
	return TokenResponse{
		AccessToken:  src.AccessToken,
		RefreshToken: src.RefreshToken,
	}
}
