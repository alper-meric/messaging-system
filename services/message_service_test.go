package services

import (
	"testing"
	"time"

	"github.com/alper.meric/messaging-system/clients"
	"github.com/alper.meric/messaging-system/config"
	mocks "github.com/alper.meric/messaging-system/mocks/repository"
	"github.com/alper.meric/messaging-system/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MessageServiceTestSuite, MessageService için test suite
type MessageServiceTestSuite struct {
	suite.Suite
	mockMsgRepo    *mocks.MessageRepository
	mockCacheRepo  *mocks.CacheRepository
	messageClient  *clients.MessageClient
	config         *config.Configuration
	messageService MessageServiceInterface
	testMessages   []models.Message
	unsentMessages []models.Message
}

// SetupTest, her test öncesi çalışacak kurulum fonksiyonu
func (suite *MessageServiceTestSuite) SetupTest() {
	// Mock nesneler oluşturuluyor
	suite.mockMsgRepo = mocks.NewMessageRepository(suite.T())
	suite.mockCacheRepo = mocks.NewCacheRepository(suite.T())

	// Test verileri hazırlanıyor
	suite.testMessages = []models.Message{
		{
			ID:            1,
			PhoneNumber:   "+90123456789",
			Content:       "Test message 1",
			IsSent:        true,
			SentAt:        time.Now().Add(-1 * time.Hour),
			ExternalMsgID: "ext-1",
			CreatedAt:     time.Now().Add(-2 * time.Hour),
		},
	}

	suite.unsentMessages = []models.Message{
		{
			ID:          2,
			PhoneNumber: "+90123456789",
			Content:     "Test message 2",
			IsSent:      false,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		},
	}

	// Test konfigürasyonu
	suite.config = &config.Configuration{
		App: config.AppConfig{
			MessageBatchSize:    2,
			WebhookURL:          "https://test.example.com",
			MaxContentLength:    1000,
			MessageSendDryRun:   true,
			MessageSendInterval: 2,
		},
	}

	// HTTP client oluştur
	suite.messageClient = clients.NewMessageClient(
		suite.config.App.WebhookURL,
		suite.config.App.MessageSendDryRun,
	)

	// MessageService oluşturuluyor
	suite.messageService = NewMessageService(suite.config, suite.mockMsgRepo, suite.mockCacheRepo, suite.messageClient)
}

// TestNewMessageService, servis başlatma testleri
func (suite *MessageServiceTestSuite) TestNewMessageService() {
	assert.NotNil(suite.T(), suite.messageService, "NewMessageService fonksiyonu nil döndürmemeli")
	assert.False(suite.T(), suite.messageService.Status(), "Yeni servis varsayılan olarak çalışır durumda olmamalı")
}

// TestStartStop, servis başlatma ve durdurma testleri
func (suite *MessageServiceTestSuite) TestStartStop() {
	// GetUnsentMessages mock ayarı (processMessages için)
	suite.mockMsgRepo.EXPECT().GetUnsentMessages(suite.config.App.MessageBatchSize).Return([]models.Message{}, nil).Maybe()

	// Başlat
	err := suite.messageService.Start()
	assert.NoError(suite.T(), err, "Start fonksiyonu hata döndürmemeli")
	assert.True(suite.T(), suite.messageService.Status(), "Start fonksiyonu çağrıldıktan sonra Status() true olmalı")

	// Zaten çalışırken tekrar başlatma
	err = suite.messageService.Start()
	assert.Error(suite.T(), err, "Zaten çalışan servisi başlatmaya çalışırken hata olmalı")

	// Durdur
	err = suite.messageService.Stop()
	assert.NoError(suite.T(), err, "Stop fonksiyonu hata döndürmemeli")
	assert.False(suite.T(), suite.messageService.Status(), "Stop fonksiyonu çağrıldıktan sonra Status() false olmalı")

	// Zaten durdurulmuşken tekrar durdurma
	err = suite.messageService.Stop()
	assert.Error(suite.T(), err, "Zaten durmuş servisi durdurmaya çalışırken hata olmalı")
}

// TestStatus, servis durumu testleri
func (suite *MessageServiceTestSuite) TestStatus() {
	// GetUnsentMessages mock ayarı (processMessages için)
	suite.mockMsgRepo.EXPECT().GetUnsentMessages(suite.config.App.MessageBatchSize).Return([]models.Message{}, nil).Maybe()

	// Başlangıç durumu
	assert.False(suite.T(), suite.messageService.Status(), "Başlangıçta servis durumu false olmalı")

	// Çalışırken durumu
	_ = suite.messageService.Start()
	assert.True(suite.T(), suite.messageService.Status(), "Başlatıldıktan sonra servis durumu true olmalı")
}

// TestGetSentMessages, gönderilmiş mesajları getirme testi
func (suite *MessageServiceTestSuite) TestGetSentMessages() {
	// Mock davranışını ayarla
	suite.mockMsgRepo.EXPECT().GetSentMessages(1, 10).Return(suite.testMessages, 1, nil)

	// GetSentMessages fonksiyonunu test et
	messages, total, err := suite.messageService.GetSentMessages(1, 10)

	assert.NoError(suite.T(), err, "GetSentMessages fonksiyonu hata döndürmemeli")
	assert.Equal(suite.T(), 1, total, "Toplam mesaj sayısı 1 olmalı")
	assert.Len(suite.T(), messages, 1, "Dönen mesaj sayısı 1 olmalı")
	assert.Equal(suite.T(), suite.testMessages[0].ID, messages[0].ID, "Mesaj ID'leri eşleşmeli")
}

// TestProcessMessages, mesaj işleme testi
func (suite *MessageServiceTestSuite) TestProcessMessages() {
	// Mock davranışlarını ayarla
	suite.mockMsgRepo.EXPECT().GetUnsentMessages(suite.config.App.MessageBatchSize).Return(suite.unsentMessages, nil)

	// Beklenen external ID (dry run modunda)
	expectedMsgID := "dry-run-id-2"
	suite.mockMsgRepo.EXPECT().MarkMessageAsSent(2, expectedMsgID).Return(nil)
	suite.mockCacheRepo.EXPECT().CacheMessageID(expectedMsgID, mock.AnythingOfType("time.Time")).Return(nil)

	// Servis tipine dönüştür
	concreteService, ok := suite.messageService.(*MessageService)
	assert.True(suite.T(), ok, "Servis, MessageService tipine dönüştürülebilmeli")

	// MessageClient'ın dry run modunda olduğunu kontrol et
	// Not: Gerçek HTTP isteği yapılmayacak, messageClient dry run modunda

	// ProcessMessages metodunu doğrudan çağır
	concreteService.processMessages()

	// Beklenen mock çağrılarının gerçekleştiğini kontrol et (mock kütüphanesi tarafından otomatik olarak yapılır)
}

// TestMessageServiceSuite çalıştırma fonksiyonu
func TestMessageServiceSuite(t *testing.T) {
	suite.Run(t, new(MessageServiceTestSuite))
}
