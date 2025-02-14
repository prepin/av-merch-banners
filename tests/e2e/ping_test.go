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
