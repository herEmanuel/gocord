package jwt

import (
	"time"

	goJwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func CreateToken(userID uuid.UUID, secretKey string) (string, error) {

	claims := goJwt.MapClaims{}
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	preToken := goJwt.NewWithClaims(goJwt.SigningMethodHS256, claims)
	token, err := preToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(token, secretKey string) (uuid.UUID, error) {

	claims := goJwt.MapClaims{}
	_, err := goJwt.ParseWithClaims(token, claims, func(receivedToken *goJwt.Token) (interface{}, error) {

		return []byte(secretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	id, _ := uuid.Parse(claims["userID"].(string))
	return id, nil
}
