package repository

import (
	"bushgardens/entity"
)

type Authorization interface {
	CreateUser(user *entity.User) error
	//DeleteUser(user *entity.User) error
	GetUser(email, password string) (entity.User, error)
	//UpdateUser(b string, m map[string]interface{}) error
}

type Repository struct {
	Authorization
}

func NewRepository() (*Repository, error) {
	store, err := NewStore()
	if err != nil {
		return nil, err
	}
	return &Repository{
		Authorization: store,
	}, nil
}

//func NewRepository() *Repository {
//	return &Repository{
//		Authorization: NewStore(),
//	}
//}
