package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var connString = os.Getenv("CONNECTION_STRING")

func main() {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		// DisableNestedTransaction will turn off the use of Savepoints if gorm
		// detects a nested transaction. Cloud Spanner does not support Savepoints,
		// so it is recommended to set this configuration option to true.
		DisableNestedTransaction: true,
		Logger:                   logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		fmt.Printf("Failed to open gorm connection: %v\n", err)
		return
	}
	_ = db
}
