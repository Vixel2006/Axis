package repositories

import (
	"context"
	"database/sql"

	"backend/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ChannelRepo interface {
	CreateChannel(ctx context.Context, channel *models.Channel) error
	GetChannelByID(ctx context.Context, channelID int) (*models.Channel, error)
	GetChannelsByWorkspaceID(ctx context.Context, workspaceID int) ([]models.Channel, error)
	UpdateChannel(ctx context.Context, channel *models.Channel) error
	DeleteChannel(ctx context.Context, channelID int) error
}

type channelRepository struct {
	db *bun.DB
}

func NewChannelRepo(db *bun.DB) ChannelRepo {
	return &channelRepository{
		db: db,
	}
}

func (cr *channelRepository) CreateChannel(ctx context.Context, channel *models.Channel) error {
	_, err := cr.db.NewInsert().Model(channel).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Str("channel_name", channel.Name).Msg("Failed to create channel")
		return err
	}
	return nil
}

func (cr *channelRepository) GetChannelByID(ctx context.Context, channelID int) (*models.Channel, error) {
	channel := new(models.Channel)
	err := cr.db.NewSelect().Model(channel).Where("id = ?", channelID).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("channel_id", channelID).Msg("Channel not found")
			return nil, nil
		}
		log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to get channel by ID")
		return nil, err
	}
	return channel, nil
}

func (cr *channelRepository) GetChannelsByWorkspaceID(ctx context.Context, workspaceID int) ([]models.Channel, error) {
	var channels []models.Channel
	err := cr.db.NewSelect().Model(&channels).Where("workspace_id = ?", workspaceID).Scan(ctx)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to get channels by workspace ID")
		return nil, err
	}
	return channels, nil
}

func (cr *channelRepository) UpdateChannel(ctx context.Context, channel *models.Channel) error {
	_, err := cr.db.NewUpdate().Model(channel).WherePK().Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("channel_id", channel.ID).Msg("Failed to update channel")
		return err
	}
	return nil
}

func (cr *channelRepository) DeleteChannel(ctx context.Context, channelID int) error {
	_, err := cr.db.NewDelete().Model(&models.Channel{}).Where("id = ?", channelID).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to delete channel")
		return err
	}
	return nil
}
