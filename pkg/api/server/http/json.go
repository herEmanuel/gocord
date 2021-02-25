package http

import "github.com/google/uuid"

type CreateServerJSON struct {
	Name string `json:"name"`
}

type CreateChannelJSON struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

type SendMessageJSON struct {
	Content     string `json:"content"`
	MessageType string `json:"messageType"`
}

type CreateRoleJSON struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	Priority uint8  `json:"priority"`
}

//TODO: check if it can convert string to uuid
type AddRoleToUserJSON struct {
	RoleID   uuid.UUID `json:"roleID"`
	UserID   uuid.UUID `json:"userID"`
	ServerID uuid.UUID `json:"serverID"`
}

type RemoveUserJSON struct {
	UserID   uuid.UUID `json:"userID"`
	ServerID uuid.UUID `json:"serverID"`
}

type PromoteToAdminJSON struct {
	UserID   uuid.UUID `json:"userID"`
	ServerID uuid.UUID `json:"serverID"`
}
