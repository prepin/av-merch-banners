package testdb

import (
	"av-merch-shop/config"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	container testcontainers.Container
	Config    config.DBConfig
	ctx       context.Context
}

func NewTestDatabase() (*TestDatabase, error) {
	ctx := context.Background()

	const (
		dbUser     = "test"
		dbPassword = "test"
		dbName     = "testdb"
	)

	// убираем излишнее логирование, а то тест
	testcontainers.Logger = log.New(io.Discard, "", 0)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		),
		AutoRemove: true,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		if termErr := container.Terminate(ctx); termErr != nil {
			return nil, fmt.Errorf("failed to get container host: %w, failed to terminate container: %w", err, termErr)
		}
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		if termErr := container.Terminate(ctx); termErr != nil {
			return nil, fmt.Errorf("failed to get container port: %w, failed to terminate container: %w", err, termErr)
		}
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	dbConfig := config.DBConfig{
		Host:     host,
		Port:     mappedPort.Int(),
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
	}

	return &TestDatabase{
		container: container,
		Config:    dbConfig,
		ctx:       ctx,
	}, nil
}

func (td *TestDatabase) ConnectionString() string {
	return td.Config.GetConnectionString()
}

func (td *TestDatabase) MigrateConnectionString() string {
	return "postgresql://" + td.ConnectionString()[len("postgres://"):]
}

func (td *TestDatabase) TerminateDB() {
	if td.container != nil {
		if err := td.container.Terminate(td.ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}
}
