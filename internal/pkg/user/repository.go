package user

import (
	"context"
	grpcSrv "x_clone_auth_svc/internal/pkg/user/grpc/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user User) error
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

type repository struct {
	userGrpcClient grpcSrv.ServiceClient
}

func NewRepository(userGrpcClient grpcSrv.ServiceClient) UserRepository {
	return &repository{
		userGrpcClient: userGrpcClient,
	}
}

func (r *repository) SaveUser(ctx context.Context, user User) error {
	_, err := r.userGrpcClient.Create(ctx, &grpcSrv.CreateRequest{Username: user.Username, Password: user.Password})
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	user, err := r.userGrpcClient.GetByUsername(ctx, &grpcSrv.GetByUsernameRequest{Username: username})
	if err != nil {
		return User{}, err
	}
	userID, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:       userID,
		Username: user.Username,
		Password: user.Password,
	}, nil
}
