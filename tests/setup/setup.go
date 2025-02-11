package setup

import (
	"av-merch-shop/config"
	"av-merch-shop/internal/app"
	"av-merch-shop/pkg/auth"
	"av-merch-shop/pkg/database"
	"av-merch-shop/tests/testdb"
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

func SetupTestEnv(t *testing.T) (*TestEnv, func()) {
	t.Helper()

	// Initialize test database
	// May be split this in future so container is create once for test session
	// Instead of separate container for each test case.
	testDB, err := testdb.NewTestDatabase()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	if err := testDB.RunMigrations(); err != nil {
		testDB.Cleanup()
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Clean database (just in case)
	if err := testDB.CleanDatabase(); err != nil {
		testDB.Cleanup()
		t.Fatalf("Failed to clean database: %v", err)
	}

	// Load fixtures
	if err := testDB.LoadFixtures(); err != nil {
		testDB.Cleanup()
		t.Fatalf("Failed to load fixtures %v", err)
	}

	// Create test config
	cfg := &config.Config{
		Logger: InitTestLogger(),
		DB:     testDB.Config, // Use the database config directly
		Server: config.ServerConfig{
			Port:           ":0",
			ReadTimeout:    5,
			WriteTimeout:   5,
			RequestTimeout: 50,
		},
		Auth: config.AuthConfig{
			SecretKey: []byte("test-secret"),
		},
	}

	// Initialize database connection
	db := database.NewDatabase(cfg.DB)

	// Create application
	application := app.New(cfg, db)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	jwtService := auth.NewJWTService(cfg)

	// Register routes
	application.Handlers.RegisterRoutes(router, jwtService)

	// Create test server
	testServer := httptest.NewServer(router)

	env := &TestEnv{
		Config:     cfg,
		App:        application,
		Server:     testServer,
		DB:         testDB,
		HTTPClient: &http.Client{},
		JWT:        jwtService,
	}

	cleanup := func() {
		testServer.Close()
		db.Close()
		testDB.Cleanup()
	}

	return env, cleanup
}

func InitTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
