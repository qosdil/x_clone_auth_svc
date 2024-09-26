package user

import (
	"context"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user User) error
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

type repository struct{}

func NewRepository() UserRepository {
	return &repository{}
}

func (r *repository) SaveUser(ctx context.Context, user User) error {
	// Some gRPC call here
	// ...

	return nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	// Some gRPC call here
	// ...

	return User{}, nil
}
