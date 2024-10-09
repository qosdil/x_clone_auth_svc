package x_clone_auth_svc

import (
	"context"
	"errors"
	jwt "x_clone_auth_svc/pkg/jwt"

	userGrpcSvc "github.com/qosdil/x_clone_user_svc/grpc/service"
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
	userGrpcClient userGrpcSvc.ServiceClient
	jwtSecret      string
}

func NewService(userGrpcClient userGrpcSvc.ServiceClient, jwtSecret string) Service {
	return &service{userGrpcClient: userGrpcClient, jwtSecret: jwtSecret}
}

func (s *service) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.userGrpcClient.GetByUsername(ctx, &userGrpcSvc.GetByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Return JWT token
	return jwt.GenerateJWT(s.jwtSecret, claimKey, user.Id)
}

func (s *service) SignUp(ctx context.Context, username, password string) (string, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user, err := s.userGrpcClient.Create(ctx, &userGrpcSvc.CreateRequest{
		Username: username, Password: string(hashedPassword),
	})
	if err != nil {
		return "", err
	}

	// Return JWT token
	return jwt.GenerateJWT(s.jwtSecret, claimKey, user.Id)
}
