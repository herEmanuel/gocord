package gateway

import (
	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/api/models"
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

func GetServer(serverID uuid.UUID) (models.Server, error) {

	var server models.Server

	err := storage.GetServer(&server, serverID)
	if err != nil {
		return models.Server{}, nil
	}

	return server, nil
}

func GetChannelMessages(userID uuid.UUID, channelID uuid.UUID) ([]models.Message, error) {

	messagesVar, err := storage.GetChannelMessages(userID, channelID)
	if err != nil {
		return []models.Message{}, nil
	}

	return messagesVar, nil
}

func SendMessage(creatorID uuid.UUID, channelID uuid.UUID, content, messageType string) error {

	err := storage.SendMessage(creatorID, channelID, content, messageType)
	if err != nil {
		return err
	}

	return nil
}
