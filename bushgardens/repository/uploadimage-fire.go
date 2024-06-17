package repository

import (
	"bush/config"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
)

type Storage struct {
	fireStorage *config.FireStorage
}

func NewStorage() (*Storage, error) {
	fireStorage, err := config.NewFireStorage()
	if err != nil {
		return nil, err
	}
	return &Storage{
		fireStorage: fireStorage,
	}, nil
}

func (s *Storage) UploadImage(userId string, file multipart.File) (string, error) {
	bucketName := os.Getenv("BUCKET_NAME")
	BucketHandle, err := s.fireStorage.Bucket(bucketName)
	if err != nil {
		return "error getting storage bucket", err
	}

	objectPath := fmt.Sprintf("Images/%s", userId)
	objectHandle := BucketHandle.Object(objectPath)

	writer := objectHandle.NewWriter(context.Background())

	id := uuid.New()
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer writer.Close()

	if _, err := io.Copy(writer, file); err != nil {
		return "error writting to cloud storage", err
	}

	return userId, nil
}
