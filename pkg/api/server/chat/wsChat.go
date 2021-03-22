package chat

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/api/models"
	"github.com/herEmanuel/gocord/pkg/api/server/storage"
	"gorm.io/gorm"
)

var activeUsers map[uuid.UUID]*websocket.Conn
var channels map[uuid.UUID][]uuid.UUID //only 1 key: the channel id, which maps to an array of users

const (
	CLOSE_MESSAGE   = "close"
	CONNECT_MESSAGE = "connect"
)

func Init() {
	activeUsers = make(map[uuid.UUID]*websocket.Conn)
	channels = make(map[uuid.UUID][]uuid.UUID)
}

//TODO: improve that code and see if theres a better way to do all of that

func retry(ws *websocket.Conn) error {
	times := 0

	err := ws.WriteMessage(websocket.TextMessage, []byte("retry"))
	for err != nil {
		times++

		err = ws.WriteMessage(websocket.TextMessage, []byte("retry"))

		if times == 3 {
			return errors.New("Failed")
		}
	}

	return nil
}

func handshake(ws *websocket.Conn) (uuid.UUID, error) {

	var user models.User
	var channel models.Channel
	connMessage := new(ConnectionMessage)
	userID := ws.Locals("userID").(uuid.UUID)

	err := ws.ReadJSON(&connMessage)
	if err != nil {
		retryErr := retry(ws)
		if retryErr != nil {
			return uuid.Nil, errors.New("Failed")
		}

		err = ws.ReadJSON(&connMessage)
		if err != nil {
			return uuid.Nil, errors.New("Failed")
		}
	}

	fmt.Printf("received a connection request from %v\n", userID)

	result := storage.Db.First(&channel, "id = ?", connMessage.ChannelID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return uuid.Nil, errors.New("This channel doesn't exist")
	}
	//check
	storage.Db.Preload("Servers", "id = ?", channel.Server).First(&user, "id = ?", userID)

	if len(user.Servers) == 0 {
		return uuid.Nil, errors.New("You are not in this server")
	}

	channels[channel.ID] = append(channels[channel.ID], userID)

	return channel.ID, nil
}

func WSConn(ws *websocket.Conn) {

	userID := ws.Locals("userID").(uuid.UUID)

	_, msg, err := ws.ReadMessage()
	if err != nil {
		fmt.Println(err)

		retryErr := retry(ws)
		if retryErr != nil {
			ws.Close()
			return
		}

		_, msg, err = ws.ReadMessage()
		if err != nil {
			ws.Close()
			return
		}
	}
	if string(msg) != CONNECT_MESSAGE {
		fmt.Println("wrong message")
		ws.Close()
	}

	channelID, err := handshake(ws)
	if err != nil {
		ws.Close()
		return
	}

	activeUsers[userID] = ws

	ws.SetCloseHandler(func(code int, text string) error {

		if code == 1001 {
			delete(activeUsers, userID)

			//remove user from channels list
			for i, user := range channels[channelID] {
				if user == userID {
					channels[channelID][i] = channels[channelID][len(channels[channelID])-1]
					channels[channelID] = channels[channelID][:len(channels[channelID])-1]
				}
			}
		}

		return nil
	})

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)

			retryErr := retry(ws)
			if retryErr != nil {
				ws.Close()
				return
			}

			_, msg, err = ws.ReadMessage()
			if err != nil {
				ws.Close()
				return
			}
		}

		switch string(msg) {
		case CLOSE_MESSAGE:

			delete(activeUsers, userID)
			//remove user from channels list
			for i, user := range channels[channelID] {
				if user == userID {
					channels[channelID][i] = channels[channelID][len(channels[channelID])-1]
					channels[channelID] = channels[channelID][:len(channels[channelID])-1]
				}
			}

			ws.Close()
		case CONNECT_MESSAGE:

			//remove user from channels list
			for i, user := range channels[channelID] {
				fmt.Println(user)
				if user == userID {
					channels[channelID][i] = channels[channelID][len(channels[channelID])-1]
					channels[channelID] = channels[channelID][:len(channels[channelID])-1]
				}
			}

			channelID, err = handshake(ws)
			if err != nil {
				ws.Close()
				return
			}
		}
	}
}

func sendMessage(user uuid.UUID, message NormalMessage) {
	ws := activeUsers[user]
	log.Println("here")

	err := ws.WriteJSON(message)
	if err != nil {
		log.Println(err)
	}
}

func TriggerSendMessage(channelID, messageID, userID uuid.UUID, userName, userAvatar, content, messageType string) {

	if len(channels[channelID]) > 0 {

		message := NormalMessage{
			Content:     content,
			MessageType: messageType,
			MessageID:   messageID.String(),
			UserID:      userID.String(),
			UserAvatar:  userAvatar,
			UserName:    userName,
		}

		for _, user := range channels[channelID] {
			go sendMessage(user, message)
		}
	}

}
