package memory

import (
	"fmt"
	"net/http"
	"sync"
)

type UserClipboard struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
	Gay  bool   `json:"gay"`
}

type User struct {
	id   uint
	name string
	age  uint
	gay  bool
}

type ServerMemory struct {
	user map[uint]User
	Mu   sync.Mutex
	w    http.ResponseWriter
}

func NewServerMemory() *ServerMemory {
	return &ServerMemory{
		user: make(map[uint]User),
	}
}

func (memory *ServerMemory) UserGet(id uint) (*UserClipboard, error) {
	if _, ok := memory.user[id]; !ok {

		return nil, fmt.Errorf("User not found")
	}
	return &UserClipboard{
		Id:   id,
		Name: memory.user[id].name,
		Age:  memory.user[id].age,
		Gay:  memory.user[id].gay,
	}, nil
}

func (memory *ServerMemory) UserAdd(u UserClipboard) error {
	memory.Mu.Lock()
	defer memory.Mu.Unlock()
	if _, ok := memory.user[u.Id]; ok {
		return fmt.Errorf("User already exist")
	}
	memory.user[u.Id] = User{
		id:   u.Id,
		name: u.Name,
		age:  u.Age,
		gay:  u.Gay,
	}
	return nil
}

func (memory *ServerMemory) UserUpdate(u UserClipboard) {
	memory.Mu.Lock()
	defer memory.Mu.Unlock()
	if _, ok := memory.user[u.Id]; !ok {
		http.Error(memory.w, "Этого id не существует", http.StatusConflict)
		return
	}
	memory.user[u.Id] = User{
		id:   u.Id,
		name: u.Name,
		age:  u.Age,
		gay:  u.Gay,
	}
}

func (memory *ServerMemory) UserDel(id uint) {
	memory.Mu.Lock()
	defer memory.Mu.Unlock()
	if _, ok := memory.user[id]; !ok {
		http.Error(memory.w, "Этого id не существует", http.StatusConflict)
		return
	}
	delete(memory.user, id)
}
