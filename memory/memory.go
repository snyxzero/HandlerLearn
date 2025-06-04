package serverMemoryPackage

import (
	"pasha/models"
	"sync"
)

type ServerMemory struct {
	user map[string]models.User
	mu   sync.Mutex
}

func NewServerMemory() *ServerMemory {
	return &ServerMemory{
		user: make(map[string]models.User),
	}
}

func (obj *ServerMemory) FoundUser(email string) (ok bool) {
	_, ok = obj.user[email]
	return ok
}

func (obj *ServerMemory) AddUser(user models.User) {
	obj.mu.Lock()
	obj.user[user.Email] = user
	obj.mu.Unlock()
	return
}

func (obj *ServerMemory) GetUser(email string) (user models.User) {
	user = obj.user[email]
	return user
}

func (obj *ServerMemory) UpdateUser(user models.User) {
	obj.mu.Lock()
	obj.user[user.Email] = user
	obj.mu.Unlock()
	return
}

func (obj *ServerMemory) DeleteUser(email string) {
	obj.mu.Lock()
	delete(obj.user, email)
	obj.mu.Unlock()
	return
}
