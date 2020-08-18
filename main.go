package main

import (
	"net/http"
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
	handler := &Handler{}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/add", handler.AddUser)
	mux.HandleFunc("/chats/add", handler.AddChat)
	mux.HandleFunc("/messages/add", handler.AddMessage)
	mux.HandleFunc("/chats/get", handler.GetChats)
	mux.HandleFunc("/messages/get", handler.GetMessages)

	http.ListenAndServe(":9000", mux)
}
