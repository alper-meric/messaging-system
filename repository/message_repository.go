package repository

import (
	"time"

	"github.com/alper.meric/messaging-system/models"
)

// MessageRepository provides abstraction for message database operations
type MessageRepository interface {
	// Retrieves unsent messages
	GetUnsentMessages(limit int) ([]models.Message, error)

	// Marks a message as sent
	MarkMessageAsSent(id int, externalMsgID string) error

	// Retrieves sent messages with pagination
	GetSentMessages(page, limit int) ([]models.Message, int, error)

	// Adds a new message
	AddMessage(message models.Message) (int, error)
}

// CacheRepository provides abstraction for message caching operations
type CacheRepository interface {
	// Caches message ID and send time
	CacheMessageID(messageID string, sentAt time.Time) error

	// Retrieves cached message information
	GetCachedMessage(messageID string) (time.Time, error)
}
