package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	attachmentService services.AttachmentService
}

func NewAttachmentHandler(as services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{
		attachmentService: as,
	}
}

func (h *AttachmentHandler) CreateAttachment(c *gin.Context) {
	var attachment models.Attachment
	if err := c.ShouldBindJSON(&attachment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAttachment, err := h.attachmentService.CreateAttachment(c.Request.Context(), &attachment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attachment"})
		return
	}

	c.JSON(http.StatusCreated, createdAttachment)
}

func (h *AttachmentHandler) GetAttachmentByID(c *gin.Context) {
	idStr := c.Param("attachmentID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
		return
	}

	attachment, err := h.attachmentService.GetAttachmentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attachment"})
		return
	}

	if attachment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	c.JSON(http.StatusOK, attachment)
}

func (h *AttachmentHandler) GetAttachmentsForMessage(c *gin.Context) {
	messageIDStr := c.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	attachments, err := h.attachmentService.GetAttachmentsForMessage(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attachments for message"})
		return
	}

	c.JSON(http.StatusOK, attachments)
}
