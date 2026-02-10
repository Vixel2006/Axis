package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ChannelHandler struct {
	channelService services.ChannelService
	log            zerolog.Logger
}

func NewChannelHandler(cs services.ChannelService, logger zerolog.Logger) *ChannelHandler {
	return &ChannelHandler{
		channelService: cs,
		log:            logger,
	}
}

func (h *ChannelHandler) CreateChannel(c *gin.Context) {
	h.log.Info().Msg("Handling CreateChannel request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in CreateChannel")
		return
	}

	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for CreateChannel")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel.CreatorID = int(userID)

	createdChannel, err := h.channelService.CreateChannel(c.Request.Context(), &channel)
	if err != nil {
		h.log.Error().Err(err).Int("creator_id", int(userID)).Str("channel_name", channel.Name).Msg("Failed to create channel via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
		return
	}

	h.log.Info().Int("channel_id", createdChannel.ID).Int("creator_id", int(userID)).Msg("Channel created successfully")
	c.JSON(http.StatusCreated, createdChannel)
}

func (h *ChannelHandler) GetChannelByID(c *gin.Context) {
	h.log.Info().Msg("Handling GetChannelByID request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in GetChannelByID")
		return
	}

	idStr := c.Param("channelID")
	h.log.Debug().Str("channelID_param", idStr).Msg("Parsing channel ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("channelID_param", idStr).Msg("Invalid channel ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	channel, err := h.channelService.GetChannelByIDAuthorized(c.Request.Context(), int(userID), id)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("channel_id", id).Msg("User forbidden from accessing channel")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("channel_id", id).Msg("Failed to retrieve channel via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channel"})
		return
	}

	if channel == nil {
		h.log.Info().Int("channel_id", id).Msg("Channel not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	h.log.Info().Int("channel_id", id).Int("user_id", int(userID)).Msg("Channel retrieved successfully")
	c.JSON(http.StatusOK, channel)
}

func (h *ChannelHandler) GetChannelsForWorkspace(c *gin.Context) {
	h.log.Info().Msg("Handling GetChannelsForWorkspace request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in GetChannelsForWorkspace")
		return
	}

	workspaceIDStr := c.Param("workspaceID")
	h.log.Debug().Str("workspaceID_param", workspaceIDStr).Msg("Parsing workspace ID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("workspaceID_param", workspaceIDStr).Msg("Invalid workspace ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	channels, err := h.channelService.GetChannelsForWorkspace(c.Request.Context(), int(userID), workspaceID)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("workspace_id", workspaceID).Msg("User forbidden from accessing channels in workspace")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("workspace_id", workspaceID).Msg("Failed to retrieve channels for workspace via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channels for workspace"})
		return
	}

	h.log.Info().Int("user_id", int(userID)).Int("workspace_id", workspaceID).Int("channels_count", len(channels)).Msg("Channels for workspace retrieved successfully")
	c.JSON(http.StatusOK, channels)
}

func (h *ChannelHandler) UpdateChannel(c *gin.Context) {
	h.log.Info().Msg("Handling UpdateChannel request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in UpdateChannel")
		return
	}

	idStr := c.Param("channelID")
	h.log.Debug().Str("channelID_param", idStr).Msg("Parsing channel ID for update")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("channelID_param", idStr).Msg("Invalid channel ID format for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for UpdateChannel")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel.ID = id

	updatedChannel, err := h.channelService.UpdateChannel(c.Request.Context(), int(userID), &channel)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("channel_id", id).Msg("User forbidden from updating channel")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("channel_id", id).Msg("Failed to update channel via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update channel"})
		return
	}

	if updatedChannel == nil {
		h.log.Info().Int("channel_id", id).Msg("Channel not found for update")
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	h.log.Info().Int("channel_id", updatedChannel.ID).Int("user_id", int(userID)).Msg("Channel updated successfully")
	c.JSON(http.StatusOK, updatedChannel)
}

func (h *ChannelHandler) DeleteChannel(c *gin.Context) {
	h.log.Info().Msg("Handling DeleteChannel request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in DeleteChannel")
		return
	}

	idStr := c.Param("channelID")
	h.log.Debug().Str("channelID_param", idStr).Msg("Parsing channel ID for deletion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("channelID_param", idStr).Msg("Invalid channel ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	err = h.channelService.DeleteChannel(c.Request.Context(), int(userID), id)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("channel_id", id).Msg("User forbidden from deleting channel")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("channel_id", id).Msg("Failed to delete channel via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete channel"})
		return
	}

	h.log.Info().Int("channel_id", id).Int("user_id", int(userID)).Msg("Channel deleted successfully")
	c.JSON(http.StatusNoContent, nil)
}
