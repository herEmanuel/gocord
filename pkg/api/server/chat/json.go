package chat

type ConnectionMessage struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
}

type NormalMessage struct {
	Content     string `json:"content"`
	MessageType string `json:"messageType"`
	MessageID   string `json:"messageID"`
	UserID      string `json:"userID"`
	UserAvatar  string `json:"userAvatar"`
	UserName    string `json:"userName"`
}
