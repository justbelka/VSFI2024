package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := redis.RunContainer(ctx,
		testcontainers.WithImage(getContainerImage("redis:latest")),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	require.NoError(t, err)

	// Clean up the container
	defer func() {
		if err = redisContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	state, err := redisContainer.State(ctx)
	require.NoError(t, err)
	fmt.Println(state.Running)
}
