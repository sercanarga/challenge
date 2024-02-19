package models

import "time"

//type WalletStruct struct {
//	Id       string    `json:"id"`
//	UserId   string    `json:"user_id"`
//	Balances []Balance `json:"balances"`
//}
//
//type Wallet struct {
//	Id      string    `gorm:"primaryKey;unique"`
//	UserId  string    `gorm:"foreignKey:Id;"`
//	Balance []Balance `gorm:"foreignKey:WalletId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
//}
//
//type Balance struct {
//	Id         uint      `gorm:"primaryKey;autoIncrement;unique"`
//	WalletId   string    `gorm:"foreignKey:Id;uniqueIndex:idx_wallet_currency"`
//	Currency   string    `gorm:"uniqueIndex:idx_wallet_currency"`
//	Amount     float64   `gorm:"type:decimal(10,2)"`
//	LastUpdate time.Time `gorm:"autoUpdateTime"`
//}

type WalletStruct struct {
	Id       string    `json:"id"`
	UserId   string    `json:"user_id"`
	Balances []Balance `json:"balances"`
}

type Wallet struct {
	Id      string    `gorm:"primaryKey;unique"`
	UserId  string    `gorm:"foreignKey:Id;"`
	Balance []Balance `gorm:"foreignKey:WalletId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Balance struct {
	Id         uint      `gorm:"primaryKey;autoIncrement;unique" json:"id,omitempty"`
	WalletId   string    `gorm:"foreignKey:Id;uniqueIndex:idx_wallet_currency" json:"wallet_id,omitempty"`
	Currency   string    `gorm:"uniqueIndex:idx_wallet_currency" json:"currency,omitempty"`
	Amount     float64   `gorm:"type:decimal(10,2)"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
}
