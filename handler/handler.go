package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	serverMemoryPackage "pasha/memory"
	"pasha/models"
)

type UserClipboard struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Age   int    `json:"age,omitempty"`
}

type UserHandler struct {
	serverMemory *serverMemoryPackage.ServerMemory
}

func NewUserHandler(serverMemory *serverMemoryPackage.ServerMemory) *UserHandler {

	return &UserHandler{
		serverMemory: serverMemory,
	}
}

func printUser(w http.ResponseWriter, user models.User) {
	fmt.Fprintf(w, `Email %d: 
Имя - %s
Возвраст - %d`, user.Email, user.Name, user.Age)
	return
}

func (obj *UserHandler) AddUserHandler(w http.ResponseWriter, r *http.Request) {

	var data UserClipboard
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := models.User{
		Email: data.Email,
		Name:  data.Name,
		Age:   data.Age,
	}

	if obj.serverMemory.FoundUser(user.Email) {
		http.Error(w, "Пользователь с этим Email уже существует", http.StatusConflict)
		return
	}

	obj.serverMemory.AddUser(user)
	fmt.Println("Записали нового пользователя:")
	printUser(w, obj.serverMemory.GetUser(user.Email))
	return
}

func (obj *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]

	if !obj.serverMemory.FoundUser(email) {
		http.Error(w, "Пользователя с этим Email не существует", http.StatusConflict)
		return
	}

	fmt.Println("Данные пользователя по запросу:")
	printUser(w, obj.serverMemory.GetUser(email))
	return
}

func (obj *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]

	if !obj.serverMemory.FoundUser(email) {
		http.Error(w, "Пользователя с этим Email не существует", http.StatusConflict)
		return
	}

	var data UserClipboard
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := models.User{
		Email: email,
		Name:  data.Name,
		Age:   data.Age,
	}

	obj.serverMemory.UpdateUser(user)

	fmt.Fprintf(w, "Пользователь Email = %d изменен:", email)
	printUser(w, obj.serverMemory.GetUser(email))

	return
}

func (obj *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]

	if !obj.serverMemory.FoundUser(email) {
		http.Error(w, "Пользователя с этим Email не существует", http.StatusConflict)
		return
	}

	obj.serverMemory.DeleteUser(email)
	fmt.Fprintln(w, "Пользователь удален")

	return
}
