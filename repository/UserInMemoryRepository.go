package repository

import (
	"fmt"
	"pasha/models"
	"sync"
)

type InMemoryRepository struct {
	user  map[int]models.User
	count int
	mu    sync.Mutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		user:  make(map[int]models.User),
		count: 1,
	}
}

func (obj *InMemoryRepository) AddUser(user models.User) error {
	obj.mu.Lock()
	user.ID = obj.count
	obj.user[obj.count] = user
	obj.count++
	obj.mu.Unlock()
	return nil
}

func (obj *InMemoryRepository) GetUser(id int) (*models.User, error) {
	user, ok := obj.user[id]
	if !ok {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &user, nil
}

func (obj *InMemoryRepository) UpdateUser(user models.User) error {
	_, ok := obj.user[user.ID]
	if !ok {
		return fmt.Errorf("user with id %d not found", user.ID)
	}
	obj.mu.Lock()
	obj.user[user.ID] = user
	obj.mu.Unlock()
	return nil
}

func (obj *InMemoryRepository) DeleteUser(id int) error {
	obj.mu.Lock()
	delete(obj.user, id)
	obj.mu.Unlock()
	return nil
}
