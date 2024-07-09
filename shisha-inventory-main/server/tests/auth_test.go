package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/internal/controllers"
	"server/internal/database"
	"server/internal/models"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var jwtKey = "your_secret_key"

// Set up PostgreSQL container
func setupPostgres(ctx context.Context) (string, func(), error) {
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

func setupDatabase(dsn string) (*gorm.DB, error) {
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

func setupTestApp(db *gorm.DB) *fiber.App {
	database.DB = db
	app := fiber.New()

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login(jwtKey))
	app.Post("/auth", controllers.AuthRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		return c.JSON(fiber.Map{"username": user.Username})
	})

	return app
}

func TestAuthControllers(t *testing.T) {
	ctx := context.Background()
	dsn, teardown, err := setupPostgres(ctx)
	assert.NoError(t, err)
	defer teardown()

	db, err := setupDatabase(dsn)
	assert.NoError(t, err)

	app := setupTestApp(db)

	t.Run("Register", func(t *testing.T) {
		payload := `{"username":"testuser","password":"testpass"}`
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "User registered successfully", body["message"])
	})

	t.Run("Login", func(t *testing.T) {
		payload := `{"username":"testuser","password":"testpass"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.NotEmpty(t, body["token"])
		assert.Equal(t, float64(100), body["coins"])
	})

	t.Run("AuthRequired", func(t *testing.T) {
		// Login to get a token
		payload := `{"username":"testuser","password":"testpass"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)

		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		token := body["token"].(string)

		// Test AuthRequired endpoint
		authPayload := `{"username":"testuser","file":"somefile"}`
		authReq, _ := http.NewRequest("POST", "/auth", strings.NewReader(authPayload))
		authReq.Header.Set("Content-Type", "application/json")
		authReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		authResp, err := app.Test(authReq)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, authResp.StatusCode)

		var authBody map[string]interface{}
		err = json.NewDecoder(authResp.Body).Decode(&authBody)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", authBody["username"])
	})
}
