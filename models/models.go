package models

import "time"

//easyjson:json
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

//easyjson:json
type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Users     []string  `json:"users"`
	CreatedAt time.Time `json:"created_at"`
}

//easyjson:json
type Message struct {
	ID        int       `json:"id"`
	ChatID    string    `json:"chat"`
	AuthorID  string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
