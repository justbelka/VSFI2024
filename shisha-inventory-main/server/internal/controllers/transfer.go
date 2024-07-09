// internal/controllers/transfer.go
package controllers

import (
	"context"
	"server/internal/database"
	"server/internal/initializers"
	"server/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

func Tranfser(topic string, brokers []string, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type TransferRequest struct {
			FromUsername string `json:"from_username"`
			ToUsername   string `json:"to_username"`
			Amount       int    `json:"amount"`
		}

		var req TransferRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		if req.Amount <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Amount must be greater than zero"})
		}

		var fromUser, toUser models.User

		if err := database.DB.Where("username = ?", req.FromUsername).First(&fromUser).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ebat' ti kto?"})
		}

		if err := database.DB.Where("username = ?", req.ToUsername).First(&toUser).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Recipient not found"})
		}

		if fromUser.Coins < req.Amount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient coins"})
		}

		err := database.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&fromUser).Update("coins", gorm.Expr("coins - ?", req.Amount)).Error; err != nil {
				return err
			}

			if err := tx.Model(&toUser).Update("coins", gorm.Expr("coins + ?", req.Amount)).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "transfer failed"})
		}

		producer, err := initializers.NewProducer(brokers, topic)
		if err != nil {
			return err
		}

		producer.SendTransferMessage(ctx, req.FromUsername, req.ToUsername, req.Amount)

		return c.JSON(fiber.Map{"message": "Transfer successful"})
	}
}
