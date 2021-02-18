package storage

import (
	"errors"

	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/api/models"
	"gorm.io/gorm"
)

func CreateServer(serverVar *models.Server, userID uuid.UUID, name string) error {

	var creator models.User

	Db.First(&creator, "id = ?", userID)

	newServer := models.Server{
		Name:       name,
		InviteCode: "",
		Members:    []models.User{creator},
		Admins:     []models.User{creator},
	}

	result := Db.Omit("Members.*", "Admins.*").
		Create(&newServer)
	if result.Error != nil {
		return result.Error
	}

	newChannel := models.Channel{
		Name:       "general",
		Permission: "public",
		Server:     newServer.ID,
	}

	result = Db.Create(&newChannel)
	if result.Error != nil {
		return result.Error
	}

	*serverVar = newServer

	return nil
}

func AddImage(imagePath string, serverID uuid.UUID, userID uuid.UUID) error {

	var server models.Server

	result := Db.Preload("Admins").
		First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	isAdmin := false
	for _, admin := range server.Admins {
		if admin.ID == userID {
			isAdmin = true
		}
	}
	if !isAdmin {
		return errors.New("You don't have permission to perform that task")
	}

	server.Picture = imagePath

	result = Db.Save(&server)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// func EditServer(serverVar *models.Server, serverID uuid.UUID) error {

// 	var server models.Server

// 	result := Db.Model(&models.Server{}).First(&server, "id = ?", serverID)

// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return errors.New("This server doesn't exist")
// 	}

// 	return nil
// }

func GetServer(serverVar *models.Server, serverID uuid.UUID) error {

	var server models.Server

	result := Db.Preload("Channels").
		Preload("Members").
		Preload("Roles").
		First(&server, "id = ?", serverID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	*serverVar = server

	return nil
}

//TODO: create a middleware to check whether a person is an admin or not
func CreateChannel(channelVar *models.Channel, serverID uuid.UUID, name, permission string) error {

	var server models.Server

	result := Db.Model(&models.Server{}).First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	newChannel := models.Channel{
		Name:       name,
		Permission: permission,
		Server:     server.ID,
	}

	result = Db.Model(&models.Channel{}).Create(&newChannel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetChannelMessages(userID uuid.UUID, channelID uuid.UUID) ([]models.Message, error) {

	var channel models.Channel
	var user models.User

	Db.Preload("Servers").
		First(&user, "id = ?", userID)

	result := Db.Preload("Messages").
		Preload("Messages.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "avatar")
		}).
		First(&channel, "id = ?", channelID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []models.Message{}, errors.New("This channel doesn't exist")
	}

	isInServer := false
	for _, server := range user.Servers {
		if server.ID == channel.Server {
			isInServer = true
		}
	}

	if !isInServer {
		return []models.Message{}, errors.New("You are not in this server")
	}

	//TODO: check whether a user has permission to see the messages of a channel or not

	return channel.Messages, nil
}

func SendMessage(creatorID uuid.UUID, channelID uuid.UUID, content, messageType string) error {

	var channel models.Channel
	var creator models.User

	Db.Preload("Servers").
		First(&creator, "id = ?", creatorID)

	result := Db.First(&channel, "id = ?", channelID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This channel doesn't exist")
	}

	isInServer := false
	for _, server := range creator.Servers {
		if server.ID == channel.Server {
			isInServer = true
		}
	}

	if !isInServer {
		return errors.New("You are not in this server")
	}

	//TODO: check if the user has permission to send the message in this channel

	newMessage := models.Message{
		Content: content,
		Type:    messageType,
		User:    creator,
		Channel: channel.ID,
	}

	result = Db.Create(&newMessage)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
