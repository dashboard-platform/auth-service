package database

import (
	"errors"
	"sync"

	"auth-service/models"
)

type MockupDatabase struct {
	mu sync.Mutex
	db map[string]models.User
}

func NewMockupDatabase() *MockupDatabase {
	return &MockupDatabase{
		mu: sync.Mutex{},
		db: make(map[string]models.User),
	}
}

func (m *MockupDatabase) AutoMigrate() error {
	return nil
}

func (m *MockupDatabase) Create(user models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.db[user.ID] = user
	return nil
}

func (m *MockupDatabase) Fetch(email string) (models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.db {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (m *MockupDatabase) FetchByID(id string) (models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user, ok := m.db[id]; ok {
		return user, nil
	}
	return models.User{}, errors.New("user not found")
}

func (m *MockupDatabase) Close() {
	return
}
