package services

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/alper.meric/messaging-system/clients"
	"github.com/alper.meric/messaging-system/config"
	"github.com/alper.meric/messaging-system/models"
	"github.com/alper.meric/messaging-system/repository"
)

// MessageServiceInterface defines the interface for the message service
type MessageServiceInterface interface {
	Start() error
	Stop() error
	Status() bool
	GetSentMessages(page, limit int) ([]models.Message, int, error)
}

// MessageService handles the message sending functionality
type MessageService struct {
	messageRepo   repository.MessageRepository
	cacheRepo     repository.CacheRepository
	messageClient *clients.MessageClient
	running       bool
	ticker        *time.Ticker
	stopChan      chan struct{}
	batchSize     int
	interval      time.Duration
	maxLength     int
	mutex         sync.Mutex
	isInitialized bool
}

// NewMessageService creates a new message service
func NewMessageService(
	cfg *config.Configuration,
	messageRepo repository.MessageRepository,
	cacheRepo repository.CacheRepository,
	messageClient *clients.MessageClient,
) MessageServiceInterface {
	return &MessageService{
		messageRepo:   messageRepo,
		cacheRepo:     cacheRepo,
		messageClient: messageClient,
		running:       false,
		batchSize:     cfg.App.MessageBatchSize,
		interval:      time.Duration(cfg.App.MessageSendInterval) * time.Minute,
		maxLength:     cfg.App.MaxContentLength,
		isInitialized: true,
	}
}

// Start begins the scheduled message sending
func (s *MessageService) Start() error {
	if !s.isInitialized {
		return errors.New("message service is not initialized")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.running {
		return errors.New("message service is already running")
	}

	log.Println("Starting message service...")
	s.ticker = time.NewTicker(s.interval)
	s.stopChan = make(chan struct{})
	s.running = true

	go s.run()
	return nil
}

// Stop stops the scheduled message sending
func (s *MessageService) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.running {
		return errors.New("message service is not running")
	}

	log.Println("Stopping message service...")
	s.ticker.Stop()
	s.stopChan <- struct{}{}
	s.running = false
	return nil
}

// Status returns whether the service is running
func (s *MessageService) Status() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.running
}

// GetSentMessages retrieves sent messages with pagination
func (s *MessageService) GetSentMessages(page, limit int) ([]models.Message, int, error) {
	return s.messageRepo.GetSentMessages(page, limit)
}

func (s *MessageService) run() {
	// Initial processing
	s.processMessages()

	for {
		select {
		case <-s.ticker.C:
			s.processMessages()
		case <-s.stopChan:
			log.Println("Message service stopped")
			return
		}
	}
}

func (s *MessageService) processMessages() {
	log.Println("Processing unsent messages...")

	// Get unsent messages from the repository
	messages, err := s.messageRepo.GetUnsentMessages(s.batchSize)
	if err != nil {
		log.Printf("Error getting unsent messages: %v", err)
		return
	}

	if len(messages) == 0 {
		log.Println("No unsent messages found")
		return
	}

	log.Printf("Found %d unsent messages to process", len(messages))

	// Process each message
	for _, msg := range messages {
		// Validate message content
		if len(msg.Content) > s.maxLength {
			log.Printf("Message %d content exceeds maximum length (%d > %d)", msg.ID, len(msg.Content), s.maxLength)
			continue
		}

		// Send the message using the HTTP client
		externalID, err := s.messageClient.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message %d: %v", msg.ID, err)
			continue
		}

		// Mark as sent in repository
		err = s.messageRepo.MarkMessageAsSent(msg.ID, externalID)
		if err != nil {
			log.Printf("Failed to mark message %d as sent: %v", msg.ID, err)
			continue
		}

		// Cache in Redis (bonus feature)
		sentAt := time.Now()
		err = s.cacheRepo.CacheMessageID(externalID, sentAt)
		if err != nil {
			log.Printf("Warning: Failed to cache message ID %s: %v", externalID, err)
			// Continue anyway, as this is a non-critical operation
		}

		log.Printf("Successfully sent message %d to %s (external ID: %s)", msg.ID, msg.PhoneNumber, externalID)
	}
}
