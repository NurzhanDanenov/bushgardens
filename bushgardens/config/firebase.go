package config

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"fmt"
	"google.golang.org/api/option"
	"os"
)

type FireDB struct {
	*db.Client
}

var fireDB FireDB

func NewFireDB() (*FireDB, error) {
	// Find home directory.
	home, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	opt := option.WithCredentialsFile(home + "\\main\\bushgardens-e7339-firebase-adminsdk-42jdi-c6d0b389d6.json")
	config := &firebase.Config{DatabaseURL: "https://bushgardens-e7339-default-rtdb.asia-southeast1.firebasedatabase.app/"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing database: %v", err)
	}
	return &FireDB{Client: client}, nil
}

//func NewFireDB() (*FireDB, error) {
//	// Find home directory.
//	home, err := os.Getwd()
//	if err != nil {
//		return nil, err
//	}
//	ctx := context.Background()
//	opt := option.WithCredentialsFile(home + "\\main\\bushgardens-e7339-firebase-adminsdk-42jdi-c6d0b389d6.json")
//	config := &firebase.Config{DatabaseURL: "https://bushgardens-e7339-default-rtdb.asia-southeast1.firebasedatabase.app/"}
//	app, err := firebase.NewApp(ctx, config, opt)
//	if err != nil {
//		return nil, fmt.Errorf("error initializing app: %v", err)
//	}
//	client, err := app.Database(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("error initializing database: %v", err)
//	}
//	return &FireDB{Client: client}, nil
//}

//func (db *FireDB) Connect() error {
//	// Find home directory.
//	home, err := os.Getwd()
//	if err != nil {
//		return err
//	}
//	ctx := context.Background()
//	opt := option.WithCredentialsFile(home + "\\main\\bushgardens-e7339-firebase-adminsdk-42jdi-c6d0b389d6.json")
//	config := &firebase.Config{DatabaseURL: "https://bushgardens-e7339-default-rtdb.asia-southeast1.firebasedatabase.app/"}
//	app, err := firebase.NewApp(ctx, config, opt)
//	if err != nil {
//		return fmt.Errorf("error initializing app: %v", err)
//	}
//	client, err := app.Database(ctx)
//	if err != nil {
//		return fmt.Errorf("error initializing database: %v", err)
//	}
//	db.Client = client
//	return nil
//}

func FirebaseDB() *FireDB {
	return &fireDB
}
