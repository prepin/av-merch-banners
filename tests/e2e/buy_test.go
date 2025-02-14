package e2e

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func (s *E2ETestSuite) TestOrder() {
	var tokenResponse tokenResponse
	var errorResponse errorResponse

	authReq := s.client.R().
		SetBody(AuthRequest{
			Username: "employee",
			Password: "password",
		}).
		SetResult(&tokenResponse).
		SetError(&errorResponse)

	resp, err := authReq.Post(s.env.Server.URL + authURL)
	s.Require().NoError(err, "Failed to get employee token")
	s.Require().Equal(http.StatusOK, resp.StatusCode(), "Failed to get employee token")
	employeeToken := tokenResponse.Token

	tests := []struct {
		name          string
		itemName      string
		expectedCode  int
		expectedError string
		token         string
		useToken      bool
	}{
		{
			name:         "successful order",
			itemName:     "t-shirt",
			expectedCode: http.StatusOK,
			token:        employeeToken,
			useToken:     true,
		},
		{
			name:          "order without token",
			itemName:      "t-shirt",
			expectedCode:  http.StatusUnauthorized,
			expectedError: "authorization header required",
			useToken:      false,
		},
		{
			name:          "order non-existent item",
			itemName:      "nonexistent",
			expectedCode:  http.StatusNotFound,
			expectedError: "not found",
			token:         employeeToken,
			useToken:      true,
		},
		{
			name:          "order with empty item name",
			itemName:      "",
			expectedCode:  http.StatusBadRequest,
			expectedError: "item is required",
			token:         employeeToken,
			useToken:      true,
		},
	}

	methods := []string{"GET", "POST"}

	for _, method := range methods {
		for _, tt := range tests {
			testName := fmt.Sprintf("%s - %s", method, tt.name)
			s.Run(testName, func() {
				req := s.client.R().
					SetError(&errorResponse)

				if tt.useToken {
					req.SetHeader("Authorization", "Bearer "+tt.token)
				}

				url := s.env.Server.URL + orderURL + tt.itemName

				var resp *resty.Response
				var err error

				if method == "GET" {
					resp, err = req.Get(url)
				} else {
					resp, err = req.Post(url)
				}

				s.Require().NoError(err, "Failed to make request")
				s.Assert().Equal(tt.expectedCode, resp.StatusCode(), "Expected status code mismatch")

				if tt.expectedError != "" {
					s.Assert().Equal(tt.expectedError, errorResponse.Error)
				}

				if method == "GET" && tt.expectedCode == http.StatusOK {
					deprecationHeader := resp.Header().Get("Deprecation")
					s.Assert().NotEmpty(deprecationHeader, "Deprecation header should be present for GET requests")
				}
			})
		}
	}
}
