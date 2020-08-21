package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aanufriev/AvitoTest/models"
	"github.com/aanufriev/AvitoTest/storage"
	_ "github.com/lib/pq"
)

func writeError(w http.ResponseWriter, status int, err error, u *url.URL) {
	fmt.Printf("error happend: %+v\nurl: %s\n", err, u)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	errorJSON := fmt.Sprintf(`{"error":"%s"}`, err)
	w.Write([]byte(errorJSON))
}

func idAsJSON(id int) []byte {
	return []byte(fmt.Sprintf(`{"id":"%v"}`, id))
}

// Handler processes user requests
type Handler struct {
	storage storage.StorageInterface
}

// AddUser creates a new user with unique username and returns id
func (h Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}
	defer r.Body.Close()

	user := &models.User{
		CreatedAt: time.Now(),
	}
	user.UnmarshalJSON(body)

	id, err := h.storage.SaveUser(user)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}

	w.WriteHeader(http.StatusOK)
	idJSON := idAsJSON(id)
	w.Write(idJSON)
}

// AddChat creates a new chat with 2 or more users. Has a unique name and returns id
func (h Handler) AddChat(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}
	defer r.Body.Close()

	chat := &models.Chat{
		CreatedAt: time.Now(),
	}
	chat.UnmarshalJSON(body)

	id, err := h.storage.SaveChat(chat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}

	w.WriteHeader(http.StatusOK)
	idJSON := idAsJSON(id)
	w.Write(idJSON)
}

// AddMessage creates a new message in chat and returns id
func (h Handler) AddMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}
	defer r.Body.Close()

	msg := &models.Message{
		CreatedAt: time.Now(),
	}
	msg.UnmarshalJSON(body)

	id, err := h.storage.SaveMessage(msg)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}

	w.WriteHeader(http.StatusOK)
	idJSON := idAsJSON(id)
	w.Write(idJSON)
}

// GetChats returns all chats the user has.
// Sorted by CreatedAt of the last message in chat (from the latest to earlies)
func (h Handler) GetChats(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if err != nil {
			writeError(w, http.StatusBadRequest, err, r.URL)
			return
		}
	}
	defer r.Body.Close()

	var userID map[string]interface{}

	err = json.Unmarshal(body, &userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}
	id, ok := userID["user"].(string)
	if !ok {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}

	chats, err := h.storage.GetChats(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}

	response := make([]byte, 0)
	for _, chat := range chats {
		chatJSON, err := chat.MarshalJSON()
		if err != nil {
			writeError(w, http.StatusBadRequest, err, r.URL)
			return
		}
		response = append(response, chatJSON...)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetMessages returns all messages from chat
// Sorted by CreatedAt of the message in chat (from the earlies to latest)
func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}
	defer r.Body.Close()

	var chatID map[string]interface{}

	err = json.Unmarshal(body, &chatID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}
	id, ok := chatID["chat"].(string)
	if !ok {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}

	msgs, err := h.storage.GetMessages(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, r.URL)
		return
	}

	response := make([]byte, 0)
	for _, msg := range msgs {
		msgJSON, err := msg.MarshalJSON()
		if err != nil {
			writeError(w, http.StatusBadRequest, err, r.URL)
			return
		}
		response = append(response, msgJSON...)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	storage := &storage.PostgresStorage{}
	err := storage.Open("host=database user=testuser password=test_password dbname=avito sslmode=disable")
	if err != nil {
		log.Fatal("can't open database connection: ", err)
	}

	err = storage.InitDatabase()
	if err != nil {
		log.Fatal("can't init database: ", err)
	}

	handler := &Handler{
		storage,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/add", handler.AddUser)
	mux.HandleFunc("/chats/add", handler.AddChat)
	mux.HandleFunc("/messages/add", handler.AddMessage)
	mux.HandleFunc("/chats/get", handler.GetChats)
	mux.HandleFunc("/messages/get", handler.GetMessages)

	http.ListenAndServe(":9000", mux)
}
