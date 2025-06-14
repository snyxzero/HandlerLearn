package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pasha/models"
	"pasha/repository"
	"pasha/validators"
	"strconv"
)

type Repository interface {
	AddUser(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int) error
}

const (
	errDataFormat = "Неверный формат данных"
	errInvalidID  = "Некорректный id"
)

type UserClipboard struct {
	ID    int    `json:"id,omitempty" validate:"omitempty"`
	Email string `json:"email" validate:"required,email,gmail_email"`
	Name  string `json:"name" validate:"required,russian_name"`
	Age   int    `json:"age" validate:"gte=0,lte=130"`
}

type UserHandler struct {
	repository Repository
	ctx        context.Context
	validator  *validators.ValidatorForUser
}

func NewUserHandler(ctx context.Context, repository Repository, validator *validators.ValidatorForUser) *UserHandler {

	return &UserHandler{
		repository: repository,
		ctx:        ctx,
		validator:  validator,
	}
}

func writeResponse(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
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

func (obj *UserHandler) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	var data UserClipboard
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		http.Error(w, errDataFormat, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = obj.validator.Struct(data)
	if err != nil {
		log.Println(err.(validator.ValidationErrors))
		http.Error(w, errDataFormat, http.StatusBadRequest)
		return
	}
	user := models.User{

		Email: data.Email,
		Name:  data.Name,
		Age:   data.Age,
	}

	user.ID, err = obj.repository.AddUser(obj.ctx, user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Не удалось добавить пользователя", http.StatusInternalServerError)
	}

	userWriteRes, err := obj.repository.GetUser(obj.ctx, user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка чтения добавленной записи", http.StatusInternalServerError)
		return
	}
	writeResponse(w, *userWriteRes)
	return
}

func (obj *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	user, err := obj.repository.GetUser(obj.ctx, id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, "Не нашли пользователя", http.StatusNotFound)
			return
		}
		http.Error(w, "Что то пошло не так", http.StatusInternalServerError)
		return
	}

	writeResponse(w, *user)

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

	err = obj.validator.Struct(data)
	if err != nil {
		log.Println(err.(validator.ValidationErrors))
		http.Error(w, errDataFormat, http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:    id,
		Email: data.Email,
		Name:  data.Name,
		Age:   data.Age,
	}

	err = obj.repository.UpdateUser(obj.ctx, user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Не получилось обновить данные", http.StatusInternalServerError)
		return
	}

	userWriteRes, err := obj.repository.GetUser(obj.ctx, user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка чтения обновленной записи", http.StatusInternalServerError)
		return
	}
	writeResponse(w, *userWriteRes)

	return
}

func (obj *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, errInvalidID, http.StatusBadRequest)
		return
	}

	err = obj.repository.DeleteUser(obj.ctx, id)
	if err != nil {
		log.Println(err)
		writeResponse(w, map[string]string{"success": "false"})
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}
	writeResponse(w, map[string]string{"success": "true"})

	return
}
