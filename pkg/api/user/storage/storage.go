package storage

import (
	"errors"

	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/api/models"
	serverStorage "github.com/herEmanuel/gocord/pkg/api/server/storage"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(userVar *models.User, name, email, password string) error {

	var emailExists models.User
	result := Db.Where("email = ?", email).First(&emailExists)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
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

	result = Db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	*userVar = newUser

	return nil
}

func Login(userVar *models.User, email, password string) error {

	var userExists models.User
	result := Db.Where("email = ?", email).First(&userExists)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("Wrong email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExists.Password), []byte(password))
	if err != nil {
		return errors.New("Wrong email or password")
	}

	*userVar = userExists

	return nil
}

func ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {

	var user models.User

	Db.First(&user, "id = ?", userID)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("Invalid old password")
	}

	if oldPassword == newPassword {
		return errors.New("The new password can not be the same as the old one")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	result := Db.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func AddImage(imagePath string, userID uuid.UUID) error {

	var user models.User

	Db.First(&user, "id = ?", userID)

	user.Avatar = imagePath

	result := Db.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetUserInfo(userVar *models.User, userID uuid.UUID) error {

	var user models.User

	result := Db.Preload("Servers").
		Select("id", "created_at", "name", "avatar").
		First(&user, "id = ?", userID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This user doesn't exist")
	}

	*userVar = user

	return nil
}

func EnterServer(serverVar *models.Server, userID uuid.UUID, inviteCode string) error {

	var user models.User
	var server models.Server

	result := Db.Where("invite_code = ?", inviteCode).
		Preload("Members").
		First(&server)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("Invalid invite code")
	}

	for _, serverMember := range server.Members {
		if serverMember.ID == userID {
			return errors.New("You are already in this server")
		}
	}

	Db.First(&user, "id = ?", userID)

	user.Servers = append(user.Servers, server)

	result = Db.Omit("Servers.*").Save(&user)
	if result.Error != nil {
		return result.Error
	}

	err := serverStorage.GetServer(serverVar, server.ID)
	if err != nil {
		return err
	}

	return nil
}

func LeaveServer(userID, serverID uuid.UUID) error {

	var user models.User
	var server models.Server

	result := Db.First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	Db.Preload("Servers").First(&user, "id = ?", userID)

	isInServer := false
	for _, userServer := range user.Servers {
		if userServer.ID == server.ID {
			isInServer = true
			break
		}
	}
	if !isInServer {
		return errors.New("You are not in this server")
	}

	err := Db.Model(&user).Association("Servers").Delete(&server)
	if err != nil {
		return err
	}

	return nil
}
