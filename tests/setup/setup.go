package setup

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/app"
	"av-merch-shop/pkg/auth"
	"av-merch-shop/pkg/database"
	"av-merch-shop/pkg/redis"
	"av-merch-shop/tests/testdb"
	"av-merch-shop/tests/testredis"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

type TestEnv struct {
	Config     *config.Config
	App        *app.Application
	Server     *httptest.Server
	DB         *testdb.TestDatabase
	HTTPClient *http.Client
	JWT        *auth.JWTService
}

// Возвращает Окружение, функцию остановки тест-контейнера,
// функцию загрузки сидов, функцию очистки сидов.
func CreateTestEnv(t *testing.T) (env *TestEnv, cleanup, loadSeeds, dropSeeds func()) {
	t.Helper()

	testDB, err := testdb.NewTestDatabase()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	testRedis, err := testredis.NewTestRedis()
	if err != nil {
		testDB.TerminateDB()
		testRedis.TerminateRedis()
		t.Fatalf("Failed to create test redis: %v", err)
	}

	loadSeedData := func() {
		if err := testDB.RunMigrations(); err != nil {
			testDB.TerminateDB()
			testRedis.TerminateRedis()
			t.Fatalf("Failed to run migrations: %v", err)
		}
		if err := testDB.LoadFixtures(); err != nil {
			testDB.TerminateDB()
			testRedis.TerminateRedis()
			t.Fatalf("Failed to load fixtures %v", err)
		}
	}

	dropSeedData := func() {
		if err := testDB.CleanDatabase(); err != nil {
			testDB.TerminateDB()
			testRedis.TerminateRedis()
			t.Fatalf("Failed to clean database: %v", err)
		}
	}

	loadSeedData()

	cfg := &config.Config{
		Logger: InitTestLogger(),
		DB:     testDB.Config,
		Redis:  testRedis.Config,
		Server: config.ServerConfig{
			Port:         ":0",
			ReadTimeout:  5,
			WriteTimeout: 5,
		},
		Auth: config.AuthConfig{
			SecretKey: []byte("test-secret"),
		},
	}

	db := database.NewDatabase(cfg.DB)
	redis := redis.NewRedis(cfg.Redis, cfg.Logger)

	application := app.New(cfg, db, redis)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	jwtService := auth.NewJWTService(cfg)

	application.Handlers.RegisterRoutes(router, jwtService)

	testServer := httptest.NewServer(router)

	env = &TestEnv{
		Config:     cfg,
		App:        application,
		Server:     testServer,
		DB:         testDB,
		HTTPClient: &http.Client{},
		JWT:        jwtService,
	}

	cleanup = func() {
		testServer.Close()
		db.Close()
	}

	return env, cleanup, loadSeedData, dropSeedData
}

func InitTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
