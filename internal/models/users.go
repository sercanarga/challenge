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
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

type Wallet struct {
	Id          string    `gorm:"primaryKey;unique"`
	UserId      string    `gorm:"foreignKey:Id"`
	Currency    string    `gorm:"primaryKey;unique"`
	Amount      float64   `gorm:"default:0"`
	LastUpdated time.Time `gorm:"autoUpdateTime"`

	//CurrencyId  uint      `gorm:"foreignKey:Currency"`
}

//type Currency struct {
//	Id   string `gorm:"primaryKey;unique"`
//	Name string
//}
