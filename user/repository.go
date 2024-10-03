package user

import (
	"context"
	"fmt"
	grpcSrv "x_clone_auth_svc/user/grpc/service"
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
	_, err := r.userGrpcClient.Create(ctx, &grpcSrv.Request{Username: user.Username, Password: user.Password})
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	// Some gRPC call here
	// ...

	return User{}, nil
}
