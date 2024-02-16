package durable

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type ConnectionInfo struct {
	Host     string
	User     string
	Password string
	Database string
	Port     string
	SSLMode  string
}

func ConnectDB(c *ConnectionInfo) error {
	var err error

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Istanbul", c.Host, c.User, c.Password, c.Database, c.Port, c.SSLMode),
		PreferSimpleProtocol: true,
	}), &gorm.Config{SkipDefaultTransaction: false})

	if err != nil {
		return err
	}
	return nil
}

func Connection() *gorm.DB {
	return db
}
