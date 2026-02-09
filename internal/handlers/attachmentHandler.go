package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AttachmentHandler struct {
	attachmentService services.AttachmentService
	log               zerolog.Logger
}

func NewAttachmentHandler(as services.AttachmentService, logger zerolog.Logger) *AttachmentHandler {
	return &AttachmentHandler{
		attachmentService: as,
		log:               logger,
	}
}

func (h *AttachmentHandler) CreateAttachment(c *gin.Context) {
	h.log.Info().Msg("Handling CreateAttachment request")
	var attachment models.Attachment
	if err := c.ShouldBindJSON(&attachment); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for CreateAttachment")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAttachment, err := h.attachmentService.CreateAttachment(c.Request.Context(), &attachment)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to create attachment via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attachment"})
		return
	}

	h.log.Info().Int("attachment_id", createdAttachment.ID).Msg("Attachment created successfully")
	c.JSON(http.StatusCreated, createdAttachment)
}

func (h *AttachmentHandler) GetAttachmentByID(c *gin.Context) {
	idStr := c.Param("attachmentID")
	h.log.Info().Str("handler", "GetAttachmentByID").Str("attachmentID_param", idStr).Msg("Received request")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("attachmentID_param", idStr).Msg("Invalid attachment ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
		return
	}

	attachment, err := h.attachmentService.GetAttachmentByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error().Err(err).Int("attachment_id", id).Msg("Failed to retrieve attachment via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attachment"})
		return
	}

	if attachment == nil {
		h.log.Info().Int("attachment_id", id).Msg("Attachment not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	h.log.Info().Int("attachment_id", id).Msg("Attachment retrieved successfully")
	c.JSON(http.StatusOK, attachment)
}

func (h *AttachmentHandler) GetAttachmentsForMessage(c *gin.Context) {
	messageIDStr := c.Param("messageID")
	h.log.Info().Str("handler", "GetAttachmentsForMessage").Str("messageID_param", messageIDStr).Msg("Received request")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("messageID_param", messageIDStr).Msg("Invalid message ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	attachments, err := h.attachmentService.GetAttachmentsForMessage(c.Request.Context(), messageID)
	if err != nil {
		h.log.Error().Err(err).Int("message_id", messageID).Msg("Failed to retrieve attachments for message via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attachments for message"})
		return
	}

	h.log.Info().Int("message_id", messageID).Int("attachments_count", len(attachments)).Msg("Attachments for message retrieved successfully")
	c.JSON(http.StatusOK, attachments)
}
