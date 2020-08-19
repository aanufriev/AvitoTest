package main

import (
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
