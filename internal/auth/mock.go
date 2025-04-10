package auth

import (
	"errors"
	"net/mail"
	"sync"
	"time"

	"auth-service/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type MockService struct {
	mu     sync.Mutex
	users  map[string]models.User
	logins map[string]string
}

func NewMockService() *MockService {
	return &MockService{
		users:  make(map[string]models.User),
		logins: make(map[string]string),
	}
}

func (m *MockService) Register(input models.RegisterAPI) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := mail.ParseAddress(input.Email); err != nil {
		log.Error().Err(err).Msg("email address is not valid")
		return "", err
	}

	if input.Password == "" {
		log.Error().Msg("password is empty")
		return "", errors.New("password is empty")
	}

	id := uuid.New().String()
	user := models.User{
		ID:           id,
		Email:        input.Email,
		PasswordHash: input.Password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	m.users[id] = user
	m.logins[input.Email] = input.Password
	return id, nil
}

func (m *MockService) Login(input models.LoginAPI) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if input.Email == "" || input.Password == "" {
		log.Error().Msg("email or password is empty")
		return "", errors.New("email or password is empty")
	}

	for id, u := range m.users {
		if u.Email == input.Email && m.logins[u.Email] == input.Password {
			return id, nil
		}
	}
	return "", errors.New("invalid credentials")
}

func (m *MockService) GetUserByID(id string) (models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, ok := m.users[id]
	if !ok {
		return models.User{}, errors.New("not found")
	}
	return u, nil
}

func (m *MockService) AddUser(input models.User) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.users[input.ID] = input
	m.logins[input.Email] = input.PasswordHash
}

func (m *MockService) DeleteUser(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.users, id)
}
