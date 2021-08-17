package middlewares

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ehabterra/flash_api/api/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/ehabterra/flash_api/internal/constants"
	"github.com/google/martian/log"
)

func ValidateHeader(bearerHeader string) (*models.Principle, error) {
	bearerToken := strings.Split(bearerHeader, " ")[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error decoding token")
		}
		return []byte(constants.JWTSecretKey), nil
	})
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	if token.Valid {
		email := claims["email"].(string)
		id := claims["id"].(string)
		username := claims["username"].(string)

		return &models.Principle{
			Email:    &email,
			ID:       &id,
			Username: &username,
		}, nil
	}
	return nil, errors.New("invalid token")
}

func GenerateJWT(userID, email, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = userID
	claims["email"] = email
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 120).Unix()

	tokenString, err := token.SignedString([]byte(constants.JWTSecretKey))
	if err != nil {
		log.Errorf("Error generating Token: " + err.Error())
		return "", err
	}
	return tokenString, nil
}
