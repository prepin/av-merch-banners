package e2e

import (
	"net/http"
)

type AuthRequest struct {
	Username string
	Password string
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (s *E2ETestSuite) TestAuthPost() {

	var tokenResponse tokenResponse
	var errorResponse errorResponse

	tests := []struct {
		name           string
		payload        AuthRequest
		expectedCode   int
		expectedError  string
		expectedResult string
		isSuccess      bool
	}{
		{
			name: "existing user",
			payload: AuthRequest{
				Username: "employee",
				Password: "password",
			},
			expectedCode: http.StatusOK,
			isSuccess:    true,
		},
		{
			name: "new user",
			payload: AuthRequest{
				Username: "new_user",
				Password: "password",
			},
			expectedCode: http.StatusOK,
			isSuccess:    true,
		},
		{
			name: "empty password",
			payload: AuthRequest{
				Username: "employee",
				Password: "",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "password is required",
			isSuccess:     false,
		},
		{
			name: "empty username",
			payload: AuthRequest{
				Username: "",
				Password: "password",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "username is required",
			isSuccess:     false,
		},
		{
			name: "bad password",
			payload: AuthRequest{
				Username: "employee",
				Password: "wrong-password",
			},
			expectedError: "wrong password",
			expectedCode:  http.StatusUnauthorized,
			isSuccess:     false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := s.client.R().
				SetBody(tt.payload).
				SetResult(&tokenResponse).
				SetError(&errorResponse)

			resp, err := req.Post(s.env.Server.URL + authURL)
			s.Require().NoError(err, "Failed to make request")

			s.Assert().Equal(tt.expectedCode, resp.StatusCode(), "Expected status code mismatch")
			if tt.isSuccess {
				claims, err := s.env.JWT.ValidateToken(tokenResponse.Token)
				s.Require().NoError(err, "Token validation failed")

				s.Assert().Equal(tt.payload.Username, claims.Username)
			} else {
				s.Assert().Equal(tt.expectedError, errorResponse.Error)
			}
		})
	}

}
