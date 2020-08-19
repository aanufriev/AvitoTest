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

func (ps PostgresStorage) SaveChat(chat *models.Chat) (int, error) {
	if len(chat.Users) < 2 {
		return 0, fmt.Errorf("SaveChat error: Expect users count: > 2\nGot users count: %v", len(chat.Users))
	}

	for _, userID := range chat.Users {
		if !ps.checkUser(userID) {
			return 0, fmt.Errorf("SaveChat error: User with id %v doesn`t exist", userID)
		}
	}

	var lastID int

	err := ps.db.QueryRow(
		"INSERT INTO chats (name, created_at) VALUES ($1, $2) RETURNING id",
		chat.Name, chat.CreatedAt,
	).Scan(&lastID)

	if err != nil {
		return 0, fmt.Errorf("SaveChat error: %s with chat: %v", err, chat)
	}
	chat.ID = lastID

	for _, userID := range chat.Users {
		_, err = ps.db.Exec(
			"INSERT INTO userchat (user_id, chat_id) VALUES ($1, $2)",
			userID, chat.ID,
		)
		if err != nil {
			return 0, fmt.Errorf("SaveChat error: %s with user_id: %s", err, userID)
		}
	}

	return lastID, nil
}

func (ps PostgresStorage) checkUser(userID string) bool {
	user := models.User{}

	row := ps.db.QueryRow(
		"SELECT * FROM users WHERE id=$1",
		userID,
	)
	switch err := row.Scan(&user.ID, &user.Username, &user.CreatedAt); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		panic(err)
	}
}

func (ps PostgresStorage) SaveMessage(msg *models.Message) (int, error) {
	if !ps.checkUser(msg.AuthorID) {
		return 0, fmt.Errorf("SaveMessage error: User with id %v doesn`t exist", msg.AuthorID)
	}

	if !ps.checkChat(msg.ChatID) {
		return 0, fmt.Errorf("SaveMessage error: Chat with id %v doesn`t exist", msg.ChatID)
	}

	var lastID int

	err := ps.db.QueryRow(
		"INSERT INTO messages (chat_id, user_id, msg_text, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		msg.ChatID, msg.AuthorID, msg.Text, msg.CreatedAt,
	).Scan(&lastID)

	if err != nil {
		return 0, fmt.Errorf("SaveMessage error: %s with message: %v", err, msg)
	}
	msg.ID = lastID

	return lastID, err
}

func (ps PostgresStorage) checkChat(chatID string) bool {
	chat := models.Chat{}

	row := ps.db.QueryRow(
		"SELECT * FROM chats WHERE id=$1",
		chatID,
	)
	switch err := row.Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		panic(err)
	}
}
