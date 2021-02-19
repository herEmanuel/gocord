package gateway

import (
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
		return []models.Message{}, nil
	}

	return messagesVar, nil
}

func SendMessage(creatorID uuid.UUID, channelID uuid.UUID, content, messageType string) error {

	message, err := storage.SendMessage(creatorID, channelID, content, messageType)
	if err != nil {
		return err
	}

	//trigger an event to broadcast the message via websockets
	go chat.TriggerSendMessage(channelID, message.ID, message.User.ID, message.User.Name, message.User.Avatar, content, messageType)

	return nil
}

func DeleteMessage(userID uuid.UUID, messageID uuid.UUID) error {

	err := storage.DeleteMessage(userID, messageID)
	if err != nil {
		return err
	}

	return nil
}
