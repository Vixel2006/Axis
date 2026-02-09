package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils" // Import the utils package
	"github.com/gin-gonic/gin"
)

type ChannelHandler struct {
	channelService services.ChannelService
}

func NewChannelHandler(cs services.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		channelService: cs,
	}
}

func (h *ChannelHandler) CreateChannel(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel.CreatorID = int(userID) // Set the CreatorID from the context

	createdChannel, err := h.channelService.CreateChannel(c.Request.Context(), &channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
		return
	}

	c.JSON(http.StatusCreated, createdChannel)
}

func (h *ChannelHandler) GetChannelByID(c *gin.Context) {
	idStr := c.Param("channelID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	channel, err := h.channelService.GetChannelByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channel"})
		return
	}

	if channel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	c.JSON(http.StatusOK, channel)
}

func (h *ChannelHandler) GetChannelsForWorkspace(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	workspaceIDStr := c.Param("workspaceID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	channels, err := h.channelService.GetChannelsForWorkspace(c.Request.Context(), int(userID), workspaceID)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channels for workspace"})
		return
	}

	c.JSON(http.StatusOK, channels)
}

func (h *ChannelHandler) UpdateChannel(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	idStr := c.Param("channelID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel.ID = id // Ensure the ID from the URL is used

	// Pass userID to the service for authorization
	updatedChannel, err := h.channelService.UpdateChannel(c.Request.Context(), int(userID), &channel)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update channel"})
		return
	}

	if updatedChannel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	c.JSON(http.StatusOK, updatedChannel)
}

func (h *ChannelHandler) DeleteChannel(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	idStr := c.Param("channelID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	// Pass userID to the service for authorization
	err = h.channelService.DeleteChannel(c.Request.Context(), int(userID), id)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete channel"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
