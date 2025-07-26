package service

import (
	"time"

	"github.com/ImanMontajabi/Gomodoro/internal/model"
	"github.com/ImanMontajabi/Gomodoro/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type authService struct {
	repo      repository.UserRepository
	secretKey string
}

func NewAuthService(repo repository.UserRepository, secretKey string) AuthService {
	return &authService{repo, secretKey}
}

func (s *authService) Register(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{Username: username, Password: string(hashed)}
	return s.repo.Create(user)
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(s.secretKey))
}
