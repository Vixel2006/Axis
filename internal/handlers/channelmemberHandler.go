package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ChannelMemberHandler struct {
	channelMemberService services.ChannelMemberService
	log                  zerolog.Logger
}

func NewChannelMemberHandler(cms services.ChannelMemberService, logger zerolog.Logger) *ChannelMemberHandler {
	return &ChannelMemberHandler{
		channelMemberService: cms,
		log:                  logger,
	}
}

func (h *ChannelMemberHandler) AddMemberToChannel(c *gin.Context) {
	h.log.Info().Msg("Handling AddMemberToChannel request")
	channelIDStr := c.Param("channelID")
	h.log.Debug().Str("channelID_param", channelIDStr).Msg("Parsing channel ID")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("channelID_param", channelIDStr).Msg("Invalid channel ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var reqBody struct {
		UserID int `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for AddMemberToChannel")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channelMember, err := h.channelMemberService.AddMemberToChannel(c.Request.Context(), channelID, reqBody.UserID)
	if err != nil {
		h.log.Error().Err(err).Int("channel_id", channelID).Int("user_id", reqBody.UserID).Msg("Failed to add member to channel via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member to channel"})
		return
	}

	h.log.Info().Int("channel_id", channelID).Int("user_id", reqBody.UserID).Msg("Member added to channel successfully")
	c.JSON(http.StatusCreated, channelMember)
}

func (h *ChannelMemberHandler) RemoveMemberFromChannel(c *gin.Context) {
	h.log.Info().Msg("Handling RemoveMemberFromChannel request")
	channelIDStr := c.Param("channelID")
	h.log.Debug().Str("channelID_param", channelIDStr).Msg("Parsing channel ID for removal")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("channelID_param", channelIDStr).Msg("Invalid channel ID format for removal")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	userIDStr := c.Param("userID")
	h.log.Debug().Str("userID_param", userIDStr).Msg("Parsing user ID for removal")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("userID_param", userIDStr).Msg("Invalid user ID format for removal")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.channelMemberService.RemoveMemberFromChannel(c.Request.Context(), channelID, userID)
	if err != nil {
		h.log.Error().Err(err).Int("channel_id", channelID).Int("user_id", userID).Msg("Failed to remove member from channel via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member from channel"})
		return
	}

	h.log.Info().Int("channel_id", channelID).Int("user_id", userID).Msg("Member removed from channel successfully")
	c.JSON(http.StatusNoContent, nil)
}

func (h *ChannelMemberHandler) GetChannelMembers(c *gin.Context) {
	h.log.Info().Msg("Handling GetChannelMembers request")
	channelIDStr := c.Param("channelID")
	h.log.Debug().Str("channelID_param", channelIDStr).Msg("Parsing channel ID")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("channelID_param", channelIDStr).Msg("Invalid channel ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	members, err := h.channelMemberService.GetChannelMembers(c.Request.Context(), channelID)
	if err != nil {
		h.log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to retrieve channel members via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channel members"})
		return
	}

	h.log.Info().Int("channel_id", channelID).Int("members_count", len(members)).Msg("Channel members retrieved successfully")
	c.JSON(http.StatusOK, members)
}
