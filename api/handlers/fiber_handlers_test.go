package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockservices "github.com/alper.meric/messaging-system/mocks/services"
	"github.com/alper.meric/messaging-system/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MessageRequest, mesaj gönderme isteği modeli
type MessageRequest struct {
	PhoneNumber string `json:"phone_number"`
	Content     string `json:"content"`
}

// MessageControllerTestSuite, MessageController için test suite
type MessageControllerTestSuite struct {
	suite.Suite
	mockService  *mockservices.MessageServiceInterface
	app          *fiber.App
	controller   *MessageController
	testMessages []models.Message
}

// SetupTest, her test öncesi çalışacak kurulum fonksiyonu
func (suite *MessageControllerTestSuite) SetupTest() {
	// Mock servis oluşturuluyor
	suite.mockService = mockservices.NewMessageServiceInterface(suite.T())

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
		{
			ID:            2,
			PhoneNumber:   "+90987654321",
			Content:       "Test message 2",
			IsSent:        true,
			SentAt:        time.Now().Add(-2 * time.Hour),
			ExternalMsgID: "ext-2",
			CreatedAt:     time.Now().Add(-3 * time.Hour),
		},
	}

	// Fiber app ve controller oluşturuluyor
	suite.app = fiber.New()
	suite.controller = NewMessageController(suite.mockService)

	// Controller rotalarını kaydet
	suite.app.Post("/api/service", suite.controller.ServiceControl)
	suite.app.Get("/api/service/status", suite.controller.ServiceStatus)
	suite.app.Get("/api/messages", suite.controller.GetSentMessages)
}

// TestServiceControl, servis kontrol endpointlerini test eder
func (suite *MessageControllerTestSuite) TestServiceControl() {
	// Start action testi
	suite.mockService.EXPECT().Status().Return(false)
	suite.mockService.EXPECT().Start().Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/service?action=start", nil)
	resp, err := suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.True(suite.T(), result["success"].(bool))
	assert.Contains(suite.T(), result["message"].(string), "started")

	// Stop action testi
	suite.mockService.EXPECT().Status().Return(true)
	suite.mockService.EXPECT().Stop().Return(nil)

	req = httptest.NewRequest(http.MethodPost, "/api/service?action=stop", nil)
	resp, err = suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.True(suite.T(), result["success"].(bool))
	assert.Contains(suite.T(), result["message"].(string), "stopped")
}

// TestServiceStatus, servis durum endpointini test eder
func (suite *MessageControllerTestSuite) TestServiceStatus() {
	// Çalışırken durumu
	suite.mockService.EXPECT().Status().Return(true).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/service/status", nil)
	resp, err := suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.True(suite.T(), result["success"].(bool))
	assert.True(suite.T(), result["running"].(bool))

	// Dururken durumu
	suite.mockService.EXPECT().Status().Return(false).Once()

	req = httptest.NewRequest(http.MethodGet, "/api/service/status", nil)
	resp, err = suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.True(suite.T(), result["success"].(bool))
	assert.False(suite.T(), result["running"].(bool))
}

// TestGetSentMessages, gönderilmiş mesajları getirme endpointini test eder
func (suite *MessageControllerTestSuite) TestGetSentMessages() {
	// Mesajları getirme başarılı senaryosu
	suite.mockService.EXPECT().GetSentMessages(1, 10).Return(suite.testMessages, len(suite.testMessages), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/messages?page=1&limit=10", nil)
	resp, err := suite.app.Test(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var result models.MessageListResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.True(suite.T(), result.Success)
	assert.Equal(suite.T(), len(suite.testMessages), result.Total)
	assert.Equal(suite.T(), len(suite.testMessages), len(result.Messages))
	assert.Equal(suite.T(), suite.testMessages[0].ID, result.Messages[0].ID)
}

// TestMessageControllerSuite çalıştırma fonksiyonu
func TestMessageControllerSuite(t *testing.T) {
	suite.Run(t, new(MessageControllerTestSuite))
}
