package repository

import (
	"context"
	"fmt"
	"pasha/models"
	"sync"
)

type UserInMemoryRepository struct {
	user  map[int]models.User
	count int
	mu    sync.Mutex
}

func NewUserInMemoryRepository(ctx context.Context) *UserInMemoryRepository {
	return &UserInMemoryRepository{
		user:  make(map[int]models.User),
		count: 1,
	}
}

func (obj *UserInMemoryRepository) AddUser(ctx context.Context, user models.User) (int, error) {
	obj.mu.Lock()
	user.ID = obj.count
	obj.user[obj.count] = user
	obj.count++
	obj.mu.Unlock()
	return user.ID, nil
}

func (obj *UserInMemoryRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	user, ok := obj.user[id]
	if !ok {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &user, nil
}

func (obj *UserInMemoryRepository) UpdateUser(ctx context.Context, user models.User) error {
	_, ok := obj.user[user.ID]
	if !ok {
		return fmt.Errorf("user with id %d not found", user.ID)
	}
	obj.mu.Lock()
	obj.user[user.ID] = user
	obj.mu.Unlock()
	return nil
}

func (obj *UserInMemoryRepository) DeleteUser(ctx context.Context, id int) error {
	obj.mu.Lock()
	delete(obj.user, id)
	obj.mu.Unlock()
	return nil
}
