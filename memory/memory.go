package serverMemoryPackage

import (
	"pasha/models"
	"sync"
)

type ServerMemory struct {
	user map[int]models.User
	mu   sync.Mutex
}

func NewServerMemory() *ServerMemory {
	return &ServerMemory{
		user: make(map[int]models.User),
	}
}

func (obj *ServerMemory) FoundUser(id int) (ok bool) {
	_, ok = obj.user[id]
	return ok
}

func (obj *ServerMemory) AddUser(user models.User) {
	obj.mu.Lock()
	obj.user[user.ID] = user
	obj.mu.Unlock()
	return
}

func (obj *ServerMemory) GetUser(id int) (user models.User) {
	user = obj.user[id]
	return user
}

func (obj *ServerMemory) UpdateUser(user models.User) {
	obj.mu.Lock()
	obj.user[user.ID] = user
	obj.mu.Unlock()
	return
}

func (obj *ServerMemory) DeleteUser(id int) {
	obj.mu.Lock()
	delete(obj.user, id)
	obj.mu.Unlock()
	return
}
