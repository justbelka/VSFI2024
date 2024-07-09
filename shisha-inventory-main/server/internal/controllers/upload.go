// controllers/upload.go
package controllers

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"server/internal/initializers"
	"server/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/minio/minio-go/v7"
	log "github.com/rs/zerolog/log"
)

type UploadController struct {
	DB             *gorm.DB
	MinioClient    *minio.Client
	RedisClient    *redis.Client
	RedPandaBroker []string
	Ctx            context.Context
}

func NewUploadController(db *gorm.DB, minioClient *minio.Client, redisClient *redis.Client, redPandaBroker []string, ctx context.Context) *UploadController {
	return &UploadController{
		DB:             db,
		MinioClient:    minioClient,
		RedisClient:    redisClient,
		RedPandaBroker: redPandaBroker,
		Ctx:            ctx,
	}
}

func (uc *UploadController) HandleUpload(c *fiber.Ctx) error {
	// Parse form data
	user := c.Locals("user").(models.User)
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to upload file"})
	}

	// Open the uploaded file
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer fileHeader.Close()

	// Read file content
	fileBytes, err := io.ReadAll(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file"})
	}

	// Compute MD5 hash of the file
	hash := md5.Sum(fileBytes)
	hashValue := hex.EncodeToString(hash[:])

	// Check if the hash exists in Redis
	exists, err := uc.RedisClient.Exists(uc.Ctx, hashValue).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check Redis"})
	}
	if exists > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "This image already exists"})
	}

	// Generate UUID for the image
	imageID := uuid.New().String()

	// Begin transaction
	tx := uc.DB.Begin()

	// Add 1 coin to user's balance
	user.Coins += 1
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user balance"})
	}

	// Upload the file to MinIO
	if _, err := fileHeader.Seek(0, 0); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to seek file"})
	}
	_, err = uc.MinioClient.PutObject(uc.Ctx, "user-images", file.Filename, fileHeader, file.Size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file to MinIO"})
	}

	// Store the image metadata in PostgreSQL
	image := models.Image{
		UUID:       imageID,
		Name:       file.Filename,
		UploadedAt: time.Now(),
		Hash:       hashValue,
		Username:   user.Username,
	}
	if err := uc.DB.Create(&image).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store image metadata"})
	}

	// Store the hash in Redis
	if err := uc.RedisClient.Set(uc.Ctx, hashValue, imageID, 0).Err(); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store hash in Redis"})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	topic := "shisha"
	brokers := uc.RedPandaBroker
	producer, err := initializers.NewProducer(brokers, topic)
	if err != nil {
		return err
	}
	producer.SendUploadMessage(uc.Ctx, user.Username, imageID)

	return c.JSON(fiber.Map{"message": "File uploaded successfully"})
}

// Upload Premium shishki
var imageUrls = []string{
	"https://cveti-piter.ru/wa-data/public/shop/products/24/75/7524/images/49507/49507.750@2x.jpeg",
	"https://www.mirmulchi.ru/images/products/shiski/shishki-sosnovye.jpg",
	"https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/Zapfen_mit_Massband.jpg/1200px-Zapfen_mit_Massband.jpg",
	"https://www.mirmulchi.ru/images/2019/08/14/shishki-sosnovye.jpg",
	"https://nawat.ru/upload/iblock/2d4/2d4c4d01ac817d2b3c623051d98554af.jpg",
}

func (uc *UploadController) HandleBulkDownload() error {
	for i, url := range imageUrls {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Read image data
		imageBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Compute MD5 hash of the image
		hash := md5.Sum(imageBytes)
		hashValue := hex.EncodeToString(hash[:])

		// Check if the hash exists in Redis
		exists, err := uc.RedisClient.Exists(uc.Ctx, hashValue).Result()
		if err != nil {
			return err
		}
		if exists > 0 {
			log.Info().Msg("Image exists in Redis")
			continue // Skip this image if it already exists
		}
		// Generate UUID for the image
		imageID := uuid.New().String()

		imageName := "Cool shishka №" + strconv.Itoa(i+1)
		// Upload the image to MinIO
		reader := bytes.NewReader(imageBytes)
		_, err = uc.MinioClient.PutObject(uc.Ctx, "premium-images", hashValue+".jpg", reader, reader.Size(), minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		})
		if err != nil {
			return err

		}

		// Store the image metadata in PostgreSQL
		image := models.PremiumImage{
			UUID:       imageID,
			Name:       imageName,
			UploadedAt: time.Now(),
			Hash:       hashValue,
			Price:      25,
		}
		if err := uc.DB.Create(&image).Error; err != nil {
			return err
		}

		// Store the hash in Redis
		if err := uc.RedisClient.Set(uc.Ctx, hashValue, imageID, 0).Err(); err != nil {

			return err
		}
	}

	log.Info().Msg("Images downloaded successfully")
	return nil
}
