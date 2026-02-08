package services

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog/log"
)

type ReactionService interface {
	AddReaction(ctx context.Context, reaction *models.Reaction) (*models.Reaction, error)
	RemoveReaction(ctx context.Context, messageID, userID int, emoji string) error
	GetReactionsForMessage(ctx context.Context, messageID int) ([]models.Reaction, error)
}

type reactionService struct {
	reactionRepo repositories.ReactionRepo
}

func NewReactionService(rr repositories.ReactionRepo) ReactionService {
	return &reactionService{
		reactionRepo: rr,
	}
}

func (s *reactionService) AddReaction(ctx context.Context, reaction *models.Reaction) (*models.Reaction, error) {
	existingReaction, err := s.reactionRepo.GetReactionByMessageUserEmoji(ctx, reaction.MessageID, reaction.UserID, reaction.Emoji)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Err(err).Int("message_id", reaction.MessageID).Int("user_id", reaction.UserID).Str("emoji", reaction.Emoji).Msg("Failed to check for existing reaction")
		return nil, err
	}
	if existingReaction != nil {
		log.Info().Int("message_id", reaction.MessageID).Int("user_id", reaction.UserID).Str("emoji", reaction.Emoji).Msg("Reaction already exists")
		return existingReaction, nil // Or return an error indicating duplicate
	}

	err = s.reactionRepo.AddReaction(ctx, reaction)
	if err != nil {
		log.Error().Err(err).Int("message_id", reaction.MessageID).Int("user_id", reaction.UserID).Str("emoji", reaction.Emoji).Msg("Failed to add reaction")
		return nil, err
	}
	log.Info().Int("reaction_id", reaction.ID).Int("message_id", reaction.MessageID).Int("user_id", reaction.UserID).Str("emoji", reaction.Emoji).Msg("Reaction added successfully")
	return reaction, nil
}

func (s *reactionService) RemoveReaction(ctx context.Context, messageID, userID int, emoji string) error {
	reaction, err := s.reactionRepo.GetReactionByMessageUserEmoji(ctx, messageID, userID, emoji)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Reaction not found for removal")
			return nil // If not found, nothing to remove, so consider it successful
		}
		log.Error().Err(err).Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Failed to get reaction for removal")
		return err
	}

	err = s.reactionRepo.RemoveReaction(ctx, reaction.ID)
	if err != nil {
		log.Error().Err(err).Int("reaction_id", reaction.ID).Msg("Failed to remove reaction")
		return err
	}
	log.Info().Int("reaction_id", reaction.ID).Msg("Reaction removed successfully")
	return nil
}

func (s *reactionService) GetReactionsForMessage(ctx context.Context, messageID int) ([]models.Reaction, error) {
	reactions, err := s.reactionRepo.GetReactionsByMessageID(ctx, messageID)
	if err != nil {
		log.Error().Err(err).Int("message_id", messageID).Msg("Failed to get reactions for message")
		return nil, err
	}
	return reactions, nil
}
