package handlers

import (
	"log"
	"strconv"

	"github.com/alper.meric/messaging-system/models"
	"github.com/alper.meric/messaging-system/services"
	"github.com/gofiber/fiber/v2"
)

// MessageController handles all message related HTTP operations
type MessageController struct {
	messageService services.MessageServiceInterface
}

// NewMessageController creates a new message controller
func NewMessageController(
	messageService services.MessageServiceInterface,
) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// ServiceControl handles service control actions for Fiber
// @Summary Controls the message service
// @Description Starts or stops the message sending service
// @Tags service
// @Accept json
// @Produce json
// @Param action query string true "Action to perform (start/stop)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /service [post]
func (mc *MessageController) ServiceControl(c *fiber.Ctx) error {
	// Parse action parameter
	action := c.Query("action")

	if action == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "action parameter is required (start/stop)",
		})
	}

	var err error
	var message string

	switch action {
	case "start":
		err = mc.messageService.Start()
		message = "Message service started successfully"
	case "stop":
		err = mc.messageService.Stop()
		message = "Message service stopped successfully"
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid action. Use 'start' or 'stop'",
		})
	}

	if err != nil {
		log.Printf("Service control error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"running": mc.messageService.Status(),
	})
}

// ServiceStatus retrieves the current status of the message service
// @Summary Gets service status
// @Description Retrieves the current running status of the message service
// @Tags service
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /service/status [get]
func (mc *MessageController) ServiceStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"running": mc.messageService.Status(),
	})
}

// GetSentMessages lists sent messages using Fiber
// @Summary Retrieves sent messages
// @Description Gets a list of sent messages with pagination
// @Tags messages
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} models.MessageListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /messages [get]
func (mc *MessageController) GetSentMessages(c *fiber.Ctx) error {
	// Parse pagination parameters
	page := 1
	limit := 10
	var err error

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
	}

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil || limit < 1 || limit > 100 {
			limit = 10
		}
	}

	// Get sent messages from service instead of repository
	messages, total, err := mc.messageService.GetSentMessages(page, limit)
	if err != nil {
		log.Printf("Error retrieving sent messages: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve messages",
		})
	}

	// Calculate total pages
	pages := total / limit
	if total%limit > 0 {
		pages++
	}

	// Create response
	response := models.MessageListResponse{
		Success:  true,
		Messages: messages,
		Page:     page,
		Limit:    limit,
		Total:    total,
		Pages:    pages,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Artık doğrudan controller kullanıldığı için uyumluluk fonksiyonlarına ihtiyaç kalmadı.
// Compatibility functions are removed as we now use controller directly.
