package gateway

import (
	"errors"
	"os"

	"github.com/google/uuid"
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

func ChangePassword(userID uuid.UUID, oldPassword, newPassword, confirmationPassword string) error {

	if confirmationPassword != newPassword {
		return errors.New("The passwords don't match")
	}

	err := storage.ChangePassword(userID, oldPassword, newPassword)
	if err != nil {
		return err
	}

	return nil
}

func AddImage(imagePath string, userID uuid.UUID) error {

	err := storage.AddImage(imagePath, userID)
	if err != nil {
		return err
	}

	return nil
}

func EnterServer(userID uuid.UUID, inviteCode string) (models.Server, error) {

	var server models.Server

	err := storage.EnterServer(&server, userID, inviteCode)
	if err != nil {
		return models.Server{}, err
	}

	return server, nil
}

func LeaveServer(userID, serverID uuid.UUID) error {

	err := storage.LeaveServer(userID, serverID)
	if err != nil {
		return err
	}

	return nil
}
