// Package auth provides core authentication functionalities, including user registration,
// login, password hashing and verification, and JWT token management. It is designed
// to handle secure authentication workflows and integrate with the application's database
// and middleware layers.
package auth

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/dashboard-platform/auth-service/internal/database"
	"github.com/dashboard-platform/auth-service/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Service represents the authentication service.
// It provides methods for user registration, login, and user retrieval.
type Service struct {
	db database.Repository // Database repository for user data.
}

// NewService creates a new instance of the authentication service.
//
// Parameters:
//   - db: The database repository to be used by the service.
//
// Returns:
//   - *Service: A pointer to the newly created authentication service.
func NewService(db database.Repository) *Service {
	return &Service{
		db: db,
	}
}

// Register registers a new user in the system.
//
// Parameters:
//   - data: The registration data containing the user's email and password.
//
// Returns:
//   - string: The ID of the newly registered user.
//   - error: An error if the registration fails.
func (s *Service) Register(data models.RegisterAPI) (string, error) {
	if _, err := mail.ParseAddress(data.Email); err != nil {
		log.Error().Err(err).Msg("email address is not valid")
		return "", err
	}

	if data.Password == "" {
		log.Error().Msg("password is empty")
		return "", errors.New("password is empty")
	}

	hashed, err := hashPassword(data.Password)
	if err != nil {
		log.Error().Err(err).Msg("hashing password failed")
		return "", err
	}

	t := time.Now()
	user := models.User{
		ID:           uuid.New().String(),
		Email:        data.Email,
		PasswordHash: hashed,
		AuthProvider: "password",
		CreatedAt:    t,
		UpdatedAt:    t,
	}
	return user.ID, s.db.Create(user)
}

// Login authenticates a user by verifying their email and password.
//
// Parameters:
//   - data: The login data containing the user's email and password.
//
// Returns:
//   - string: The ID of the authenticated user.
//   - error: An error if the authentication fails.
func (s *Service) Login(data models.LoginAPI) (string, error) {
	if data.Email == "" || data.Password == "" {
		log.Error().Msg("email or password is empty")
		return "", errors.New("email or password is empty")
	}

	fetched, err := s.db.Fetch(data.Email)
	if err != nil {
		log.Error().Err(err).Msg("email not found")
		return "", err
	}

	if fetched.AuthProvider != "password" {
		return "", fmt.Errorf("this email is registered via %s login", fetched.AuthProvider)
	}

	ok, err := verifyPassword(fetched.PasswordHash, data.Password)
	if err != nil {
		log.Error().Err(err).Msgf("password validation failed for %s", data.Email)
		return "", err
	}

	if ok {
		return fetched.ID, nil
	}

	log.Warn().Msgf("credentials validation failed for %s", data.Email)
	return "", errors.New("email or password is wrong")
}

// LoginOrRegisterOAuth handles user login or registration via OAuth.
//
// This method is used when a user logs in using an OAuth provider (for now, Google).
// If the user already exists in the database, their ID is returned. If the user
// does not exist, a new user is created with the provided email, and their ID
// is returned.
//
// Parameters:
//   - email: The email address of the user obtained from the OAuth provider.
//
// Returns:
//   - string: The ID of the user (existing or newly created).
//   - error: An error if the operation fails (e.g., database issues).
func (s *Service) LoginOrRegisterOAuth(email string) (string, error) {
	user, err := s.db.Fetch(email)
	if err == nil {
		if user.AuthProvider != "google" {
			return "", fmt.Errorf("this email is registered with %s login", user.AuthProvider)
		}
		return user.ID, nil
	}

	t := time.Now()
	user = models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: "",
		AuthProvider: "google",
		CreatedAt:    t,
		UpdatedAt:    t,
	}
	return user.ID, s.db.Create(user)
}

// GetUserByID retrieves a user by their ID.
//
// Parameters:
//   - id: The ID of the user to retrieve.
//
// Returns:
//   - models.User: The user data.
//   - error: An error if the retrieval fails or the ID is empty.
func (s *Service) GetUserByID(id string) (models.User, error) {
	if id == "" {
		return models.User{}, errors.New("id cannot be empty")
	}

	return s.db.FetchByID(id)
}
