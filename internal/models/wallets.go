package models

import "time"

type Balance struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type WalletStruct struct {
	Id       string    `json:"id"`
	UserId   string    `json:"user_id"`
	Balances []Balance `json:"balances"`
}

type Wallet struct {
	Id          string    `gorm:"primaryKey;unique"`
	UserId      string    `gorm:"foreignKey:Id"`
	Currency    string    `gorm:"primaryKey;unique"`
	Amount      float64   `gorm:"default:0;type:numeric(8,2)"`
	LastUpdated time.Time `gorm:"autoUpdateTime"`

	//CurrencyId  uint      `gorm:"foreignKey:Currency"`
}

//type Currency struct {
//	Id   string `gorm:"primaryKey;unique"`
//	Name string
//}
