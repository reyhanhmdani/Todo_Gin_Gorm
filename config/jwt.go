package config

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// secret key untuk signing token
var JwtKey = []byte("secret_key")

// payload untuk token
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// fungsi untuk membuat token
func CreateToken(username string) (string, error) {
	// mengatur waktu kadaluwarsa token
	expirationTime := time.Now().Add(10 * time.Minute)

	// membuat claims
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// membuat token dengan signing method HS256 dan secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
