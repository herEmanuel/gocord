package http

import "github.com/google/uuid"

type CreateServerJSON struct {
	Name string `json:"name"`
}

type GetServerJSON struct {
	ServerID uuid.UUID `json:"serverId"`
}

type SendMessageJSON struct {
	Content     string `json:"content"`
	MessageType string `json:"messageType"`
}
