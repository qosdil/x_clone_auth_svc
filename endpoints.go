package x_clone_auth_svc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// AuthRequest defines the structure for the auth request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Endpoints struct {
	SignUpEndpoint endpoint.Endpoint
}

// Response defines the structure for the response
type Response struct {
	Token string `json:"token"`
	Err   error  `json:"error,omitempty"`
}

// SignUpRequest defines the structure for the sign-up request
type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=8"`
	Password string `json:"password" validate:"required,min=8"`
}

// makeAuthEndpoint creates the endpoint for the auth service
func makeAuthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthRequest)
		token, err := s.Authenticate(ctx, req.Username, req.Password)
		if err != nil {
			return Response{Err: err}, nil
		}
		return Response{Token: token, Err: nil}, nil
	}
}

func makeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		SignUpEndpoint: makeSignUpEndpoint(s),
	}
}

// makeSignUpEndpoint creates the endpoint for the sign-up service
func makeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SignUpRequest)
		token, e := s.SignUp(ctx, req.Username, req.Password)
		return Response{Token: token, Err: e}, nil
	}
}
