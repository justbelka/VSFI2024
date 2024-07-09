package tests

import (
	"context"
	"fmt"
	"log"
	"server/internal/database"
	"server/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestDBStorage(t *testing.T) {
	ctx := context.Background()

	dbName := "shishaDB"
	dbUser := "postgres"
	dbPassword := "yourpassword"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(getContainerImage("postgres:16.3-bookworm")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithInitScripts("../../create_extensions.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	// Clean up the container
	defer func() {
		if err = postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	fmt.Println(dsn)
	require.NoError(t, err)

	db := database.Connect(dsn)
	assert.Nil(t, db)
	t.Run("RunMigrate", func(t *testing.T) {
		err = database.Migrate(&models.User{}, &models.Image{}, &models.PremiumImage{}, &models.Purchase{})
		require.NoError(t, err)
	})

	t.Run("CreateUser", func(t *testing.T) {
		var user models.User
		user.Username = "123"
		user.Password = "123"
		user.Coins = 100
		err = database.Migrate(&models.User{})
		err = database.DB.Create(&user).Error
		require.NoError(t, err)
	})
}
