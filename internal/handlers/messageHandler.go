package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils" // Import the utils package
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type MessageHandler struct {
	messageService services.MessageService
	log            zerolog.Logger
}

func NewMessageHandler(ms services.MessageService, logger zerolog.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: ms,
		log:            logger,
	}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	h.log.Info().Msg("Handling CreateMessage request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in CreateMessage")
		return
	}

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for CreateMessage")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message.SenderID = int(userID)
	h.log.Debug().Int("sender_id", message.SenderID).Int("meeting_id", message.MeetingID).Msg("CreateMessage request body")

	createdMessage, err := h.messageService.CreateMessage(c.Request.Context(), &message)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("sender_id", message.SenderID).Int("meeting_id", message.MeetingID).Msg("User forbidden from creating message in meeting")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("meeting_id", message.MeetingID).Msg("Meeting not found for message creation")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("sender_id", message.SenderID).Int("meeting_id", message.MeetingID).Msg("Failed to create message via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	h.log.Info().Int("message_id", createdMessage.ID).Int("sender_id", createdMessage.SenderID).Msg("Message created successfully")
	c.JSON(http.StatusCreated, createdMessage)
}

func (h *MessageHandler) GetMessageByID(c *gin.Context) {
	h.log.Info().Msg("Handling GetMessageByID request")
	idStr := c.Param("messageID")
	h.log.Debug().Str("messageID_param", idStr).Msg("Parsing message ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("messageID_param", idStr).Msg("Invalid message ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := h.messageService.GetMessageByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error().Err(err).Int("message_id", id).Msg("Failed to retrieve message via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve message"})
		return
	}

	if message == nil {
		h.log.Info().Int("message_id", id).Msg("Message not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	h.log.Info().Int("message_id", id).Msg("Message retrieved successfully")
	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) GetMessagesInMeeting(c *gin.Context) {
	h.log.Info().Msg("Handling GetMessagesInMeeting request")
	meetingIDStr := c.Param("meetingID")
	h.log.Debug().Str("meetingID_param", meetingIDStr).Msg("Parsing meeting ID")
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("meetingID_param", meetingIDStr).Msg("Invalid meeting ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "100") // Default limit 100
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.log.Error().Err(err).Str("limit_param", limitStr).Msg("Invalid limit parameter format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offsetStr := c.DefaultQuery("offset", "0") // Default offset 0
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.log.Error().Err(err).Str("offset_param", offsetStr).Msg("Invalid offset parameter format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}
	h.log.Debug().Int("meeting_id", meetingID).Int("limit", limit).Int("offset", offset).Msg("Retrieving messages in meeting")

	messages, err := h.messageService.GetMessagesInMeeting(c.Request.Context(), meetingID, limit, offset)
	if err != nil {
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("meeting_id", meetingID).Msg("Meeting not found for message retrieval")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to retrieve messages for meeting via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages for meeting"})
		return
	}

	h.log.Info().Int("meeting_id", meetingID).Int("messages_count", len(messages)).Msg("Messages in meeting retrieved successfully")
	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	h.log.Info().Msg("Handling UpdateMessage request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in UpdateMessage")
		return // GetUserIDFromContext already handles the error response
	}

	idStr := c.Param("messageID")
	h.log.Debug().Str("messageID_param", idStr).Msg("Parsing message ID for update")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("messageID_param", idStr).Msg("Invalid message ID format for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for UpdateMessage")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message.ID = id // Ensure the ID from the URL is used
	h.log.Debug().Int("message_id", id).Int("user_id", int(userID)).Interface("request_body", message).Msg("UpdateMessage request body")

	// Pass userID to the service for authorization
	updatedMessage, err := h.messageService.UpdateMessage(c.Request.Context(), int(userID), &message)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("message_id", id).Msg("User forbidden from updating message")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("message_id", id).Msg("Message not found for update")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("message_id", id).Msg("Failed to update message via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	if updatedMessage == nil {
		h.log.Info().Int("message_id", id).Msg("Message not found for update (service returned nil)")
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	h.log.Info().Int("message_id", updatedMessage.ID).Int("user_id", int(userID)).Msg("Message updated successfully")
	c.JSON(http.StatusOK, updatedMessage)
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	h.log.Info().Msg("Handling DeleteMessage request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in DeleteMessage")
		return // GetUserIDFromContext already handles the error response
	}

	idStr := c.Param("messageID")
	h.log.Debug().Str("messageID_param", idStr).Msg("Parsing message ID for deletion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("messageID_param", idStr).Msg("Invalid message ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	h.log.Debug().Int("message_id", id).Int("user_id", int(userID)).Msg("Attempting to delete message")

	// Pass userID to the service for authorization
	err = h.messageService.DeleteMessage(c.Request.Context(), int(userID), id)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("message_id", id).Msg("User forbidden from deleting message")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("message_id", id).Msg("Message not found for deletion")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("message_id", id).Msg("Failed to delete message via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	h.log.Info().Int("message_id", id).Int("user_id", int(userID)).Msg("Message deleted successfully")
	c.JSON(http.StatusNoContent, nil)
}
