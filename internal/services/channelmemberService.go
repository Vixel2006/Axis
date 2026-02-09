package services

import (
	"context"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog"
)

type ChannelMemberService interface {
	AddMemberToChannel(ctx context.Context, channelID, userID int) (*models.ChannelMember, error)
	RemoveMemberFromChannel(ctx context.Context, channelID, userID int) error
	GetChannelMembers(ctx context.Context, channelID int) ([]models.ChannelMember, error)
}

type channelMemberService struct {
	channelMemberRepo repositories.ChannelMemberRepo
	log               zerolog.Logger
}

func NewChannelMemberService(cmr repositories.ChannelMemberRepo, logger zerolog.Logger) ChannelMemberService {
	return &channelMemberService{
		channelMemberRepo: cmr,
		log:               logger,
	}
}

func (s *channelMemberService) AddMemberToChannel(ctx context.Context, channelID, userID int) (*models.ChannelMember, error) {
	err := s.channelMemberRepo.AddMemberToChannel(ctx, channelID, userID)
	if err != nil {
		s.log.Error().Err(err).Int("channel_id", channelID).Int("user_id", userID).Msg("Failed to add member to channel")
		return nil, err
	}
	channelMember := &models.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
	}
	s.log.Info().Int("channel_id", channelID).Int("user_id", userID).Msg("Member added to channel successfully")
	return channelMember, nil
}

func (s *channelMemberService) RemoveMemberFromChannel(ctx context.Context, channelID, userID int) error {
	err := s.channelMemberRepo.RemoveMemberFromChannel(ctx, channelID, userID)
	if err != nil {
		s.log.Error().Err(err).Int("channel_id", channelID).Int("user_id", userID).Msg("Failed to remove member from channel")
		return err
	}
	s.log.Info().Int("channel_id", channelID).Int("user_id", userID).Msg("Member removed from channel successfully")
	return nil
}

func (s *channelMemberService) GetChannelMembers(ctx context.Context, channelID int) ([]models.ChannelMember, error) {
	s.log.Debug().Int("channel_id", channelID).Msg("Calling ChannelMemberRepo.GetChannelMembers")
	members, err := s.channelMemberRepo.GetChannelMembers(ctx, channelID)
	if err != nil {
		s.log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to get channel members")
		return nil, err
	}
	s.log.Debug().Int("channel_id", channelID).Int("members_count", len(members)).Msg("Received channel members from repo")
	return members, nil
}
