package config

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"firebase.google.com/go/v4/storage"
	"fmt"
	"google.golang.org/api/option"
	"os"
)

type FireDB struct {
	*db.Client
}

type FireStorage struct {
	*storage.Client
}

var fireDB FireDB

func InitializeFirebaseApp() (*firebase.App, error) {
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
	//client, err := app.Database(ctx)
	//if err != nil {
	//	return nil, fmt.Errorf("error initializing database: %v", err)
	//}
	//return &FireDB{Client: client}, nil
	return app, nil
}

func NewFireDB() (*FireDB, error) {
	// Initialize Firebase Realtime Database client.
	ctx := context.Background()
	app, err := InitializeFirebaseApp()
	client, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing database client: %v", err)
	}

	return &FireDB{Client: client}, nil
}

// NewFireStorage initializes the FireStorage client and returns it.
func NewFireStorage() (*FireStorage, error) {
	// Initialize Firebase Storage client.
	ctx := context.Background()
	app, err := InitializeFirebaseApp()
	client, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing storage client: %v", err)
	}

	return &FireStorage{Client: client}, nil
}

func FirebaseDB() *FireDB {
	return &fireDB
}
