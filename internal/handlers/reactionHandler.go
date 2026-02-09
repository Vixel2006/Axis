package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ReactionHandler struct {
	reactionService services.ReactionService
	log             zerolog.Logger
}

func NewReactionHandler(rs services.ReactionService, logger zerolog.Logger) *ReactionHandler {
	return &ReactionHandler{
		reactionService: rs,
		log:             logger,
	}
}

func (h *ReactionHandler) AddReaction(c *gin.Context) {
	h.log.Info().Msg("Handling AddReaction request")
	messageIDStr := c.Param("messageID")
	h.log.Debug().Str("messageID_param", messageIDStr).Msg("Parsing message ID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("messageID_param", messageIDStr).Msg("Invalid message ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var reqBody struct {
		UserID int    `json:"user_id"`
		Emoji  string `json:"emoji"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for AddReaction")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Int("message_id", messageID).Int("user_id", reqBody.UserID).Str("emoji", reqBody.Emoji).Msg("AddReaction request body")

	reaction := &models.Reaction{
		MessageID: messageID,
		UserID:    reqBody.UserID,
		Emoji:     reqBody.Emoji,
	}

	addedReaction, err := h.reactionService.AddReaction(c.Request.Context(), reaction)
	if err != nil {
		h.log.Error().Err(err).Int("message_id", messageID).Int("user_id", reqBody.UserID).Str("emoji", reqBody.Emoji).Msg("Failed to add reaction via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}

	h.log.Info().Int("reaction_id", addedReaction.ID).Int("message_id", messageID).Int("user_id", reqBody.UserID).Msg("Reaction added successfully")
	c.JSON(http.StatusCreated, addedReaction)
}

func (h *ReactionHandler) RemoveReaction(c *gin.Context) {
	h.log.Info().Msg("Handling RemoveReaction request")
	messageIDStr := c.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("messageID_param", messageIDStr).Msg("Invalid message ID format for removal")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("userID_param", userIDStr).Msg("Invalid user ID format for removal")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	emoji := c.Param("emoji")
	h.log.Debug().Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Attempting to remove reaction")

	err = h.reactionService.RemoveReaction(c.Request.Context(), messageID, userID, emoji)
	if err != nil {
		h.log.Error().Err(err).Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Failed to remove reaction via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
		return
	}

	h.log.Info().Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Reaction removed successfully")
	c.JSON(http.StatusNoContent, nil)
}

func (h *ReactionHandler) GetReactionsForMessage(c *gin.Context) {
	h.log.Info().Msg("Handling GetReactionsForMessage request")
	messageIDStr := c.Param("messageID")
	h.log.Debug().Str("messageID_param", messageIDStr).Msg("Parsing message ID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("messageID_param", messageIDStr).Msg("Invalid message ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	reactions, err := h.reactionService.GetReactionsForMessage(c.Request.Context(), messageID)
	if err != nil {
		h.log.Error().Err(err).Int("message_id", messageID).Msg("Failed to retrieve reactions for message via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reactions for message"})
		return
	}

	h.log.Info().Int("message_id", messageID).Int("reactions_count", len(reactions)).Msg("Reactions for message retrieved successfully")
	c.JSON(http.StatusOK, reactions)
}
