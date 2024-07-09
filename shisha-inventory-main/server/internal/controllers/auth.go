// internal/controllers/auth.go
package controllers

import (
	"server/internal/database"
	"server/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	zlog "github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		zlog.Error().Err(err).Msg("failed body parsed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := user.HashPassword(); err != nil {
		zlog.Error().Err(err).Msg("user password hash")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user.Coins = 100
	if err := database.DB.Create(&user).Error; err != nil {
		zlog.Error().Err(err).Msg("db create error")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username already exists"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User registered successfully"})
}

func Login(jwtKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input models.User
		var user models.User
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}
		if !user.CheckPassword(input.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect password"})
		}
		expirationTime := time.Now().Add(72 * time.Hour)
		claims := &Claims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenString, "coins": user.Coins})
	}
}
func AuthRequired(c *fiber.Ctx) error {
	type AuthRequest struct {
		Username string `json:"username"`
		File     string `json:"file"`
	}

	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized. JSON Problem",
		})
	}

	if req.Username == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Verify the user exists in the database
	var user models.User
	result := database.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	c.Locals("user", user)
	return c.Next()
}
