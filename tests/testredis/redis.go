package testredis

import (
	"av-merch-shop/config"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestRedis struct {
	container testcontainers.Container
	Config    config.RedisConfig
	ctx       context.Context
}

func NewTestRedis() (*TestRedis, error) {
	ctx := context.Background()

	testcontainers.Logger = log.New(io.Discard, "", 0)

	req := testcontainers.ContainerRequest{
		Image:        "redis:7",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready to accept connections"),
			wait.ForListeningPort("6379/tcp"),
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

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		if termErr := container.Terminate(ctx); termErr != nil {
			return nil, fmt.Errorf("failed to get container port: %w, failed to terminate container: %w", err, termErr)
		}
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	redisConfig := config.RedisConfig{
		Addr:     fmt.Sprintf("%s:%d", host, mappedPort.Int()),
		Password: "",
		DB:       0,
	}

	return &TestRedis{
		container: container,
		Config:    redisConfig,
		ctx:       ctx,
	}, nil
}

func (tr *TestRedis) TerminateRedis() {
	if tr.container != nil {
		if err := tr.container.Terminate(tr.ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}

}
