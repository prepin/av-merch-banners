package e2e

import "net/http"

func (s *E2ETestSuite) TestCreditPost() {
	var tokenResponse tokenResponse
	var errorResponse errorResponse

	authReq := s.client.R().
		SetBody(AuthRequest{
			Username: "director",
			Password: "password",
		}).
		SetResult(&tokenResponse).
		SetError(&errorResponse)

	resp, err := authReq.Post(s.env.Server.URL + authUrl)
	s.Require().NoError(err, "Failed to get auth token")
	s.Require().Equal(http.StatusOK, resp.StatusCode(), "Failed to get auth token")
	directorToken := tokenResponse.Token

	authReq = s.client.R().
		SetBody(AuthRequest{
			Username: "employee",
			Password: "password",
		}).
		SetResult(&tokenResponse).
		SetError(&errorResponse)

	resp, err = authReq.Post(s.env.Server.URL + authUrl)
	s.Require().NoError(err, "Failed to get employee token")
	s.Require().Equal(http.StatusOK, resp.StatusCode(), "Failed to get employee token")
	employeeToken := tokenResponse.Token

	type CreditRequest struct {
		Username string `json:"username"`
		Amount   int    `json:"amount"`
	}

	type CreditResponse struct {
		NewAmount   int    `json:"new_amount"`
		ReferenceID string `json:"reference_id"`
	}

	var creditResponse CreditResponse

	tests := []struct {
		name              string
		payload           CreditRequest
		expectedCode      int
		expectedError     string
		isSuccess         bool
		useToken          bool
		useEmployeeToken  bool
		expectedNewAmount int
	}{
		{
			name: "add credit to existing user",
			payload: CreditRequest{
				Username: "employee",
				Amount:   100,
			},
			expectedCode:      http.StatusCreated,
			isSuccess:         true,
			useToken:          true,
			expectedNewAmount: 100,
		},
		{
			name: "add credit without token",
			payload: CreditRequest{
				Username: "employee",
				Amount:   50,
			},
			expectedCode:  http.StatusUnauthorized,
			expectedError: "authorization header required",
			isSuccess:     false,
			useToken:      false,
		},
		{
			name: "add credit to non-existing user",
			payload: CreditRequest{
				Username: "nonexistent",
				Amount:   75,
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "user not exists",
			isSuccess:     false,
			useToken:      true,
		},
		{
			name: "reduce credit with negative amount",
			payload: CreditRequest{
				Username: "employee",
				Amount:   -50,
			},
			expectedCode:      http.StatusCreated,
			isSuccess:         true,
			useToken:          true,
			expectedNewAmount: 50,
		},
		{
			name: "reduce credit beyond balance",
			payload: CreditRequest{
				Username: "employee",
				Amount:   -1000,
			},
			expectedCode:      http.StatusCreated,
			isSuccess:         true,
			useToken:          true,
			expectedNewAmount: 0,
		},
		{
			name: "add credit with employee token",
			payload: CreditRequest{
				Username: "employee",
				Amount:   100,
			},
			expectedCode:     http.StatusForbidden,
			expectedError:    "admin access required",
			isSuccess:        false,
			useEmployeeToken: true,
		},
		{
			name: "invalid request with missing username",
			payload: CreditRequest{
				Username: "",
				Amount:   100,
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "username is required",
			isSuccess:     false,
			useToken:      true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := s.client.R().
				SetBody(tt.payload).
				SetResult(&creditResponse).
				SetError(&errorResponse)

			if tt.useToken {
				req.SetHeader("token", directorToken)
			} else if tt.useEmployeeToken {
				req.SetHeader("token", employeeToken)
			}

			resp, err := req.Post(s.env.Server.URL + "/api/credit")
			s.Require().NoError(err, "Failed to make request")

			s.Assert().Equal(tt.expectedCode, resp.StatusCode(), "Expected status code mismatch")

			if tt.isSuccess {
				s.Assert().Equal(tt.expectedNewAmount, creditResponse.NewAmount, "New amount should match expected")
				s.Assert().NotEmpty(creditResponse.ReferenceID, "Reference ID should not be empty")
			} else {
				s.Assert().Equal(tt.expectedError, errorResponse.Error)
			}
		})
	}
}
