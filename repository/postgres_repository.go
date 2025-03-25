package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/alper.meric/messaging-system/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresRepository implements the MessageRepository interface using PostgreSQL database
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new PostgresRepository instance
func NewPostgresRepository(host string, port int, user, password, dbname string) (*PostgresRepository, error) {
	// Create connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(log.Writer(), "GORM: ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Automatic schema migration
	err = db.AutoMigrate(&models.Message{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	log.Println("Database connection established and migration completed")

	return &PostgresRepository{
		db: db,
	}, nil
}

// GetDB provides access to the database object (for testing if needed)
func (r *PostgresRepository) GetDB() *gorm.DB {
	return r.db
}

// GetUnsentMessages retrieves unsent messages
func (r *PostgresRepository) GetUnsentMessages(limit int) ([]models.Message, error) {
	var messages []models.Message
	result := r.db.Where("is_sent = ?", false).
		Order("created_at asc").
		Limit(limit).
		Find(&messages)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve unsent messages: %w", result.Error)
	}

	return messages, nil
}

// MarkMessageAsSent marks a message as sent
func (r *PostgresRepository) MarkMessageAsSent(id int, externalMsgID string) error {
	result := r.db.Model(&models.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_sent":         true,
			"sent_at":         time.Now(),
			"external_msg_id": externalMsgID,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to mark message as sent: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("message with ID %d not found", id)
	}

	return nil
}

// GetSentMessages retrieves sent messages with pagination
func (r *PostgresRepository) GetSentMessages(page, limit int) ([]models.Message, int, error) {
	var messages []models.Message
	var total int64

	// Calculate offset
	offset := (page - 1) * limit

	// Count total sent messages
	r.db.Model(&models.Message{}).Where("is_sent = ?", true).Count(&total)

	// Get sent messages with pagination
	result := r.db.Where("is_sent = ?", true).
		Order("sent_at desc").
		Offset(offset).
		Limit(limit).
		Find(&messages)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to retrieve sent messages: %w", result.Error)
	}

	return messages, int(total), nil
}

// AddMessage adds a new message
func (r *PostgresRepository) AddMessage(message models.Message) (int, error) {
	// Set defaults
	message.IsSent = false
	message.CreatedAt = time.Now()

	// Add message to database
	result := r.db.Create(&message)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to add message: %w", result.Error)
	}

	return message.ID, nil
}
