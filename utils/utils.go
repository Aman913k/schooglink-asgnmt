/*
Package utils provides utility functions for JWT handling and validation.
*/
package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("jwt-key")

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`

	jwt.StandardClaims
}


func GenerateJWT(email, name string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Name:  name,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}


func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}


func StrongPassword(password string) bool {
	return len(password) >= 8
}


func IsValidGmail(email string) bool {
	return strings.HasSuffix(email, "@gmail.com")
}
