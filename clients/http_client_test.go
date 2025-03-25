package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alper.meric/messaging-system/models"
	"github.com/stretchr/testify/assert"
)

func TestMessageClientSendMessage(t *testing.T) {
	// Test setup
	msg := models.Message{
		ID:          1,
		PhoneNumber: "+90123456789",
		Content:     "Test message",
	}

	t.Run("successful send", func(t *testing.T) {
		// Mock server setup
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request method
			assert.Equal(t, http.MethodPost, r.Method)

			// Check content type
			contentType := r.Header.Get("Content-Type")
			assert.Equal(t, "application/json", contentType)

			// Check request body
			var requestBody map[string]string
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			assert.NoError(t, err)
			assert.Equal(t, msg.PhoneNumber, requestBody["to"])
			assert.Equal(t, msg.Content, requestBody["content"])

			// Send response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]string{
				"message":   "Success",
				"messageId": "test-id-123",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create client with mock server URL
		client := NewMessageClient(server.URL, false)

		// Test send message
		externalID, err := client.SendMessage(msg)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, "test-id-123", externalID)
	})

	t.Run("server error", func(t *testing.T) {
		// Mock server setup that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		// Create client with mock server URL
		client := NewMessageClient(server.URL, false)

		// Test send message
		_, err := client.SendMessage(msg)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "returned error status: 500")
	})

	t.Run("invalid response", func(t *testing.T) {
		// Mock server setup that returns invalid JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{invalid json"))
		}))
		defer server.Close()

		// Create client with mock server URL
		client := NewMessageClient(server.URL, false)

		// Test send message
		_, err := client.SendMessage(msg)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode response")
	})

	t.Run("missing message ID", func(t *testing.T) {
		// Mock server setup that returns a response without message ID
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]string{
				"message": "Success",
				// messageId is missing
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create client with mock server URL
		client := NewMessageClient(server.URL, false)

		// Test send message
		_, err := client.SendMessage(msg)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "did not return a valid message ID")
	})

	t.Run("dry run mode", func(t *testing.T) {
		// Create client in dry run mode
		client := NewMessageClient("http://example.com", true)

		// Test send message in dry run mode
		externalID, err := client.SendMessage(msg)

		// Assertions
		assert.NoError(t, err)
		assert.Contains(t, externalID, "dry-run-id-1")
	})
}
