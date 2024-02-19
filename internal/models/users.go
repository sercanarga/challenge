package models

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	Id        string `gorm:"primaryKey;unique"`
	Name      string
	Wallets   []Wallet  `gorm:"foreignKey:UserId"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt
}
