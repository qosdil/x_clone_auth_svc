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

// Response defines the structure for the response
type Response struct {
	Token string `json:"token"`
	Err   error  `json:"error,omitempty"`
}

// MakeSignUpEndpoint creates the endpoint for the sign-up service
func makeSignUpEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)
		token, err := svc.SignUp(ctx, req.Username, req.Password)
		return Response{Token: token, Err: err}, nil
	}
}

// MakeLoginEndpoint creates the endpoint for the login service
func makeLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		token, err := svc.Login(ctx, req.Username, req.Password)
		if err != nil {
			return Response{Err: err}, nil
		}
		return Response{Token: token, Err: nil}, nil
	}
}
