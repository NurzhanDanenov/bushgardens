package repository

import (
	"bush/config"
	"bush/entity"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Store struct {
	fireDB *config.FireDB
}

func NewStore() (*Store, error) {
	fireDB, err := config.NewFireDB()
	if err != nil {
		return nil, err
	}
	return &Store{
		fireDB: fireDB,
	}, nil
}

func (s *Store) GetUserByEmail(email string) (*entity.User, error) {
	ctx := context.Background()
	var users map[string]entity.User

	// Получаем всех пользователей
	err := s.fireDB.NewRef("users").Get(ctx, &users)
	if err != nil {
		return nil, err
	}

	// Ищем пользователя с соответствующим email
	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func (s *Store) CreateUser(user *entity.User) error {
	if s.fireDB == nil {
		return errors.New("fireDB is nil")
	}

	ctx := context.Background()
	// Проверяем, существует ли пользователь с таким же email
	existingUser, err := s.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already in use")
	}

	// Создаем ссылку на узел users
	usersRef := s.fireDB.NewRef("users")
	// Генерируем уникальный ключ для нового пользователя
	newUserRef, err := usersRef.Push(ctx, nil)
	if err != nil {
		return err
	}

	// Устанавливаем сгенерированный уникальный ключ как ID пользователя
	user.Id = newUserRef.Key

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	user.LastAttendanceTime = formattedTime
	user.FirstAttendance = formattedTime

	user.TotalAttendance = 0

	// Устанавливаем данные пользователя по сгенерированному уникальному ключу
	if err := newUserRef.Set(ctx, user); err != nil {
		return err
	}

	// Выводим уникальный ID
	fmt.Println("Successfully created user with ID:", user.Id)

	return nil
}

func (r *Store) GetUser(email, password string) (entity.User, error) {
	var userMap map[string]entity.User
	ref := r.fireDB.NewRef("users")
	query := ref.OrderByChild("email").EqualTo(email)
	err := query.Get(context.Background(), &userMap)
	if err != nil {
		log.Println("Error fetching user from Firebase:", err)
		return entity.User{}, err
	}

	if len(userMap) == 0 {
		log.Println("User not found for email:", email)
		return entity.User{}, errors.New("user not found")
	}

	// Extract the first user (there should be only one due to unique email constraint)
	var user entity.User
	for _, u := range userMap {
		user = u
		break
	}

	if password != user.Password {
		return entity.User{}, errors.New("Password is incorrect")
	}
	return user, nil
}
