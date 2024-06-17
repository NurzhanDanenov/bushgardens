package repository

import (
	"bush/entity"
	"mime/multipart"
)

type Authorization interface {
	CreateUser(user *entity.User) error
	//DeleteUser(user *entity.User) error
	GetUser(email, password string) (entity.User, error)
	//UpdateUser(b string, m map[string]interface{}) error
}

type UploadImage interface {
	UploadImage(userId string, file multipart.File) (string, error)
}

type Repository struct {
	Authorization
	UploadImage
}

func NewRepository() (*Repository, error) {
	store, err := NewStore()
	if err != nil {
		return nil, err
	}
	storage, err := NewStorage()
	if err != nil {
		return nil, err
	}
	return &Repository{
		Authorization: store,
		UploadImage:   storage,
	}, nil
}
