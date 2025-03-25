package models

import (
	"time"

	"gorm.io/gorm"
)

// Message represents a message in the system
type Message struct {
	ID            int            `json:"id" gorm:"primaryKey"`
	Content       string         `json:"content" gorm:"type:text;not null"`
	PhoneNumber   string         `json:"phoneNumber" gorm:"type:varchar(20);not null"`
	IsSent        bool           `json:"isSent" gorm:"default:false"`
	SentAt        time.Time      `json:"sentAt,omitempty" gorm:"default:null"`
	ExternalMsgID string         `json:"externalMsgId,omitempty" gorm:"default:null"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName sets the table name for the Message model
func (Message) TableName() string {
	return "messages"
}

// MessageResponse represents a response from the external message service
type MessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

// MessageRequest represents a request to the external message service
type MessageRequest struct {
	To      string `json:"to" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// MessageListResponse represents a paginated list of messages
type MessageListResponse struct {
	Success  bool      `json:"success"`
	Messages []Message `json:"messages"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
	Pages    int       `json:"pages"`
}

// ServiceStatus represents the status of the message service
type ServiceStatus struct {
	IsRunning   bool      `json:"isRunning"`
	LastRunTime time.Time `json:"lastRunTime,omitempty"`
}
