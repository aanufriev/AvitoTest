package main

import (
	"log"
	"net/http"

	"github.com/aanufriev/AvitoTest/storage"
	_ "github.com/lib/pq"
)

type Handler struct{}

func (h Handler) AddUser(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) AddChat(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) AddMessage(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) GetChats(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {

}

func main() {
	storage := &storage.PostgresStorage{}
	err := storage.Open("host=127.0.0.1 user=testuser password=test_password dbname=avito sslmode=disable")
	if err != nil {
		log.Fatal("can't open database connection: ", err)
	}

	err = storage.InitDatabase()
	if err != nil {
		log.Fatal("can't init database: ", err)
	}

	handler := &Handler{}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/add", handler.AddUser)
	mux.HandleFunc("/chats/add", handler.AddChat)
	mux.HandleFunc("/messages/add", handler.AddMessage)
	mux.HandleFunc("/chats/get", handler.GetChats)
	mux.HandleFunc("/messages/get", handler.GetMessages)

	http.ListenAndServe(":9000", mux)
}
