package x_clone_auth_svc

import (
	"context"
	"errors"
	jwt "x_clone_auth_svc/pkg/jwt"

	user "github.com/qosdil/x_clone_user_svc/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	claimKey = "user_id"
)

type Service interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
	SignUp(ctx context.Context, username, password string) (string, error)
}

type service struct {
	userRepo  Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{userRepo: repo, jwtSecret: jwtSecret}
}

func (s *service) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.FirstByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Return JWT token
	return jwt.GenerateJWT(s.jwtSecret, claimKey, user.ID)
}

func (s *service) SignUp(ctx context.Context, username, password string) (string, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := user.User{
		Username: username,
		Password: string(hashedPassword),
	}
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	// Return JWT token
	return jwt.GenerateJWT(s.jwtSecret, claimKey, user.ID)
}
