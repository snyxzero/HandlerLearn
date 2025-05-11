package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	serverMemoryPackage "pasha/memory"
	"pasha/models"
	"strconv"
)

type UserClipboard struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
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
	fmt.Fprintf(w, `ID %d: 
Имя - %s
Возвраст - %d`, user.ID, user.Name, user.Age)
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
		ID:   data.ID,
		Name: data.Name,
		Age:  data.Age,
	}

	if obj.serverMemory.FoundUser(user.ID) {
		http.Error(w, "Пользователь с этим ID уже существует", http.StatusConflict)
		return
	}

	obj.serverMemory.AddUser(user)
	fmt.Println("Записали нового пользователя:")
	printUser(w, obj.serverMemory.GetUser(user.ID))
	return
}

func (obj *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !obj.serverMemory.FoundUser(id) {
		http.Error(w, "Пользователя с этим ID не существует", http.StatusConflict)
		return
	}

	fmt.Println("Данные пользователя по запросу:")
	printUser(w, obj.serverMemory.GetUser(id))
	return
}

func (obj *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !obj.serverMemory.FoundUser(id) {
		http.Error(w, "Пользователя с этим ID не существует", http.StatusConflict)
		return
	}

	var data UserClipboard
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := models.User{
		ID:   id,
		Name: data.Name,
		Age:  data.Age,
	}

	obj.serverMemory.UpdateUser(user)

	fmt.Fprintf(w, "Пользователь ID = %d изменен:", id)
	printUser(w, obj.serverMemory.GetUser(id))

	return
}

func (obj *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !obj.serverMemory.FoundUser(id) {
		http.Error(w, "Пользователя с этим ID не существует", http.StatusConflict)
		return
	}

	obj.serverMemory.DeleteUser(id)
	fmt.Fprintln(w, "Пользователь удален")

	return
}
