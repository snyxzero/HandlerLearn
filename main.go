package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pasha/handler"
	"pasha/repository"
)

func CheckIpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}
		ip := r.RemoteAddr
		fmt.Println("Middleware пропустил запрос от этого ip -", ip)
		if ip == "123.0.0.0:8080" {
			http.Error(w, "Попался педик", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

	pgRepo := repository.NewUserSQLRepository()
	inmemoryRepo := repository.NewUserInMemoryRepository()

	/*	mem := serverMemoryPackage.NewServerMemory()*/
	userHandler := handler.NewUserHandler(inmemoryRepo)
	rout := mux.NewRouter()
	rout.Use(CheckIpMiddleware)

	rout.HandleFunc("/user", userHandler.AddUserHandler).Methods(http.MethodPost)
	rout.HandleFunc("/user/{id:[0-9]+}", userHandler.GetUserHandler).Methods(http.MethodGet)
	rout.HandleFunc("/user/{id:[0-9]+}", userHandler.UpdateUserHandler).Methods(http.MethodPut)
	rout.HandleFunc("/user/{id:[0-9]+}", userHandler.DeleteUserHandler).Methods(http.MethodDelete)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", rout))

	pgRepo.Close()
}
