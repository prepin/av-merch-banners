package e2e

import (
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
	env     *TestEnv
	cleanup func()
	client  *resty.Client
}

var bannerUrl = "/v1/banner"

type bannerItem struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	TagIds    []int  `json:"tag_ids"`
	FeatureId int    `json:"feature_id"`
	IsActive  bool   `json:"is_active"`
}

type bannerCreatedResponse struct {
	BannerId int `json:"banner_id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *E2ETestSuite) SetupTest() {
	s.env, s.cleanup = setupTestEnv(s.T())
	s.client = resty.New()
}

func (s *E2ETestSuite) TearDownTest() {
	s.cleanup()
}

func TestE2Esuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
