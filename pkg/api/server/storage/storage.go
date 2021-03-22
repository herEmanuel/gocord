package storage

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/herEmanuel/gocord/pkg/api/models"
	"gorm.io/gorm"
)

//TODO: Debug all of that

func CreateServer(serverVar *models.Server, userID uuid.UUID, name string) error {

	var creator models.User

	Db.First(&creator, "id = ?", userID)

	//Generate invite code
	var inviteCode string
	letters := "abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 8; i++ {
		inviteCode += string(letters[rand.Intn(26)])
	}

	newServer := models.Server{
		Name:       name,
		InviteCode: inviteCode,
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

func DeleteServer(serverID uuid.UUID, ts ...*gorm.DB) error {

	db := Db
	if len(ts) > 0 {
		db = ts[0]
	}

	var server models.Server

	result := db.First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	result = db.Select("Members", "Admins").Delete(&server, "id = ?", server.ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func AddImage(imagePath string, serverID uuid.UUID) error {

	var server models.Server

	result := Db.First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
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

// 	result := Db.First(&server, "id = ?", serverID)

// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return errors.New("This server doesn't exist")
// 	}

// 	return nil
// }

func GetServer(serverVar *models.Server, serverID uuid.UUID) error {

	var server models.Server
	//TODO: check if it works (ordering roles by priority)
	result := Db.
		Preload("Channels", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "permission", "server")
		}).
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "color", "priority", "server").Order("priority DESC")
		}).
		Preload("Roles.Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "avatar")
		}).
		Select("id", "name", "picture", "invite_code").
		First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	*serverVar = server

	return nil
}

func CreateChannel(channelVar *models.Channel, serverID uuid.UUID, name, permission string) error {

	var server models.Server

	result := Db.First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	newChannel := models.Channel{
		Name:       name,
		Permission: permission,
		Server:     server.ID,
	}

	result = Db.Create(&newChannel)
	if result.Error != nil {
		return result.Error
	}

	*channelVar = newChannel

	return nil
}

func DeleteChannel(channelID uuid.UUID) error {

	var channel models.Channel

	result := Db.First(&channel, "id = ?", channelID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This channel doesn't exist")
	}

	Db.Delete(&channel, "id = ?", channel.ID)

	return nil
}

func GetChannelMessages(userID, channelID uuid.UUID) ([]models.Message, error) {

	var channel models.Channel
	var user models.User

	result := Db.Preload("Messages").
		Preload("Messages.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "avatar")
		}).
		First(&channel, "id = ?", channelID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []models.Message{}, errors.New("This channel doesn't exist")
	}

	Db.Preload("Servers").
		First(&user, "id = ?", userID)

	isInServer := false
	for _, server := range user.Servers {
		if server.ID == channel.Server {
			isInServer = true
		}
	}

	if !isInServer {
		return []models.Message{}, errors.New("You are not in this server")
	}

	if channel.Permission == "admin-only" {
		var server models.Server
		Db.Preload("Admins").First(&server, "id = ?", channel.Server)

		isServerAdmin := false
		for _, serverAdmin := range server.Admins {
			if serverAdmin.ID == user.ID {
				isServerAdmin = true
				break
			}
		}
		if !isServerAdmin {
			return []models.Message{}, errors.New("You don't have permission to see the messages of this channel")
		}
	}

	return channel.Messages, nil
}

func SendMessage(creatorID uuid.UUID, channelID uuid.UUID, content, messageType string) (map[string]interface{}, error) {

	var channel models.Channel
	var creator models.User

	result := Db.First(&channel, "id = ?", channelID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("This channel doesn't exist")
	}

	Db.Preload("Servers").
		First(&creator, "id = ?", creatorID)

	isInServer := false
	for _, server := range creator.Servers {
		if server.ID == channel.Server {
			isInServer = true
		}
	}

	if !isInServer {
		return nil, errors.New("You are not in this server")
	}

	if channel.Permission != "public" {
		var server models.Server
		Db.Preload("Admins").First(&server, "id = ?", channel.Server)

		isServerAdmin := false
		for _, serverAdmin := range server.Admins {
			if serverAdmin.ID == creator.ID {
				isServerAdmin = true
				break
			}
		}
		if !isServerAdmin {
			return nil, errors.New("You don't have permission to send a message on this channel")
		}
	}

	newMessage := models.Message{
		Content: content,
		Type:    messageType,
		UserID:  creator.ID,
		Channel: channel.ID,
	}

	result = Db.Create(&newMessage)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"newMessage": newMessage,
		"userName":   creator.Name,
		"userAvatar": creator.Avatar,
	}, nil
}

func DeleteMessage(userID uuid.UUID, messageID uuid.UUID) error {

	var message models.Message

	result := Db.First(&message, "id = ?", messageID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This message doesn't exist")
	}

	if message.UserID != userID {
		return errors.New("You did not send this message")
	}

	result = Db.Delete(&message, "id = ?", message.ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CreateRole(roleVar *models.Role, serverID uuid.UUID, priority uint8, name, color string) error {

	var server models.Server

	result := Db.First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	newRole := models.Role{
		Name:     name,
		Color:    color,
		Priority: priority,
		Server:   server.ID,
	}

	result = Db.Create(&newRole)
	if result.Error != nil {
		return result.Error
	}

	*roleVar = newRole

	return nil
}

func DeleteRole(roleID, serverID uuid.UUID) error {

	var server models.Server

	result := Db.Preload("Roles", "id = ?", roleID).
		First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	if len(server.Roles) == 0 {
		return errors.New("This role doesn't exist")
	}

	result = Db.Select("Users").Delete(&server.Roles[0], "id = ?", server.Roles[0].ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func AddRoleToUser(roleID, userID, serverID uuid.UUID) error {

	var server models.Server
	var user models.User

	result := Db.Preload("Roles", "id = ?", roleID).
		First(&server, "id = ?", serverID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This server doesn't exist")
	}

	if len(server.Roles) == 0 {
		return errors.New("This role doesn't exist")
	}

	Db.First(&user, "id = ?", userID)

	server.Roles[0].Users = append(server.Roles[0].Users, user)

	result = Db.Omit("Users.*").Save(&server.Roles[0])
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func RemoveUser(userID, serverID uuid.UUID) error {

	var user models.User
	var server models.Server

	result := Db.Preload("Servers").Preload("Roles").First(&user, "id = ?", userID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This user doesn't exist")
	}

	isInServer := false
	for _, userServer := range user.Servers {
		if userServer.ID == serverID {
			isInServer = true
			break
		}
	}

	if !isInServer {
		return errors.New("This user is not in the server")
	}

	Db.First(&server, "id = ?", serverID)

	err := Db.Transaction(func(ts *gorm.DB) error {

		err := ts.Model(&server).Association("Members").Delete(&user)
		if err != nil {
			return err
		}

		for _, role := range user.Roles {
			if role.Server == server.ID {
				err = ts.Model(&user).Association("Roles").Delete(&role)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func PromoteToAdmin(userID, serverID uuid.UUID) error {

	var user models.User
	var server models.Server

	result := Db.Preload("Servers").First(&user, "id = ?", userID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("This user doesn't exist")
	}

	isInServer := false
	for _, userServer := range user.Servers {
		if userServer.ID == serverID {
			isInServer = true
			break
		}
	}

	if !isInServer {
		return errors.New("This user is not in the server")
	}

	Db.First(&server, "id = ?", serverID)

	server.Admins = append(server.Admins, user)

	result = Db.Omit("Admins.*").Save(&server)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
