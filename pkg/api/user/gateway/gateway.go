package gateway

import (
	"errors"
	"os"

	"github.com/herEmanuel/gocord/pkg/api/models"
	"github.com/herEmanuel/gocord/pkg/api/user/storage"
	"github.com/herEmanuel/gocord/pkg/util/jwt"
)

func Register(name, email, password, confirmPassword string) (map[string]interface{}, error) {

	if password != confirmPassword {
		return nil, errors.New("The passwords don't match")
	}

	var userVar models.User

	err := storage.Register(&userVar, name, email, password)
	if err != nil {
		return nil, err
	}

	token, err := jwt.CreateToken(userVar.ID, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"newUser": userVar, "token": token}, nil
}

func Login(email, password string) (map[string]interface{}, error) {

	var userVar models.User

	err := storage.Login(&userVar, email, password)
	if err != nil {
		return nil, err
	}

	token, err := jwt.CreateToken(userVar.ID, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"user": userVar, "token": token}, nil
}
