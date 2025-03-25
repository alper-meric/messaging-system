package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alper.meric/messaging-system/models"
)

// MessageClient, mesaj gönderimi için HTTP istemcisini temsil eder
type MessageClient struct {
	webhookURL string
	client     *http.Client
	dryRun     bool
}

// NewMessageClient, yeni bir MessageClient oluşturur
func NewMessageClient(webhookURL string, dryRun bool) *MessageClient {
	return &MessageClient{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second, // 10 saniyelik timeout
		},
		dryRun: dryRun,
	}
}

// MessageResponse, mesaj gönderimi yanıtını temsil eder
type MessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

// SendMessage, belirtilen mesajı dış servise gönderir ve mesaj ID'sini döndürür
func (c *MessageClient) SendMessage(msg models.Message) (string, error) {
	// Eğer dry run modunda ise, mesajları gerçekten göndermez
	if c.dryRun {
		log.Printf("DRY RUN: Would send message to %s: %s", msg.PhoneNumber, msg.Content)
		return fmt.Sprintf("dry-run-id-%d", msg.ID), nil
	}

	// İstek verilerini hazırla
	payload := map[string]string{
		"to":      msg.PhoneNumber,
		"content": msg.Content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// HTTP POST isteği gönder
	req, err := http.NewRequest(http.MethodPost, c.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Yanıt durumunu kontrol et
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("external service returned error status: %d", resp.StatusCode)
	}

	// Dış mesaj ID'sini almak için yanıtı ayrıştır
	var responseData MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Dış mesaj ID'sini çıkar
	if responseData.MessageID == "" {
		return "", errors.New("external service did not return a valid message ID")
	}

	return responseData.MessageID, nil
}
