// models/image.go
package models

import (
	"time"
)

type Image struct {
	ID         uint      `gorm:"primaryKey"`
	UUID       string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Name       string    `json:"name"`
	UploadedAt time.Time `json:"uploaded_at"`
	Hash       string    `json:"hash"`
	Username   string    `json:"username"`
	CreatedAt  time.Time
}

type PremiumImage struct {
	ID         uint      `gorm:"primaryKey"`
	UUID       string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Name       string    `json:"name"`
	UploadedAt time.Time `json:"uploaded_at"`
	Hash       string    `json:"hash"`
	Price      int       `json:"price" gorm:"default:25"`
	CreatedAt  time.Time
}

type Purchase struct {
	ID        uint   `gorm:"primaryKey"`
	UserName  string `gorm:"not null"`
	ImageID   uint   `gorm:"not null"`
	ImageUUID string `gorm:"type:uuid;default:uuid_generate_v4()" json:"imageuuid"`
	ImageName string `json:"imagename"`
	Hash      string `json:"hash"`
	CreatedAt time.Time
}
