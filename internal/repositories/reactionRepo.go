package repositories

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type ReactionRepo interface {
	CreateReaction(ctx context.Context, reaction *models.Reaction) error
	DeleteReaction(ctx context.Context, messageID, userID int, emoji string) error
	GetReactionsByMessageID(ctx context.Context, messageID int) ([]models.Reaction, error)
	GetReactionByMessageUserEmoji(ctx context.Context, messageID, userID int, emoji string) (*models.Reaction, error)
}

type reactionRepository struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewReactionRepo(db *bun.DB, logger zerolog.Logger) ReactionRepo {
	return &reactionRepository{
		db:  db,
		log: logger,
	}
}

func (rr *reactionRepository) CreateReaction(ctx context.Context, reaction *models.Reaction) error {
	_, err := rr.db.NewInsert().Model(reaction).Exec(ctx)
	if err != nil {
		rr.log.Error().Err(err).Int("message_id", reaction.MessageID).Int("user_id", reaction.UserID).
			Str("emoji", reaction.Emoji).Msg("Failed to create reaction")
		return err
	}
	return nil
}

func (rr *reactionRepository) DeleteReaction(ctx context.Context, messageID, userID int, emoji string) error {
	_, err := rr.db.NewDelete().
		Model(&models.Reaction{}).
		Where("message_id = ?", messageID).
		Where("user_id = ?", userID).
		Where("emoji = ?", emoji).
		Exec(ctx)
	if err != nil {
		rr.log.Error().Err(err).Int("message_id", messageID).Int("user_id", userID).
			Str("emoji", emoji).Msg("Failed to delete reaction")
		return err
	}
	return nil
}

func (rr *reactionRepository) GetReactionsByMessageID(ctx context.Context, messageID int) ([]models.Reaction, error) {
	var reactions []models.Reaction
	err := rr.db.NewSelect().Model(&reactions).Where("message_id = ?", messageID).Scan(ctx)
	if err != nil {
		rr.log.Error().Err(err).Int("message_id", messageID).Msg("Failed to get reactions by message ID")
		return nil, err
	}
	return reactions, nil
}

func (rr *reactionRepository) GetReactionByMessageUserEmoji(ctx context.Context, messageID, userID int, emoji string) (*models.Reaction, error) {
	reaction := new(models.Reaction)
	err := rr.db.NewSelect().
		Model(reaction).
		Where("message_id = ?", messageID).
		Where("user_id = ?", userID).
		Where("emoji = ?", emoji).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			rr.log.Info().Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Reaction not found")
			return nil, nil
		}
		rr.log.Error().Err(err).Int("message_id", messageID).Int("user_id", userID).
			Str("emoji", emoji).Msg("Failed to get reaction by message, user, and emoji")
		return nil, err
	}
	return reaction, nil
}
