package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pasha/memory"
	"strconv"
)

/*type Paulus struct {
	Age      uint `json:"сколько лет"`
	Baldness uint `json:"какая лысовость"`
	Money    uint `json:"сколько денег"`
}*/

/*
var memory []Paulus

const headerTruth = "Truth"
*/
//var memory1 *memory.ServerMemory

func pashaCreate(mem *memory.ServerMemory) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data memory.UserClipboard
		err := json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//if r.Header.Get(headerTruth) == "True" {
		//	data = Paulus{data.Age, 100, 99999999999999}
		//}

		err = mem.UserAdd(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(data)
		return
	}

}

func pashaRead(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := memory1.UserGet(uint(id))
	fmt.Fprintf(w, `Юзер номер %d: 
Id: %d
Имя: %s
Сколько лет - %d
Он гей? - %t`, id, user.Id, user.Name, user.Age, user.Gay)

	return
}

func pashaUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data memory.UserClipboard
	err = json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memory1.UserUpdate(data)

	fmt.Fprintf(w, `Юзер номер %d изменен: 
Id: %d
Имя: %s
Сколько лет - %d
Он гей? - %t`, id, data.Id, data.Name, data.Age, data.Gay)

	return
}

func pashaDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	memory1.UserDel(uint(id))
	fmt.Fprintln(w, "Удален")

	return
}

func CheckIpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			next.ServeHTTP(w, r)
		}
		ip := r.RemoteAddr
		fmt.Println("Запрос от этого ip -", ip)
		if ip == "123.0.0.0:8080" {
			http.Error(w, "Попался педик", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mem := memory.NewServerMemory()
	rout := mux.NewRouter()
	rout.Use(CheckIpMiddleware)

	rout.HandleFunc("/pasha", pashaCreate(mem)).Methods(http.MethodPost)
	rout.HandleFunc("/pasha/{id:[0-9]+}", pashaRead).Methods(http.MethodGet)
	rout.HandleFunc("/pasha/{id:[0-9]+}", pashaUpdate).Methods("PUT")
	rout.HandleFunc("/pasha/{id:[0-9]+}", pashaDelete).Methods("DELETE")

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", rout))

}
