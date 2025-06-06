package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pasha/handler"
	"pasha/repository"
	"syscall"
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

	ctx := context.Background()

	pgRepo := repository.NewUserSQLRepository(ctx)
	defer pgRepo.Close(ctx)
	inmemoryRepo := repository.NewUserInMemoryRepository(ctx)

	userHandler := handler.NewUserHandler(ctx, inmemoryRepo)
	rout := mux.NewRouter()
	rout.Use(CheckIpMiddleware)

	rout.HandleFunc("/user", userHandler.AddUserHandler).Methods(http.MethodPost)
	rout.HandleFunc("/user/{id:[0-9]+}", userHandler.GetUserHandler).Methods(http.MethodGet)
	rout.HandleFunc("/user/{id:[0-9]+}", userHandler.UpdateUserHandler).Methods(http.MethodPut)
	rout.HandleFunc("/user/{id:[0-9]+}", userHandler.DeleteUserHandler).Methods(http.MethodDelete)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: rout}
	go func() {
		fmt.Println("Сервер запущен на http://localhost:8080")
		log.Fatal(srv.ListenAndServe())
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Ожидаем сигнала завершения
	<-stop

	// Graceful Shutdown

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
