package x_clone_auth_svc

import (
	"context"
	"errors"
	"time"
	"x_clone_auth_svc/internal/pkg/user"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUp(ctx context.Context, username, password string) (string, error)
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type service struct {
	repo      user.UserRepository
	jwtSecret string
}

func NewService(repo user.UserRepository, jwtSecret string) Service {
	return &service{repo: repo, jwtSecret: jwtSecret}
}

func (s *service) SignUp(ctx context.Context, username, password string) (string, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := user.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: string(hashedPassword),
	}
	err := s.repo.SaveUser(ctx, user)
	if err != nil {
		return "", err
	}

	// Return JWT token
	return s.generateJWT(user)
}

func (s *service) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Return JWT token
	return s.generateJWT(user)
}

func (s *service) generateJWT(user user.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
