package x_clone_auth_svc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	user "github.com/qosdil/x_clone_user_svc/model"
)

var (
	ErrCodeBadRequest = errors.New("bad_request")
	ErrInvalidToken   = errors.New("invalid token")
)

type bodyErrField struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (r *bodyErrField) Error() string {
	return r.Code + "_" + strings.ReplaceAll(r.Message, " ", "_")
}

type errorer interface {
	error() error
}

func decodeAuthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeSignUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	// TODO Implement this as middleware.
	// Validate user inputs
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(req); err != nil {
		if vErr, ok := err.(validator.ValidationErrors); ok {
			// As per standard, we take only the first bad input to respond with.
			// Tag() returns "min" from "min=8".
			msg := strings.ToLower(vErr[0].Field()) + " field " + vErr[0].Tag()

			// Param() returns "8" from "min=8"
			if vErr[0].Param() != "" {
				msg += " " + vErr[0].Param()
			}

			// Interrupt request flow, respond with the custom error type.
			return req, &bodyErrField{
				Code:    ErrCodeBadRequest.Error(),
				Message: msg,
			}
		}
	}

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Respond to bad user inputs
	if v, ok := err.(*bodyErrField); ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": bodyErrField{Code: v.Code, Message: v.Message},
		})
		return
	}

	httpCode, bodyCode, bodyMessage := getErrorStatusAndBody(err)
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": bodyErrField{Code: bodyCode, Message: bodyMessage},
	})
}

// Response encoder
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func getErrorStatusAndBody(err error) (httpCode int, bodyCode, bodyMessage string) {
	if userError, ok := err.(user.Error); ok {
		switch userError.Code {
		case user.ErrCodeUsernameNotAvailable:
			return http.StatusConflict, userError.Code, userError.Error()
		}
	}

	// Others = 500
	bodyMessage = strings.ToLower(http.StatusText(http.StatusInternalServerError))
	bodyCode = strings.ReplaceAll(bodyMessage, " ", "_")
	return http.StatusInternalServerError, bodyCode, bodyMessage
}

// JWTAuthMiddleware is a middleware to validate the JWT token
func JWTAuthMiddleware(secret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Split the token to get the Bearer part
			tokenParts := strings.Split(tokenString, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString = tokenParts[1]

			// Parse the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, nil
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Set user ID in context
			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MakeHTTPHandler sets up the HTTP routes for authentication
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := makeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// Sign up endpoint
	signUpHandler := httptransport.NewServer(
		e.SignUpEndpoint,
		decodeSignUpRequest,
		encodeResponse,
		options...,
	)

	// Auth endpoint
	authHandler := httptransport.NewServer(
		makeAuthEndpoint(s),
		decodeAuthRequest,
		encodeResponse,
	)

	// Register routes
	pathPrefix := "/auth"
	v1Path := "/v1" + pathPrefix
	r.Handle(v1Path, authHandler).Methods("POST")
	r.Handle(v1Path+"/sign-up", signUpHandler).Methods("POST")

	return r
}
