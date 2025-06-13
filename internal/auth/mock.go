// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import (
	"errors"
	"net/mail"
	"sync"
	"time"

	"github.com/dashboard-platform/auth-service/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// MockService represents a mock implementation of the authentication service.
// It is used for testing purposes and simulates user registration, login, and retrieval.
type MockService struct {
	mu     sync.Mutex             // Mutex to ensure thread-safe operations.
	users  map[string]models.User // Map to store user data by user ID.
	logins map[string]string      // Map to store login credentials by email.
}

// NewMockService creates a new instance of the mock authentication service.
//
// Returns:
//   - *MockService: A pointer to the newly created mock service.
func NewMockService() *MockService {
	return &MockService{
		users:  make(map[string]models.User),
		logins: make(map[string]string),
	}
}

// Register registers a new user in the mock service.
//
// Parameters:
//   - input: The registration data containing the user's email and password.
//
// Returns:
//   - string: The ID of the newly registered user.
//   - error: An error if the registration fails.
func (m *MockService) Register(input models.RegisterAPI) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if input.Name == "" {
		log.Error().Msg("name is empty for mock registration")
		return "", errors.New("name is empty")
	}

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
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: input.Password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	m.users[id] = user
	m.logins[input.Email] = input.Password
	return id, nil
}

// Login authenticates a user in the mock service.
//
// Parameters:
//   - input: The login data containing the user's email and password.
//
// Returns:
//   - string: The ID of the authenticated user.
//   - error: An error if the authentication fails.
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

// LoginOrRegisterOAuth handles user login or registration via OAuth in the mock service.
//
// This method simulates the behavior of logging in or registering a user using an OAuth provider.
// If the user already exists in the mock database, their ID is returned. If the user does not exist,
// a new user is created with the provided email, and their ID is returned.
//
// Parameters:
//   - email: The email address of the user obtained from the OAuth provider.
//
// Returns:
//   - string: The ID of the user (existing or newly created).
//   - error: An error if the operation fails (e.g., database issues).
func (m *MockService) LoginOrRegisterOAuth(data models.RegisterAPI) (string, error) {
	defer m.mu.Unlock()
	m.mu.Lock()

	for id, u := range m.users {
		if u.Email == data.Email {
			return id, nil
		}
	}

	t := time.Now()
	user := models.User{
		ID:           uuid.New().String(),
		Name:         data.Name, // Adaugă numele din data
		Email:        data.Email,
		PasswordHash: "",
		AuthProvider: "google", // Asigură-te că setezi provider-ul corect
		CreatedAt:    t,
		UpdatedAt:    t,
	}

	m.users[user.ID] = user
	return user.ID, nil
}

// GetUserByID retrieves a user by their ID from the mock service.
//
// Parameters:
//   - id: The ID of the user to retrieve.
//
// Returns:
//   - models.User: The user data.
//   - error: An error if the user is not found.
func (m *MockService) GetUserByID(id string) (models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, ok := m.users[id]
	if !ok {
		return models.User{}, errors.New("not found")
	}
	return u, nil
}

// AddUser adds a user to the mock service.
//
// Parameters:
//   - input: The user data to add.
func (m *MockService) AddUser(input models.User) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.users[input.ID] = input
	m.logins[input.Email] = input.PasswordHash
}

// DeleteUser deletes a user from the mock service.
//
// Parameters:
//   - id: The ID of the user to delete.
func (m *MockService) DeleteUser(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.users, id)
}
