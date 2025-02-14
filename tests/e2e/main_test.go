package e2e

import (
	"av-merch-shop/tests/setup"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

type E2ETestSuite struct {
	suite.Suite
	env       *setup.TestEnv
	cleanup   func()
	loadData  func()
	cleanData func()
	client    *resty.Client
}

var authURL = "/api/v1/auth"
var sendCoinURL = "/api/v1/sendCoin"
var orderURL = "/api/v1/buy/"

type errorResponse struct {
	Error string `json:"errors"`
}

func (s *E2ETestSuite) SetupTest() {
	s.env, s.cleanup, s.loadData, s.cleanData = setup.CreateTestEnv(s.T())
	s.client = resty.New()
}

func (s *E2ETestSuite) SetupSubTest() {
	s.cleanData()
	s.loadData()
}

func (s *E2ETestSuite) TearDownTest() {
	s.cleanup()
}

func TestE2Esuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
