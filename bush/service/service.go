package service

import (
	"bush/entity"
	"bush/repository"
	"mime/multipart"
)

type Authorization interface {
	CreateUser(user *entity.User) error
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (string, error)
	//Users(ctx context.Context) ([]*entity.User, error)
}

type UploadImage interface {
	Upload(userId string, file multipart.File) (string, error)
	Download(userId, date string) error
}

type Service struct {
	Authorization
	UploadImage
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		UploadImage:   NewUploadImageService(repos.UploadImage),
	}
}
