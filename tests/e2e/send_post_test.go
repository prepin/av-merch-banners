package e2e

import "net/http"

func (s *E2ETestSuite) TestSendCoin() {
	var tokenResponse tokenResponse
	var errorResponse errorResponse

	authReq := s.client.R().
		SetBody(AuthRequest{
			Username: "employee",
			Password: "password",
		}).
		SetResult(&tokenResponse).
		SetError(&errorResponse)

	resp, err := authReq.Post(s.env.Server.URL + authUrl)
	s.Require().NoError(err, "Failed to get employee token")
	s.Require().Equal(http.StatusOK, resp.StatusCode(), "Failed to get employee token")
	employeeToken := tokenResponse.Token

	type SendRequest struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}

	tests := []struct {
		name          string
		payload       SendRequest
		expectedCode  int
		expectedError string
		token         string
		useToken      bool
	}{
		{
			name: "successful transfer",
			payload: SendRequest{
				ToUser: "director",
				Amount: 100,
			},
			expectedCode: http.StatusOK,
			token:        employeeToken,
			useToken:     true,
		},
		{
			name: "transfer without token",
			payload: SendRequest{
				ToUser: "director",
				Amount: 50,
			},
			expectedCode:  http.StatusUnauthorized,
			expectedError: "authorization header required",
			useToken:      false,
		},
		{
			name: "transfer to non-existent user",
			payload: SendRequest{
				ToUser: "nonexistent",
				Amount: 75,
			},
			expectedCode:  http.StatusNotFound,
			expectedError: "not found",
			token:         employeeToken,
			useToken:      true,
		},
		{
			name: "transfer with insufficient funds",
			payload: SendRequest{
				ToUser: "director",
				Amount: 999999,
			},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "insufficient funds",
			token:         employeeToken,
			useToken:      true,
		},
		{
			name: "transfer with negative amount",
			payload: SendRequest{
				ToUser: "director",
				Amount: -50,
			},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "incorrect amount",
			token:         employeeToken,
			useToken:      true,
		},
		{
			name: "transfer with zero amount",
			payload: SendRequest{
				ToUser: "director",
				Amount: 0,
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "amount is required",
			token:         employeeToken,
			useToken:      true,
		},
		{
			name: "transfer with missing recipient",
			payload: SendRequest{
				ToUser: "",
				Amount: 100,
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "toUser is required",
			token:         employeeToken,
			useToken:      true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := s.client.R().
				SetBody(tt.payload).
				SetError(&errorResponse)

			if tt.useToken {
				req.SetHeader("token", tt.token)
			}

			resp, err := req.Post(s.env.Server.URL + sendCoinUrl)
			s.Require().NoError(err, "Failed to make request")

			s.Assert().Equal(tt.expectedCode, resp.StatusCode(), "Expected status code mismatch")

			if tt.expectedError != "" {
				s.Assert().Equal(tt.expectedError, errorResponse.Error)
			}
		})
	}
}
