package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ihrk/rest-api-task/internal/app/handlers"
	"github.com/ihrk/rest-api-task/internal/app/requests"
)

type CompanyTestSuite struct {
	HandlersTestSuite
}

func TestCompanyHandlers(t *testing.T) {
	suite.Run(t, new(CompanyTestSuite))
}

func (s *CompanyTestSuite) TestCRUD() {
	token := s.signIn()

	s.Run("create new company", func() {
		resp, err := doRequest(
			http.MethodPost,
			"/companies",
			&requests.CreateCompany{
				Name:    "TestComp",
				Code:    "test-code-123",
				Country: "USA",
				Website: "https://test-website.com",
				Phone:   "+1(23)456-78-90",
			},
			token.AccessToken,
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
	})

	var companyID int64

	s.Run("list all companies", func() {
		resp, err := doRequest(
			http.MethodGet,
			"/companies",
			nil,
			"",
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		var result []handlers.CompanyResponse

		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)
		s.Require().Len(result, 1)
		s.Require().Equal("TestComp", result[0].Name)
		s.Require().Equal("test-code-123", result[0].Code)
		s.Require().Equal("USA", result[0].Country)
		s.Require().Equal("https://test-website.com", result[0].Website)
		s.Require().Equal("+1(23)456-78-90", result[0].Phone)

		companyID = result[0].ID
	})

	s.Run("get single company", func() {
		resp, err := doRequest(
			http.MethodGet,
			fmt.Sprintf("/companies/%d", companyID),
			nil,
			"",
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		var result handlers.CompanyResponse

		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)
		s.Require().Equal(companyID, result.ID)
		s.Require().Equal("TestComp", result.Name)
		s.Require().Equal("test-code-123", result.Code)
		s.Require().Equal("USA", result.Country)
		s.Require().Equal("https://test-website.com", result.Website)
		s.Require().Equal("+1(23)456-78-90", result.Phone)

		s.Require().NoError(resp.Body.Close())
	})

	s.Run("update company", func() {
		resp, err := doRequest(
			http.MethodPut,
			fmt.Sprintf("/companies/%d", companyID),
			&requests.UpdateCompany{
				Name:    "NewTestCompany",
				Code:    "v2-test-code",
				Country: "Mexico",
				Website: "https://upd.test-website.com",
				Phone:   "+52(43)678-90-21",
			},
			"",
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		s.Require().NoError(resp.Body.Close())
	})

	s.Run("check new values", func() {
		resp, err := doRequest(
			http.MethodGet,
			fmt.Sprintf("/companies/%d", companyID),
			nil,
			"",
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		var result handlers.CompanyResponse

		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)
		s.Require().Equal(companyID, result.ID)
		s.Require().Equal("NewTestCompany", result.Name)
		s.Require().Equal("v2-test-code", result.Code)
		s.Require().Equal("Mexico", result.Country)
		s.Require().Equal("https://upd.test-website.com", result.Website)
		s.Require().Equal("+52(43)678-90-21", result.Phone)

		s.Require().NoError(resp.Body.Close())
	})

	s.Run("delete company", func() {
		resp, err := doRequest(
			http.MethodDelete,
			fmt.Sprintf("/companies/%d", companyID),
			nil,
			token.AccessToken,
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		s.Require().NoError(resp.Body.Close())
	})

	s.Run("get company after deletion", func() {
		resp, err := doRequest(
			http.MethodGet,
			fmt.Sprintf("/companies/%d", companyID),
			nil,
			"",
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusNotFound, resp.StatusCode)

		s.Require().NoError(resp.Body.Close())
	})

	s.Run("list again all companies", func() {
		resp, err := doRequest(
			http.MethodGet,
			"/companies",
			nil,
			"",
		)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		var result []handlers.CompanyResponse

		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)
		s.Require().Len(result, 0)

		s.Require().NoError(resp.Body.Close())
	})
}
