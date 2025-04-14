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

type Database struct {
	db     *gorm.DB
	logger zerolog.Logger
}

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

func (d *Database) AutoMigrate() error {
	if err := d.db.AutoMigrate(&models.User{}); err != nil {
		d.logger.Error().Err(err).Msgf("failed to execute automigration")
		return err
	}
	return nil
}

func (d *Database) Create(user models.User) error {
	if err := d.db.Create(&user).Error; err != nil {
		d.logger.Error().Err(err).Msgf("failed to create user")
		return err
	}
	return nil
}

func (d *Database) Fetch(email string) (models.User, error) {
	var user models.User
	if err := d.db.Where("email = ?", email).First(&user).Error; err != nil {
		d.logger.Error().Err(err).Msgf("failed to fetch user with email %s", email)
		return models.User{}, err
	}
	return user, nil
}

func (d *Database) FetchByID(id string) (models.User, error) {
	var user models.User
	if err := d.db.Where("id = ?", id).First(&user).Error; err != nil {
		d.logger.Error().Err(err).Msgf("failed to fetch user with id %s", id)
		return models.User{}, err
	}
	return user, nil
}

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
