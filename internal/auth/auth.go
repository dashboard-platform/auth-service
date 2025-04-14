package auth

import (
	"errors"
	"net/mail"
	"time"

	"github.com/dashboard-platform/auth-service/internal/database"
	"github.com/dashboard-platform/auth-service/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service struct {
	db database.Repository
}

func NewService(db database.Repository) *Service {
	return &Service{
		db: db,
	}
}

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
		CreatedAt:    t,
		UpdatedAt:    t,
	}
	return user.ID, s.db.Create(user)
}

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

func (s *Service) GetUserByID(id string) (models.User, error) {
	if id == "" {
		return models.User{}, errors.New("id cannot be empty")
	}

	return s.db.Fetch(id)
}
