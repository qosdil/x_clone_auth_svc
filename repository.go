package x_clone_auth_svc

import (
	"context"

	userGrpcSvc "github.com/qosdil/x_clone_user_svc/grpc/service"
	user "github.com/qosdil/x_clone_user_svc/model"
)

type Repository interface {
	Create(ctx context.Context, user user.User) error
	FirstByUsername(ctx context.Context, username string) (user.User, error)
}

type repository struct {
	userGrpcClient userGrpcSvc.ServiceClient
}

func NewRepository(userGrpcClient userGrpcSvc.ServiceClient) Repository {
	return &repository{
		userGrpcClient: userGrpcClient,
	}
}

func (r *repository) Create(ctx context.Context, user user.User) error {
	_, err := r.userGrpcClient.Create(ctx, &userGrpcSvc.CreateRequest{Username: user.Username, Password: user.Password})
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FirstByUsername(ctx context.Context, username string) (user.User, error) {
	resp, err := r.userGrpcClient.GetByUsername(ctx, &userGrpcSvc.GetByUsernameRequest{Username: username})
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		ID:       resp.Id,
		Username: resp.Username,
		Password: resp.Password,
	}, nil
}
