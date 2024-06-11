package repository

import (
	"bushgardens/config"
	"bushgardens/entity"
	"context"
	"errors"
	"fmt"
	"log"
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

//type Store struct {
//	*config.FireDB
//}
//
//func NewStore(db *config.FireDB) *Store {
//	//db := config.FirebaseDB()
//	return &Store{
//		FireDB: db,
//	}
//}

//type Store struct {
//	Client *firebase.App
//}
//
//func NewStore(ctx context.Context) (*Store, error) {
//	client, err := firebase.NewApp(ctx, &firebase.Config{
//		ProjectID: "bushgardens-e7339",
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return &Store{
//		Client: client,
//	}, nil
//}

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

	// Устанавливаем данные пользователя по сгенерированному уникальному ключу
	if err := newUserRef.Set(ctx, user); err != nil {
		return err
	}

	// Выводим уникальный ID
	fmt.Println("Successfully created user with ID:", user.Id)

	return nil
}

//func (s *Store) CreateUser(user *entity.User) error {
//	if s.fireDB == nil {
//		return errors.New("fireDB is nil")
//	}
//	if err := s.fireDB.NewRef("users/"+user.Name).Set(context.Background(), user); err != nil {
//		return err
//	}
//
//	return nil
//}

//func (s *Store) CreateUser(user *entity.User) error {
//	if s.Client == nil {
//		return errors.New("client is nil")
//	}
//	if err := s.NewRef("users/"+user.Name).Set(context.Background(), user); err != nil {
//		return err
//	}
//	return nil
//}

//func (s *Store) DeleteUser(user *entity.User) error {
//	return s.NewRef("user/" + user.Name).Delete(context.Background())
//}

//	func (s *Store) GetByName(b string) (*entity.User, error) {
//		bin := &entity.User{}
//		if err := s.NewRef("user/"+b).Get(context.Background(), bin); err != nil {
//			return nil, err
//		}
//		if bin.Name == "" {
//			return nil, nil
//		}
//		return bin, nil
//	}

//func (r *Store) GetUser(email, password string) (entity.User, error) {
//	var user entity.User
//	ref := r.fireDB.NewRef("users")
//	query := ref.OrderByChild("email").StartAt(email).EndAt(email + "\uf8ff")
//	err := query.Get(context.Background(), &user)
//	if err != nil {
//		return user, err
//	}
//
//	// Verify password
//	if user.Password != password {
//		return user, errors.New("invalid password")
//	}
//
//	return user, nil
//}

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

	// Log the input plain text password and stored hashed password for debugging
	log.Println("Input Plain Text Password:", password)
	log.Println("Stored Hashed Password:", user.Password)

	// Verify password
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	//if err != nil {
	//	log.Println("Error hashing password:", err)
	//	return entity.User{}, err
	//}

	//err = bcrypt.CompareHashAndPassword([]byte(user.Password), hashedPassword)
	//if err != nil {
	//	log.Println("Password comparison failed:", err)
	//	return entity.User{}, errors.New("invalid credentials")
	//}
	if password != user.Password {
		return entity.User{}, errors.New("Password is incorrect")
	}
	return user, nil
}

//func (s *Store) UpdateUser(b string, m map[string]interface{}) error {
//	return s.NewRef("user/"+b).Update(context.Background(), m)
//}
