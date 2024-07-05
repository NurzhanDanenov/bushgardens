package repository

import (
	"bush/config"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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
	bucketName := "bushgardens-e7339.appspot.com"
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

func (s *Storage) listFilesInFolder(folder string) ([]string, error) {
	ctx := context.Background()
	bucketName := "bushgardens-e7339.appspot.com"
	BucketHandle, err := s.fireStorage.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("error getting default bucket: %v", err)
	}

	query := &storage.Query{Prefix: folder}
	it := BucketHandle.Objects(ctx, query)
	var files []string

	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating through objects: %v", err)
		}
		files = append(files, objAttrs.Name)
	}

	return files, nil
}

func (s *Storage) downloadFile(objectName string) error {
	ctx := context.Background()
	bucketName := "bushgardens-e7339.appspot.com"
	BucketHandle, err := s.fireStorage.Bucket(bucketName)
	if err != nil {
		return err
	}

	obj := BucketHandle.Object(objectName)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("error creating reader: %v", err)
	}
	defer r.Close()

	localPath := filepath.Join("downloads", objectName)
	if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
		return fmt.Errorf("error creating directories: %v", err)
	}

	// Check if the localPath is a directory, and append a file name if it is
	info, err := os.Stat(localPath)
	if err == nil && info.IsDir() {
		localPath = filepath.Join(localPath, "default_filename")
	}

	// Create the local file
	f, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("error copying data: %v", err)
	}

	return nil
}

func (s *Storage) DownloadImage(userId, date string) error {
	folder := userId + "/" + date
	files, err := s.listFilesInFolder(folder)
	if err != nil {
		return fmt.Errorf("error listing files in folder: %v", err)
	}

	for _, file := range files {
		err := s.downloadFile(file)
		if err != nil {
			return fmt.Errorf("error downloading file %s: %v", file, err)
		}
	}

	return nil
}
