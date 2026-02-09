package repositories

import (
	"context"

	"axis/internal/models"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type ChannelMemberRepo interface {
	AddMemberToChannel(ctx context.Context, channelID, userID int) error
	RemoveMemberFromChannel(ctx context.Context, channelID, userID int) error
	GetChannelMembers(ctx context.Context, channelID int) ([]models.ChannelMember, error)
	GetChannelsForUser(ctx context.Context, userID int) ([]models.ChannelMember, error)
	UpdateLastReadMessageID(ctx context.Context, channelID, userID int, messageID *int) error
	IsMemberOfChannel(ctx context.Context, channelID, userID int) (bool, error)
}

type channelMemberRepository struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewChannelMemberRepo(db *bun.DB, logger zerolog.Logger) ChannelMemberRepo {
	return &channelMemberRepository{
		db:  db,
		log: logger,
	}
}

func (cmr *channelMemberRepository) AddMemberToChannel(ctx context.Context, channelID, userID int) error {
	channelMember := &models.ChannelMember{
		ChannelID: channelID,
		UserID:    userID,
	}
	_, err := cmr.db.NewInsert().Model(channelMember).Exec(ctx)
	if err != nil {
		cmr.log.Error().Err(err).Int("channel_id", channelID).Int("user_id", userID).Msg("Failed to add member to channel")
		return err
	}
	return nil
}

func (cmr *channelMemberRepository) RemoveMemberFromChannel(ctx context.Context, channelID, userID int) error {
	_, err := cmr.db.NewDelete().
		Model(&models.ChannelMember{}).
		Where("channel_id = ?", channelID).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		cmr.log.Error().Err(err).Int("channel_id", channelID).Int("user_id", userID).Msg("Failed to remove member from channel")
		return err
	}
	return nil
}

func (cmr *channelMemberRepository) GetChannelMembers(ctx context.Context, channelID int) ([]models.ChannelMember, error) {
	cmr.log.Debug().Int("channel_id", channelID).Msg("Fetching channel members from database")
	var members []models.ChannelMember
	err := cmr.db.NewSelect().
		Model(&members).
		Where("channel_id = ?", channelID).
		Relation("User"). // Eager load user details
		Scan(ctx)
	if err != nil {
		cmr.log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to get channel members")
		return nil, err
	}
	cmr.log.Debug().Int("channel_id", channelID).Int("members_count", len(members)).Msg("Successfully fetched channel members")
	return members, nil
}

func (cmr *channelMemberRepository) GetChannelsForUser(ctx context.Context, userID int) ([]models.ChannelMember, error) {
	var memberships []models.ChannelMember
	err := cmr.db.NewSelect().
		Model(&memberships).
		Where("user_id = ?", userID).
		Relation("Channel"). // Eager load channel details
		Scan(ctx)
	if err != nil {
		cmr.log.Error().Err(err).Int("user_id", userID).Msg("Failed to get channels for user")
		return nil, err
	}
	return memberships, nil
}

func (cmr *channelMemberRepository) UpdateLastReadMessageID(ctx context.Context, channelID, userID int, messageID *int) error {
	_, err := cmr.db.NewUpdate().
		Model(&models.ChannelMember{LastReadMessageID: messageID}).
		Where("channel_id = ?", channelID).
		Where("user_id = ?", userID).
		Set("last_read_message_id = ?", messageID).
		Exec(ctx)
	if err != nil {
		cmr.log.Error().Err(err).Int("channel_id", channelID).Int("user_id", userID).
			Interface("message_id", messageID).Msg("Failed to update last read message ID")
		return err
	}
	return nil
}

func (cmr *channelMemberRepository) IsMemberOfChannel(ctx context.Context, channelID, userID int) (bool, error) {
	count, err := cmr.db.NewSelect().
		Model((*models.ChannelMember)(nil)).
		Where("channel_id = ?", channelID).
		Where("user_id = ?", userID).
		Count(ctx)
	if err != nil {
		cmr.log.Error().Err(err).Int("channel_id", channelID).Int("user_id", userID).Msg("Failed to check if user is member of channel")
		return false, err
	}
	return count > 0, nil
}
