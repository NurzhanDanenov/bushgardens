package service

import (
	"bushgardens/entity"
	"bushgardens/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	salt       = "kdowkdew0321e2ndis"
	signingKey = "koswnsndq930231niisnm"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user *entity.User) error {
	hashedPassword := GeneratePasswordHash(user.Password)
	//if err != nil {
	//	return err
	//}
	user.Password = hashedPassword
	//user.Password, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	return entity.User{}, err
	//}
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(email, password string) (string, error) {
	hashedPassword := GeneratePasswordHash(password)
	//if err != nil {
	//	return "", err
	//}
	user, err := s.repo.GetUser(email, hashedPassword)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return "null", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "null", errors.New("token claims are not of type *tokenClaims")
	}
	//fmt.Println(claims.UserId)
	return claims.UserId, nil
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//if err != nil {
	//	return "", err
	//}
	//return string(hashedPassword), nil
}

//func (s *AuthService) Users(ctx context.Context) ([]*entity.User, error) {
//	return s.repo.GetUsers(ctx)
//}
