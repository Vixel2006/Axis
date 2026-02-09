package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils" // Import the utils package
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService services.MessageService
}

func NewMessageHandler(ms services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: ms,
	}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message.SenderID = int(userID)

	createdMessage, err := h.messageService.CreateMessage(c.Request.Context(), &message)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "meeting not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, createdMessage)
}

func (h *MessageHandler) GetMessageByID(c *gin.Context) {
	idStr := c.Param("messageID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := h.messageService.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve message"})
		return
	}

	if message == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) GetMessagesInMeeting(c *gin.Context) {
	meetingIDStr := c.Param("meetingID")
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100") // Default limit 100
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offsetStr := c.DefaultQuery("offset", "0") // Default offset 0
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	messages, err := h.messageService.GetMessagesInMeeting(c.Request.Context(), meetingID, limit, offset)
	if err != nil {
		if err.Error() == "meeting not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages for meeting"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	idStr := c.Param("messageID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message.ID = id // Ensure the ID from the URL is used

	// Pass userID to the service for authorization
	updatedMessage, err := h.messageService.UpdateMessage(c.Request.Context(), int(userID), &message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	if updatedMessage == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, updatedMessage)
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	idStr := c.Param("messageID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	// Pass userID to the service for authorization
	err = h.messageService.DeleteMessage(c.Request.Context(), int(userID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
