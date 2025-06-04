package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pasha/models"
	"strconv"
)

type Repository interface {
	AddUser(user models.User) error
	GetUser(id int) (*models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(id int) error
}

const (
	errDataFormat = "Неверный формат данных"
	errInvalidID  = "Некорректный id"
)

type UserClipboard struct {
	ID    int    `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Age   int    `json:"age,omitempty"`
}

type UserHandler struct {
	repository Repository
}

func NewUserHandler(repository Repository) *UserHandler {

	return &UserHandler{
		repository: repository,
	}
}

func (obj *UserHandler) AddUserHandler(w http.ResponseWriter, r *http.Request) {

	var data UserClipboard
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		http.Error(w, errDataFormat, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := models.User{

		Email: data.Email,
		Name:  data.Name,
		Age:   data.Age,
	}

	err = obj.repository.AddUser(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Не удалось добавить пользователя", http.StatusBadRequest)
	}
	fmt.Fprintln(w, "Пользователь добавлен")
	return
}

func (obj *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	user, err := obj.repository.GetUser(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Не нашли пользователя", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println("Error writing response:", err)
		http.Error(w, "Ошибка сервера при ответе", http.StatusInternalServerError)
		return
	}

}

func (obj *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	var data UserClipboard
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		http.Error(w, errDataFormat, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := models.User{
		ID:    id,
		Email: data.Email,
		Name:  data.Name,
		Age:   data.Age,
	}

	err = obj.repository.UpdateUser(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Не получилось обновить данные", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Пользователь ID = %d изменен\n", id)
	return
}

func (obj *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	err = obj.repository.DeleteUser(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Не получилось удалить пользователя", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "Пользователь удален")

	return
}
