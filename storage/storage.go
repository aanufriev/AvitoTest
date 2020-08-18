package storage

import "github.com/aanufriev/AvitoTest/models"

type StorageInterface interface {
	SaveUser(user *models.User) (int, error)
	SaveChat(chat *models.Chat) (int, error)
	SaveMessage(message *models.Message) (int, error)
	GetChats(userID string) ([]models.Chat, error)
	GetMessages(chatID string) ([]models.Message, error)
}
