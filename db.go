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
	createUser(context.Context, io.Writer, string) (string, error)
	addItemToUser(context.Context, io.Writer, Users, ItemParams) error
	userItems(context.Context, io.Writer, string) ([]map[string]interface{}, error)
}

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Users struct {
	BaseModel
	UserID string `gorm:"primaryKey;autoIncrement:false"`
	Name   string
}

type ItemParams struct {
	ItemID string `gorm:"primaryKey;autoIncrement:false"`
}

type userItems struct {
	BaseModel
	ItemParams
	UserID string `gorm:"primaryKey;autoIncrement:false"`
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
func (d dbClient) createUser(ctx context.Context, w io.Writer, u string) (string, error) {

	randomId := genId()

	user := Users{
		BaseModel: BaseModel{},
		UserID:    randomId,
		Name:      u,
	}
	res := d.sc.Debug().Create(&user)

	if res.Error != nil {
		return "", res.Error
	}

	return randomId, nil
}

// add item specified item_id to specific user
func (d dbClient) addItemToUser(ctx context.Context, w io.Writer, u Users, i ItemParams) error {

	ui := userItems{
		BaseModel:  BaseModel{},
		ItemParams: i,
		UserID:     u.UserID,
	}
	log.Printf("%+v", ui)
	res := d.sc.Debug().Create(&ui)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

// get what items the user has
func (d dbClient) userItems(ctx context.Context, w io.Writer, userId string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
