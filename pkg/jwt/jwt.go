package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(secret, claimKey, claimVal string) (string, error) {
	claims := jwt.MapClaims{
		claimKey: claimVal,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
