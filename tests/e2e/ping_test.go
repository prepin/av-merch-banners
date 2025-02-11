package e2e

import (
	"net/http"
)

type pingResponse struct {
	Message string `json:"message"`
}

func (s *E2ETestSuite) TestPingEndpoint() {
	var response pingResponse

	resp, err := s.client.R().
		SetResult(&response).
		Get(s.env.Server.URL + "/api/v1/ping")

	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode())
	s.Assert().Equal("pong", response.Message)
}

func (s *E2ETestSuite) TestTeapotEndpoint() {
	var response pingResponse

	resp, err := s.client.R().
		SetError(&response).
		Get(s.env.Server.URL + "/api/v1/teapot")
	s.Require().NoError(err)

	s.Assert().Equal(http.StatusTeapot, resp.StatusCode())
	s.Assert().Equal("teapot mode", response.Message)
}

// func (s *E2ETestSuite) TestSleepEndpoint() {
// 	var response pingResponse
// 	var errorResponse errorResponse

// 	tests := []struct {
// 		name         string
// 		timeout      string
// 		expectedCode int
// 		expectedMsg  string
// 		isError      bool
// 	}{
// 		{
// 			name:         "should return success for timeout under 45ms",
// 			timeout:      "40",
// 			expectedCode: http.StatusOK,
// 			expectedMsg:  "slept 40 ms",
// 			isError:      false,
// 		},
// 		{
// 			name:         "should return error for timeout over 55ms",
// 			timeout:      "60",
// 			expectedCode: http.StatusRequestTimeout,
// 			expectedMsg:  "timeout error",
// 			isError:      true,
// 		},
// 		{
// 			name:         "should return error, default timeout value",
// 			timeout:      "",
// 			expectedCode: http.StatusRequestTimeout,
// 			expectedMsg:  "timeout error",
// 			isError:      true,
// 		},
// 		{
// 			name:         "should return error, non-int timeout value",
// 			timeout:      "abc",
// 			expectedCode: http.StatusRequestTimeout,
// 			expectedMsg:  "timeout error",
// 			isError:      true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		s.Run(tt.name, func() {

// 			req := s.client.R()
// 			if tt.isError {
// 				req.SetError(&errorResponse)
// 			} else {
// 				req.SetResult(&response)
// 			}

// 			url := s.env.Server.URL + "/api/v1/sleep"
// 			if len(tt.timeout) > 0 {
// 				req.SetQueryParam("timeout", tt.timeout)
// 			}

// 			resp, err := req.Get(url)

// 			s.Require().NoError(err, "Failed to make request")
// 			s.Assert().Equal(tt.expectedCode, resp.StatusCode(), "Expected status code mismatch")

// 			if !tt.isError {
// 				s.Assert().Equal(tt.expectedMsg, response.Message)
// 			}
// 		})
// 	}
// }
