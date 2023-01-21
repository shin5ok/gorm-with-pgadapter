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
	userItems(context.Context, io.Writer, string) ([]ItemParams, error)
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
	ItemID   string `gorm:"primaryKey;autoIncrement:false"`
	ItemName string
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
func (d dbClient) userItems(ctx context.Context, w io.Writer, userId string) ([]ItemParams, error) {
	/*
		sql := `select users.name,items.item_name,user_items.item_id
			from user_items join items on items.item_id = user_items.item_id join users on users.user_id = user_items.user_id
			where user_items.user_id = @user_id`
	*/
	rows, err := d.sc.Debug().Table("user_items").
		Select("users.name as name,items.item_name as item_name,user_items.item_id as item_id, items.updated_at, items.created_at").
		Joins("join items on items.item_id = user_items.item_id join users on users.user_id = user_items.user_id").
		Where("user_items.user_id = ?", userId).
		Rows()

	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}

	defer rows.Close()

	var r ItemParams
	var resultUserItems []ItemParams

	for rows.Next() {
		d.sc.ScanRows(rows, &r)
		resultUserItems = append(resultUserItems, r)
	}

	return resultUserItems, nil
}
