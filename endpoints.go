package x_clone_auth_srv

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// SignUpRequest defines the structure for the sign-up request
type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest defines the structure for the login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignUpResponse defines the structure for the sign-up response
type SignUpResponse struct {
	UserID string `json:"user_id"`
	Err    error  `json:"error,omitempty"`
}

// LoginResponse defines the structure for the login response
type LoginResponse struct {
	Token string `json:"token"`
	Err   error  `json:"error,omitempty"`
}

// MakeSignUpEndpoint creates the endpoint for the sign-up service
func makeSignUpEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)
		userID, err := svc.SignUp(ctx, req.Username, req.Password)
		return SignUpResponse{UserID: userID, Err: err}, nil
	}
}

// MakeLoginEndpoint creates the endpoint for the login service
func makeLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		token, err := svc.Login(ctx, req.Username, req.Password)
		if err != nil {
			return LoginResponse{Err: err}, nil
		}
		return LoginResponse{Token: token, Err: nil}, nil
	}
}
