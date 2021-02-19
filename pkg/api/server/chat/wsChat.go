package chat

import (
	"fmt"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

var activeUsers map[uuid.UUID]*websocket.Conn
var channels map[uuid.UUID][]uuid.UUID //only 1 key: the channel id, which maps to an array of users

func Init() {
	activeUsers = make(map[uuid.UUID]*websocket.Conn)
	channels = make(map[uuid.UUID][]uuid.UUID)
}

//TODO: improve that code and see if theres a better way to do all of that

func WSConn(ws *websocket.Conn) {

	connMessage := new(ConnectionMessage)
	err := ws.ReadJSON(&connMessage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("received a connection request from " + connMessage.UserID)
	channelID, _ := uuid.Parse(connMessage.ChannelID)
	userID, _ := uuid.Parse(connMessage.UserID)

	activeUsers[userID] = ws

	channels[channelID] = append(channels[channelID], userID)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			ws.Close()
		}

		if string(msg) != "close" {
			continue
		}

		delete(activeUsers, userID)

		//remove user from channels list
		for i, user := range channels[channelID] {
			if user == userID {
				channels[channelID][i-1] = channels[channelID][len(channels[channelID])-1]
				channels[channelID] = channels[channelID][:len(channels[channelID])-1]
			}
		}

		ws.Close()
	}

}

func TriggerSendMessage(channelID, messageID, userID uuid.UUID, userName, userAvatar, content, messageType string) {

	if len(channels[channelID]) > 0 {
		for _, user := range channels[channelID] {
			ws := activeUsers[user]
			fmt.Println("here")
			message := NormalMessage{
				Content:     content,
				MessageType: messageType,
				MessageID:   messageID.String(),
				UserID:      userID.String(),
				UserAvatar:  userAvatar,
				UserName:    userName,
			}

			err := ws.WriteJSON(message)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
