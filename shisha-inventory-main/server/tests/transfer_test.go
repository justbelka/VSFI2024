package tests

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"server/internal/controllers"
	"server/internal/database"
	"server/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Set up PostgreSQL container
func setupPostgresTransfer(ctx context.Context) (string, func(), error) {

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
	if err != nil {
		return "", nil, err
	}

	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable", "application_name=test")

	teardown := func() {
		if err = postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}

	return dsn, teardown, nil
}

func setupDatabaseTransfer(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.User{}).Error
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupTestAppTransfer(db *gorm.DB, brokers []string, topic string, ctx context.Context) *fiber.App {
	database.DB = db
	app := fiber.New()

	app.Post("/transfer", controllers.Tranfser(topic, brokers, ctx))

	return app
}

func TestTransfer(t *testing.T) {
	ctx := context.Background()
	dsn, teardown, err := setupPostgresTransfer(ctx)
	assert.NoError(t, err)
	defer teardown()

	db, err := setupDatabaseTransfer(dsn)
	assert.NoError(t, err)

	// Create test users
	fromUser := models.User{
		Username: "sender",
		Password: "password",
		Coins:    100,
	}
	toUser := models.User{
		Username: "receiver",
		Password: "password",
		Coins:    50,
	}

	err = db.Create(&fromUser).Error
	assert.NoError(t, err)
	err = db.Create(&toUser).Error
	assert.NoError(t, err)

	// Brokers and topic for testing purposes
	brokers := []string{"redpanda:9092"}
	topic := "transfers"
	app := setupTestAppTransfer(db, brokers, topic, ctx)

	t.Run("Transfer success", func(t *testing.T) {
		payload := `{"from_username":"sender","to_username":"receiver","amount":50}`
		req, _ := http.NewRequest("POST", "/transfer", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "Transfer successful", body["message"])

		var updatedFromUser, updatedToUser models.User
		db.Where("username = ?", fromUser.Username).First(&updatedFromUser)
		db.Where("username = ?", toUser.Username).First(&updatedToUser)

		assert.Equal(t, fromUser.Coins-50, updatedFromUser.Coins)
		assert.Equal(t, toUser.Coins+50, updatedToUser.Coins)
	})

	t.Run("Transfer with insufficient coins", func(t *testing.T) {
		payload := `{"from_username":"sender","to_username":"receiver","amount":100}`
		req, _ := http.NewRequest("POST", "/transfer", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "Insufficient coins", body["error"])
	})

	t.Run("Transfer with invalid amount", func(t *testing.T) {
		payload := `{"from_username":"sender","to_username":"receiver","amount":-10}`
		req, _ := http.NewRequest("POST", "/transfer", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "Amount must be greater than zero", body["error"])
	})

	t.Run("Transfer from non-existing user", func(t *testing.T) {
		payload := `{"from_username":"nonexistent","to_username":"receiver","amount":10}`
		req, _ := http.NewRequest("POST", "/transfer", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "Ebat' ti kto?", body["error"])
	})

	t.Run("Transfer to non-existing user", func(t *testing.T) {
		payload := `{"from_username":"sender","to_username":"nonexistent","amount":10}`
		req, _ := http.NewRequest("POST", "/transfer", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "Recipient not found", body["error"])
	})
}
