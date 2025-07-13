package store

import (
	"errors"
	"sync"

	"bookapp/internal/models"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	mu              sync.RWMutex
	users           = make(map[string]*models.User)
	sessions        = make(map[string]string)
)

type MemoryStore struct{}

func NewInMemoryUserStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Create(u *models.User) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := users[u.Email]; ok {
		return ErrUserExists
	}
	users[u.Email] = u
	return nil
}

func (m *MemoryStore) FindByEmail(email string) (*models.User, error) {
	mu.RLock()
	defer mu.RUnlock()
	u, ok := users[email]
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (m *MemoryStore) Update(u *models.User) error {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := users[u.Email]; !ok {
		return ErrUserNotFound
	}
	users[u.Email] = u
	return nil
}
