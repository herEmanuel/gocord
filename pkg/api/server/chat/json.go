package chat

import "github.com/google/uuid"

type ConnectionMessage struct {
	ChannelID uuid.UUID `json:"channelId"`
}

type NormalMessage struct {
	Content     string `json:"content"`
	MessageType string `json:"messageType"`
	MessageID   string `json:"messageID"`
	UserID      string `json:"userID"`
	UserAvatar  string `json:"userAvatar"`
	UserName    string `json:"userName"`
}
