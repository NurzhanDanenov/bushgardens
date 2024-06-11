package service

import (
	"bushgardens/entity"
	"bushgardens/repository"
)

type Authorization interface {
	CreateUser(user *entity.User) error
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (string, error)
	//Users(ctx context.Context) ([]*entity.User, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
