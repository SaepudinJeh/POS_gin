package services

import (
	"POS/internal/models"
	"POS/internal/repositories"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     "user",
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {

	user, err := s.userRepo.FindByEmail(email)

	if err != nil {
		return "", errors.New("email tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("password salah")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", errors.New("gagal membuat token")
	}

	return tokenString, nil
}
