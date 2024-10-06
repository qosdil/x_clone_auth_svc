package x_clone_auth_svc

import (
	"context"
	"errors"
	"x_clone_auth_svc/internal/pkg/user"
	jwt "x_clone_auth_svc/pkg/jwt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	claimKey = "user_id"
)

type Service interface {
	SignUp(ctx context.Context, username, password string) (string, error)
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type service struct {
	userRepo  user.Repository
	jwtSecret string
}

func NewService(repo user.Repository, jwtSecret string) Service {
	return &service{userRepo: repo, jwtSecret: jwtSecret}
}

func (s *service) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Return JWT token
	return jwt.GenerateJWT(s.jwtSecret, claimKey, user.ID.Hex())
}

func (s *service) SignUp(ctx context.Context, username, password string) (string, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := user.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: string(hashedPassword),
	}
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	// Return JWT token
	return jwt.GenerateJWT(s.jwtSecret, claimKey, user.ID.Hex())
}
