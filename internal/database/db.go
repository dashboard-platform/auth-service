// Package database provides functionality for managing database connections and operations.
// It includes methods for initializing the database, performing migrations, and executing
// CRUD operations on user data. This package is essential for interacting with the application's
// persistent storage layer.

package database

import (
	"os"
	"time"

	"github.com/dashboard-platform/auth-service/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

// Database represents the database connection and provides methods for interacting with it.
// It includes a GORM database instance and a logger for logging database operations.
type Database struct {
	db     *gorm.DB       // GORM database instance.
	logger zerolog.Logger // Logger for database operations.
}

// Init initializes a new database connection using the provided DSN (Data Source Name).
// It retries the connection multiple times if the database is not ready.
//
// Parameters:
//   - dsn: The Data Source Name for connecting to the database.
//   - logger: A logger instance for logging database operations.
//
// Returns:
//   - *Database: A pointer to the initialized Database instance.
//   - error: An error if the connection fails after retries.
func Init(dsn string, logger zerolog.Logger) (*Database, error) {
	var (
		db  *gorm.DB
		err error
	)

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLog.Default.LogMode(gormLog.Silent),
		})

		if err == nil {
			break
		}

		log.Warn().Err(err).Msgf("DB is not ready, retrying... (%d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database after retries")
		os.Exit(1)
	}

	return &Database{
		db:     db,
		logger: logger,
	}, nil
}

// AutoMigrate performs database migrations for the User model.
//
// Returns:
//   - error: An error if the migration fails.
func (d *Database) AutoMigrate() error {
	if err := d.db.AutoMigrate(&models.User{}); err != nil {
		d.logger.Error().Err(err).Msgf("failed to execute automigration")
		return err
	}
	return nil
}

// Create inserts a new user record into the database.
//
// Parameters:
//   - user: The user data to insert.
//
// Returns:
//   - error: An error if the insertion fails.
func (d *Database) Create(user models.User) error {
	if err := d.db.Create(&user).Error; err != nil {
		d.logger.Error().Err(err).Msgf("failed to create user")
		return err
	}
	return nil
}

// Fetch retrieves a user record from the database by email.
//
// Parameters:
//   - email: The email of the user to retrieve.
//
// Returns:
//   - models.User: The retrieved user data.
//   - error: An error if the retrieval fails.
func (d *Database) Fetch(email string) (models.User, error) {
	var user models.User
	if err := d.db.Where("email = ?", email).First(&user).Error; err != nil {
		d.logger.Error().Err(err).Msgf("failed to fetch user with email %s", email)
		return models.User{}, err
	}
	return user, nil
}

// FetchByID retrieves a user record from the database by ID.
//
// Parameters:
//   - id: The ID of the user to retrieve.
//
// Returns:
//   - models.User: The retrieved user data.
//   - error: An error if the retrieval fails.
func (d *Database) FetchByID(id string) (models.User, error) {
	var user models.User
	if err := d.db.Where("id = ?", id).First(&user).Error; err != nil {
		d.logger.Error().Err(err).Msgf("failed to fetch user with id %s", id)
		return models.User{}, err
	}
	return user, nil
}

// Close closes the database connection.
//
// Logs an error if the connection cannot be closed.
func (d *Database) Close() {
	sqlDB, err := d.db.DB()
	if err != nil {
		d.logger.Error().Err(err).Msgf("failed to extract sql.DB")
		return
	}
	if err = sqlDB.Close(); err != nil {
		d.logger.Error().Err(err).Msgf("failed to close database")
	}
}
