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
	GetChannelByIDAuthorized(ctx context.Context, userID, channelID int) (*models.Channel, error)
	GetChannelsForWorkspace(ctx context.Context, userID int, workspaceID int) ([]models.Channel, error)
	UpdateChannel(ctx context.Context, userID int, channel *models.Channel) (*models.Channel, error)
	DeleteChannel(ctx context.Context, userID int, id int) error
}

type channelService struct {
	channelRepo         repositories.ChannelRepo
	channelMemberRepo   repositories.ChannelMemberRepo
	workspaceMemberRepo repositories.WorkspaceMemberRepo
}

func NewChannelService(cr repositories.ChannelRepo, cmr repositories.ChannelMemberRepo, wmr repositories.WorkspaceMemberRepo) ChannelService {
	return &channelService{
		channelRepo:         cr,
		channelMemberRepo:   cmr,
		workspaceMemberRepo: wmr,
	}
}

func (s *channelService) CreateChannel(ctx context.Context, channel *models.Channel) (*models.Channel, error) {
	err := s.channelRepo.CreateChannel(ctx, channel)
	if err != nil {
		log.Error().Err(err).Str("channel_name", channel.Name).Msg("Failed to create channel")
		return nil, err
	}
	log.Info().Int("channel_id", channel.ID).Str("channel_name", channel.Name).Msg("Channel created successfully")

	// Add the creator as a member of the new channel
	err = s.channelMemberRepo.AddMemberToChannel(ctx, channel.ID, channel.CreatorID)
	if err != nil {
		log.Error().Err(err).Int("channel_id", channel.ID).Int("user_id", channel.CreatorID).Msg("Failed to add creator as member to channel")
		// Decide if this error should prevent channel creation or just be logged
		// For now, we'll return the error as it's a critical step
		return nil, err
	}
	log.Info().Int("channel_id", channel.ID).Int("user_id", channel.CreatorID).Msg("Creator added to channel members")

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

func (s *channelService) GetChannelByIDAuthorized(ctx context.Context, userID, channelID int) (*models.Channel, error) {
	channel, err := s.channelRepo.GetChannelByID(ctx, channelID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("channel_id", channelID).Msg("Channel not found")
			return nil, nil
		}
		log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to get channel by ID for authorization check")
		return nil, err
	}
	if channel == nil {
		return nil, nil // Channel not found
	}

	// First, check if the user is a member of the workspace this channel belongs to
	isWorkspaceMember, err := s.workspaceMemberRepo.IsMemberOfWorkspace(ctx, channel.WorkspaceID, userID)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", channel.WorkspaceID).Int("user_id", userID).Msg("Failed to check workspace membership for channel access")
		return nil, err
	}
	if isWorkspaceMember {
		return channel, nil // User is a member of the workspace, so they have access to this channel
	}

	// If not a workspace member, then check if they are a direct channel member (e.g., for private channels not tied to workspace general access)
	isChannelMember, err := s.channelMemberRepo.IsMemberOfChannel(ctx, channelID, userID)
	if err != nil {
		log.Error().Err(err).Int("channel_id", channelID).Int("user_id", userID).Msg("Failed to check direct channel membership")
		return nil, err
	}
	if isChannelMember {
		return channel, nil
	}

	return nil, &ForbiddenError{Message: "User not authorized to access this channel"}
}

func (s *channelService) GetChannelsForWorkspace(ctx context.Context, userID int, workspaceID int) ([]models.Channel, error) {
	// Authorization check: User must be a member of the workspace
	isMember, err := s.workspaceMemberRepo.IsMemberOfWorkspace(ctx, workspaceID, int(userID))
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", int(userID)).Msg("Failed to check workspace membership for channels")
		return nil, err
	}
	if !isMember {
		return nil, &ForbiddenError{Message: "User not authorized to view channels in this workspace"}
	}

	channels, err := s.channelRepo.GetChannelsByWorkspaceID(ctx, workspaceID)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to get channels for workspace")
		return nil, err
	}
	return channels, nil
}

func (s *channelService) UpdateChannel(ctx context.Context, userID int, channel *models.Channel) (*models.Channel, error) {
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

	// Authorization check: Only the creator of the channel or an admin of the workspace can update it
	// First, check if the user is the creator
	if existingChannel.CreatorID != int(userID) {
		// If not the creator, check if the user is an admin member of the workspace
		isWorkspaceAdmin, err := s.workspaceMemberRepo.IsMemberOfWorkspace(ctx, existingChannel.WorkspaceID, int(userID))
		if err != nil {
			log.Error().Err(err).Int("workspace_id", existingChannel.WorkspaceID).Int("user_id", int(userID)).Msg("Failed to check workspace membership for channel update")
			return nil, err
		}
		if !isWorkspaceAdmin {
			return nil, &ForbiddenError{Message: "User not authorized to update this channel"}
		}
		// Further check for admin role
		member, err := s.workspaceMemberRepo.GetWorkspaceMember(ctx, existingChannel.WorkspaceID, int(userID))
		if err != nil && err != sql.ErrNoRows {
			log.Error().Err(err).Int("workspace_id", existingChannel.WorkspaceID).Int("user_id", int(userID)).Msg("Failed to get workspace member role for channel update")
			return nil, err
		}
		if member == nil || member.Role != models.Admin {
			return nil, &ForbiddenError{Message: "User not authorized to update this channel"}
		}
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

func (s *channelService) DeleteChannel(ctx context.Context, userID int, id int) error {
	existingChannel, err := s.channelRepo.GetChannelByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("channel_id", id).Msg("Channel not found for deletion")
			return nil
		}
		log.Error().Err(err).Int("channel_id", id).Msg("Failed to get channel for deletion")
		return err
	}

	// Authorization check: Only the creator of the channel or an admin of the workspace can delete it
	// First, check if the user is the creator
	if existingChannel.CreatorID != int(userID) {
		// If not the creator, check if the user is an admin member of the workspace
		isWorkspaceAdmin, err := s.workspaceMemberRepo.IsMemberOfWorkspace(ctx, existingChannel.WorkspaceID, int(userID))
		if err != nil {
			log.Error().Err(err).Int("workspace_id", existingChannel.WorkspaceID).Int("user_id", int(userID)).Msg("Failed to check workspace membership for channel deletion")
			return err
		}
		if !isWorkspaceAdmin {
			return &ForbiddenError{Message: "User not authorized to delete this channel"}
		}
		// Further check for admin role
		member, err := s.workspaceMemberRepo.GetWorkspaceMember(ctx, existingChannel.WorkspaceID, int(userID))
		if err != nil && err != sql.ErrNoRows {
			log.Error().Err(err).Int("workspace_id", existingChannel.WorkspaceID).Int("user_id", int(userID)).Msg("Failed to get workspace member role for channel deletion")
			return err
		}
		if member == nil || member.Role != models.Admin {
			return &ForbiddenError{Message: "User not authorized to delete this channel"}
		}
	}

	err = s.channelRepo.DeleteChannel(ctx, id)
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
