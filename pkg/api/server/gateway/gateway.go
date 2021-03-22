package gateway

import (
	"errors"

	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/api/models"
	"github.com/herEmanuel/gocord/pkg/api/server/chat"
	"github.com/herEmanuel/gocord/pkg/api/server/storage"
)

func CreateServer(userID uuid.UUID, name string) (models.Server, error) {

	var newServer models.Server

	err := storage.CreateServer(&newServer, userID, name)
	if err != nil {
		return models.Server{}, err
	}

	return newServer, nil
}

func DeleteServer(serverID uuid.UUID) error {

	err := storage.DeleteServer(serverID)
	if err != nil {
		return err
	}

	return nil
}

func AddImage(imagePath string, serverID uuid.UUID) error {

	err := storage.AddImage(imagePath, serverID)
	if err != nil {
		return err
	}

	return nil
}

func GetServer(serverID uuid.UUID) (models.Server, error) {

	var server models.Server

	err := storage.GetServer(&server, serverID)
	if err != nil {
		return models.Server{}, nil
	}

	return server, nil
}

func CreateChannel(serverID uuid.UUID, name, permission string) (models.Channel, error) {

	var channel models.Channel

	err := storage.CreateChannel(&channel, serverID, name, permission)
	if err != nil {
		return models.Channel{}, err
	}

	return channel, nil
}

func DeleteChannel(channelID uuid.UUID) error {

	err := storage.DeleteChannel(channelID)
	if err != nil {
		return err
	}

	return nil
}

func GetChannelMessages(userID uuid.UUID, channelID uuid.UUID) ([]models.Message, error) {

	messagesVar, err := storage.GetChannelMessages(userID, channelID)
	if err != nil {
		return []models.Message{}, err
	}

	return messagesVar, nil
}

func SendMessage(creatorID uuid.UUID, channelID uuid.UUID, content, messageType string) (models.Message, error) {

	message, err := storage.SendMessage(creatorID, channelID, content, messageType)
	if err != nil {
		return models.Message{}, err
	}

	newMessage := message["newMessage"].(models.Message)

	//trigger an event to broadcast the message via websockets
	go chat.TriggerSendMessage(channelID, newMessage.ID, newMessage.UserID, message["userName"].(string), message["userAvatar"].(string), content, messageType)

	return newMessage, nil
}

func DeleteMessage(userID uuid.UUID, messageID uuid.UUID) error {

	err := storage.DeleteMessage(userID, messageID)
	if err != nil {
		return err
	}

	return nil
}

func CreateRole(serverID uuid.UUID, priority uint8, name, color string) (models.Role, error) {

	var role models.Role

	err := storage.CreateRole(&role, serverID, priority, name, color)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func DeleteRole(roleID, serverID uuid.UUID) error {

	err := storage.DeleteRole(roleID, serverID)
	if err != nil {
		return err
	}

	return nil
}

func AddRoleToUser(roleID, userID, serverID uuid.UUID) error {

	err := storage.AddRoleToUser(roleID, userID, serverID)
	if err != nil {
		return err
	}

	return nil
}

func RemoveUser(userID, adminID, serverID uuid.UUID) error {

	if userID == adminID {
		return errors.New("You can not remove yourself from the server")
	}

	err := storage.RemoveUser(userID, serverID)
	if err != nil {
		return err
	}

	return nil
}

func PromoteToAdmin(userID, serverID uuid.UUID) error {

	err := storage.PromoteToAdmin(userID, serverID)
	if err != nil {
		return err
	}

	return nil
}
