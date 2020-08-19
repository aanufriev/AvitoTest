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

type Handler struct {
	storage storage.StorageInterface
}

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

func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
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
