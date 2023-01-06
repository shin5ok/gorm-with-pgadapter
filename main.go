package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var connString = os.Getenv("CONNECTION_STRING")

func main() {
	mode := flag.String("mode", "", "")
	firstName := flag.String("firstname", "", "")
	lastName := flag.String("lastname", "", "")
	flag.Parse()

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

	switch *mode {
	case "createsinger":
		if *firstName != "" && *lastName != "" {
			_, err := CreateSinger(db, *firstName, *lastName)
			if err != nil {
				log.Println(err)
				return
			}
		}
	default:
		fmt.Println("tako")
	}
}

func CreateSinger(db *gorm.DB, firstName, lastName string) (string, error) {
	singer := Singer{
		BaseModel: BaseModel{ID: uuid.NewString()},
		FirstName: sql.NullString{String: firstName, Valid: true},
		LastName:  lastName,
	}
	log.Printf("%+v\n", singer)
	res := db.Create(&singer)
	log.Printf("%+v\n", res)
	if singer.FullName != firstName+" "+lastName {
		return "", fmt.Errorf("unexpected full name for singer: %v", singer.FullName)
	}

	return singer.ID, res.Error
}
