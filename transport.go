package x_clone_auth_srv

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var ErrInvalidToken = errors.New("invalid token")

// MakeHTTPHandler sets up the HTTP routes for authentication
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	// Sign up endpoint
	signUpHandler := httptransport.NewServer(
		makeSignUpEndpoint(s),
		decodeSignUpRequest,
		encodeResponse,
	)

	// Login endpoint
	loginHandler := httptransport.NewServer(
		makeLoginEndpoint(s),
		decodeLoginRequest,
		encodeResponse,
	)

	// Register routes
	r.Handle("/auth/signup", signUpHandler).Methods("POST")
	r.Handle("/auth/login", loginHandler).Methods("POST")

	return r
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

func decodeSignUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// Response encoder
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
