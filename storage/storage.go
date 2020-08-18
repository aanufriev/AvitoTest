package storage

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/aanufriev/AvitoTest/models"
)

const (
	postgres = "postgres"
)

type StorageInterface interface {
	SaveUser(user *models.User) (int, error)
	SaveChat(chat *models.Chat) (int, error)
	SaveMessage(message *models.Message) (int, error)
	GetChats(userID string) ([]models.Chat, error)
	GetMessages(chatID string) ([]models.Message, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func (ps *PostgresStorage) Open(dataSourceName string) error {
	db, err := sql.Open(postgres, dataSourceName)
	if err != nil {
		return err
	}
	ps.db = db

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (ps PostgresStorage) InitDatabase() error {
	file, err := ioutil.ReadFile("./storage/init.sql")
	if err != nil {
		return err
	}

	_, err = ps.db.Query(string(file))
	if err != nil {
		return err
	}
	return nil
}

func (ps PostgresStorage) SaveUser(user *models.User) (int, error) {
	var lastID int

	err := ps.db.QueryRow(
		"INSERT INTO users (username, created_at) VALUES ($1, $2) RETURNING id",
		user.Username, user.CreatedAt,
	).Scan(&lastID)

	if err != nil {
		return 0, fmt.Errorf("SaveUser error: %s with user: %v", err, user)
	}
	user.ID = lastID

	return lastID, nil
}
