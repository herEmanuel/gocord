package storage

import (
	"errors"

	"github.com/herEmanuel/gocord/pkg/api/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(userVar *models.User, name, email, password string) error {

	var emailExists models.User
	Db.Where("email = ?", email).First(&emailExists)

	if emailExists.Name != "" {
		return errors.New("This email is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	result := Db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	*userVar = newUser

	return nil
}

func Login(userVar *models.User, email, password string) error {

	var userExists models.User
	Db.Where("email = ?", email).First(&userExists)

	if userExists.Name == "" {
		return errors.New("Wrong email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExists.Password), []byte(password))
	if err != nil {
		return errors.New("Wrong email or password")
	}

	*userVar = userExists

	return nil
}
