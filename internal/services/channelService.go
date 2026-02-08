package services

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog/log"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, channel *models.Channel) (*models.Channel, error)
	GetChannelByID(ctx context.Context, id int) (*models.Channel, error)
	GetChannelsForWorkspace(ctx context.Context, workspaceID int) ([]models.Channel, error)
	UpdateChannel(ctx context.Context, channel *models.Channel) (*models.Channel, error)
	DeleteChannel(ctx context.Context, id int) error
}

type channelService struct {
	channelRepo repositories.ChannelRepo
}

func NewChannelService(cr repositories.ChannelRepo) ChannelService {
	return &channelService{
		channelRepo: cr,
	}
}

func (s *channelService) CreateChannel(ctx context.Context, channel *models.Channel) (*models.Channel, error) {
	err := s.channelRepo.CreateChannel(ctx, channel)
	if err != nil {
		log.Error().Err(err).Str("channel_name", channel.Name).Msg("Failed to create channel")
		return nil, err
	}
	log.Info().Int("channel_id", channel.ID).Str("channel_name", channel.Name).Msg("Channel created successfully")
	return channel, nil
}

func (s *channelService) GetChannelByID(ctx context.Context, id int) (*models.Channel, error) {
	channel, err := s.channelRepo.GetChannelByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("channel_id", id).Msg("Channel not found")
			return nil, nil
		}
		log.Error().Err(err).Int("channel_id", id).Msg("Failed to get channel by ID")
		return nil, err
	}
	return channel, nil
}

func (s *channelService) GetChannelsForWorkspace(ctx context.Context, workspaceID int) ([]models.Channel, error) {
	channels, err := s.channelRepo.GetChannelsByWorkspaceID(ctx, workspaceID)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to get channels for workspace")
		return nil, err
	}
	return channels, nil
}

func (s *channelService) UpdateChannel(ctx context.Context, channel *models.Channel) (*models.Channel, error) {
	// First, check if the channel exists
	existingChannel, err := s.channelRepo.GetChannelByID(ctx, channel.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("channel_id", channel.ID).Msg("Channel not found for update")
			return nil, nil
		}
		log.Error().Err(err).Int("channel_id", channel.ID).Msg("Failed to get channel for update")
		return nil, err
	}

	// Update fields
	existingChannel.Name = channel.Name // Assuming Name is always provided for update
	existingChannel.Description = channel.Description

	err = s.channelRepo.UpdateChannel(ctx, existingChannel)
	if err != nil {
		log.Error().Err(err).Int("channel_id", existingChannel.ID).Msg("Failed to update channel")
		return nil, err
	}
	log.Info().Int("channel_id", existingChannel.ID).Msg("Channel updated successfully")
	return existingChannel, nil
}

func (s *channelService) DeleteChannel(ctx context.Context, id int) error {
	err := s.channelRepo.DeleteChannel(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("channel_id", id).Msg("Channel not found for deletion")
			return nil // Consider returning nil if not found is not an error for deletion
		}
		log.Error().Err(err).Int("channel_id", id).Msg("Failed to delete channel")
		return err
	}
	log.Info().Int("channel_id", id).Msg("Channel deleted successfully")
	return nil
}
