package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"github.com/gin-gonic/gin"
)

type ReactionHandler struct {
	reactionService services.ReactionService
}

func NewReactionHandler(rs services.ReactionService) *ReactionHandler {
	return &ReactionHandler{
		reactionService: rs,
	}
}

func (h *ReactionHandler) AddReaction(c *gin.Context) {
	messageIDStr := c.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var reqBody struct {
		UserID int    `json:"user_id"`
		Emoji  string `json:"emoji"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reaction := &models.Reaction{
		MessageID: messageID,
		UserID:    reqBody.UserID,
		Emoji:     reqBody.Emoji,
	}

	addedReaction, err := h.reactionService.AddReaction(c.Request.Context(), reaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	c.JSON(http.StatusCreated, addedReaction)
}

func (h *ReactionHandler) RemoveReaction(c *gin.Context) {
	messageIDStr := c.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	emoji := c.Param("emoji")

	err = h.reactionService.RemoveReaction(c.Request.Context(), messageID, userID, emoji)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ReactionHandler) GetReactionsForMessage(c *gin.Context) {
	messageIDStr := c.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	reactions, err := h.reactionService.GetReactionsForMessage(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reactions for message"})
		return
	}

	c.JSON(http.StatusOK, reactions)
}
