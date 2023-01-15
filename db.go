package main

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GameUserOperation interface {
	createUser(context.Context, io.Writer, userParams) error
	addItemToUser(context.Context, io.Writer, userParams, itemParams) error
	userItems(context.Context, io.Writer, string) ([]map[string]interface{}, error)
}

type userParams struct {
	userID   string
	userName string
}

type itemParams struct {
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
		// DisableNestedTransaction will turn off the use of Savepoints if gorm
		// detects a nested transaction. Cloud Spanner does not support Savepoints,
		// so it is recommended to set this configuration option to true.
		DisableNestedTransaction: true,
		Logger:                   logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		return dbClient{}, err
	}
	return dbClient{
		sc: db,
	}, nil
}

// create a user
func (d dbClient) createUser(ctx context.Context, w io.Writer, u userParams) error {

	_, err := d.sc.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		sqlToUsers := `INSERT users (user_id, name, created_at, updated_at)
		  VALUES (@userID, @userName, @timestamp, @timestamp)`
		t := time.Now().Format("2006-01-02 15:04:05")
		params := map[string]interface{}{
			"userID":    u.userID,
			"userName":  u.userName,
			"timestamp": t,
		}
		stmtToUsers := spanner.Statement{
			SQL:    sqlToUsers,
			Params: params,
		}
		rowCountToUsers, err := txn.Update(ctx, stmtToUsers)
		_ = rowCountToUsers
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// add item specified item_id to specific user
func (d dbClient) addItemToUser(ctx context.Context, w io.Writer, u userParams, i itemParams) error {

	_, err := d.sc.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		sqlToUsers := `INSERT user_items (user_id, item_id, created_at, updated_at)
		  VALUES (@userID, @itemID, @timestamp, @timestamp)`
		t := time.Now().Format("2006-01-02 15:04:05")
		params := map[string]interface{}{
			"userID":    u.userID,
			"itemId":    i.itemID,
			"timestamp": t,
		}
		stmtToUsers := spanner.Statement{
			SQL:    sqlToUsers,
			Params: params,
		}
		rowCountToUsers, err := txn.Update(ctx, stmtToUsers)
		_ = rowCountToUsers
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// get what items the user has
func (d dbClient) userItems(ctx context.Context, w io.Writer, userID string) ([]map[string]interface{}, error) {

	txn := d.sc.ReadOnlyTransaction()
	defer txn.Close()
	sql := `select users.name,items.item_name,user_items.item_id
		from user_items join items on items.item_id = user_items.item_id join users on users.user_id = user_items.user_id
		where user_items.user_id = @user_id`
	stmt := spanner.Statement{
		SQL: sql,
		Params: map[string]interface{}{
			"user_id": userID,
		},
	}

	iter := txn.Query(ctx, stmt)
	defer iter.Stop()

	results := []map[string]interface{}{}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return results, err
		}
		var userName string
		var itemNames string
		var itemIds string
		if err := row.Columns(&userName, &itemNames, &itemIds); err != nil {
			return results, err
		}

		results = append(results,
			map[string]interface{}{
				"user_name": userName,
				"item_name": itemNames,
				"item_id":   itemIds,
			})

	}

	return results, nil
}
