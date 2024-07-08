package util

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("MY_SECRET_KEY"))

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().AddDate(0, 1, 0) 
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iss": "strive",
		"aud": "users",
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := claims.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (int, error) {
	fmt.Println("INSIDE OF VALIDATE JWT")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(float64); ok {
			return int(userID), nil
		}
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return 0, fmt.Errorf("invalid token")
}
