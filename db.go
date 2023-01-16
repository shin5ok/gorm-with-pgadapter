package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GameUserOperation interface {
	createUser(context.Context, io.Writer, string) error
	addItemToUser(context.Context, io.Writer, Users, itemParams) error
	userItems(context.Context, io.Writer, string) ([]map[string]interface{}, error)
}

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Users struct {
	BaseModel
	UserId string `gorm:"primaryKey;autoIncrement:false"`
	Name   string
}

type itemParams struct {
	BaseModel
	itemID    string
	itemPrice int
}

type dbClient struct {
	sc *gorm.DB
}

func genId() string {
	newUUID, _ := uuid.NewRandom()
	return newUUID.String()
}

func newClient(ctx context.Context, spannerString string) (dbClient, error) {

	db, err := gorm.Open(postgres.Open(spannerString), &gorm.Config{
		DisableNestedTransaction: true,
		//Logger:                   logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		return dbClient{}, err
	}
	return dbClient{
		sc: db,
	}, nil
}

// create a user
func (d dbClient) createUser(ctx context.Context, w io.Writer, u string) error {

	randomId := genId()
	log.Printf("%+v\n", randomId)

	user := Users{
		BaseModel: BaseModel{},
		UserId:    randomId,
		Name:      u,
	}
	log.Printf("%+v\n", user)
	res := d.sc.Debug().Create(&user)
	log.Printf("%+v\n", res)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

// add item specified item_id to specific user
func (d dbClient) addItemToUser(ctx context.Context, w io.Writer, u Users, i itemParams) error {
	return nil
}

// get what items the user has
func (d dbClient) userItems(ctx context.Context, w io.Writer, userId string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
