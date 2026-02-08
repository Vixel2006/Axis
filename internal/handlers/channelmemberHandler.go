package handlers

import (
	"net/http"
	"strconv"


	"axis/internal/services"
	"github.com/gin-gonic/gin"
)

type ChannelMemberHandler struct {
	channelMemberService services.ChannelMemberService
}

func NewChannelMemberHandler(cms services.ChannelMemberService) *ChannelMemberHandler {
	return &ChannelMemberHandler{
		channelMemberService: cms,
	}
}

func (h *ChannelMemberHandler) AddMemberToChannel(c *gin.Context) {
	channelIDStr := c.Param("channelID")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var reqBody struct {
		UserID int `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channelMember, err := h.channelMemberService.AddMemberToChannel(c.Request.Context(), channelID, reqBody.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member to channel"})
		return
	}

	c.JSON(http.StatusCreated, channelMember)
}

func (h *ChannelMemberHandler) RemoveMemberFromChannel(c *gin.Context) {
	channelIDStr := c.Param("channelID")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.channelMemberService.RemoveMemberFromChannel(c.Request.Context(), channelID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member from channel"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ChannelMemberHandler) GetChannelMembers(c *gin.Context) {
	channelIDStr := c.Param("channelID")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	members, err := h.channelMemberService.GetChannelMembers(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve channel members"})
		return
	}

	c.JSON(http.StatusOK, members)
}
