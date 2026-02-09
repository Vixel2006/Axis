package repositories

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type MeetingRepo interface {
	CreateMeeting(ctx context.Context, meeting *models.Meeting) error
	GetMeetingByID(ctx context.Context, meetingID int) (*models.Meeting, error)
	GetMeetingsByChannelID(ctx context.Context, channelID int) ([]models.Meeting, error)
	UpdateMeeting(ctx context.Context, meeting *models.Meeting) error
	DeleteMeeting(ctx context.Context, meetingID int) error
	AddParticipantToMeeting(ctx context.Context, meetingID, userID int) error
	RemoveParticipantFromMeeting(ctx context.Context, meetingID, userID int) error
	IsParticipantInMeeting(ctx context.Context, meetingID, userID int) (bool, error)
}

type meetingRepository struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewMeetingRepo(db *bun.DB, logger zerolog.Logger) MeetingRepo {
	return &meetingRepository{
		db:  db,
		log: logger,
	}
}

func (mr *meetingRepository) CreateMeeting(ctx context.Context, meeting *models.Meeting) error {
	_, err := mr.db.NewInsert().Model(meeting).Exec(ctx)
	if err != nil {
		mr.log.Error().Err(err).Str("meeting_name", meeting.Name).Msg("Failed to create meeting")
		return err
	}
	return nil
}

func (mr *meetingRepository) GetMeetingByID(ctx context.Context, meetingID int) (*models.Meeting, error) {
	meeting := new(models.Meeting)
	err := mr.db.NewSelect().
		Model(meeting).
		Where("m.id = ?", meetingID).
		Relation("Creator").
		Relation("Channel").
		Relation("Participants"). // Load participants through the m2m relation
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			mr.log.Info().Int("meeting_id", meetingID).Msg("Meeting not found")
			return nil, nil
		}
		mr.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to get meeting by ID")
		return nil, err
	}
	return meeting, nil
}

func (mr *meetingRepository) GetMeetingsByChannelID(ctx context.Context, channelID int) ([]models.Meeting, error) {
	var meetings []models.Meeting
	err := mr.db.NewSelect().
		Model(&meetings).
		Where("m.channel_id = ?", channelID).
		Relation("Creator").
		Relation("Channel").
		Relation("Participants").
		Scan(ctx)
	if err != nil {
		mr.log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to get meetings by channel ID")
		return nil, err
	}
	return meetings, nil
}

func (mr *meetingRepository) UpdateMeeting(ctx context.Context, meeting *models.Meeting) error {
	_, err := mr.db.NewUpdate().Model(meeting).WherePK().Exec(ctx)
	if err != nil {
		mr.log.Error().Err(err).Int("meeting_id", meeting.ID).Msg("Failed to update meeting")
		return err
	}
	return nil
}

func (mr *meetingRepository) DeleteMeeting(ctx context.Context, meetingID int) error {
	_, err := mr.db.NewDelete().Model(&models.Meeting{}).Where("id = ?", meetingID).Exec(ctx)
	if err != nil {
		mr.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to delete meeting")
		return err
	}
	return nil
}

func (mr *meetingRepository) AddParticipantToMeeting(ctx context.Context, meetingID, userID int) error {
	meetingMember := &models.MeetingMember{
		MeetingID: meetingID,
		UserID:    userID,
	}
	_, err := mr.db.NewInsert().Model(meetingMember).Exec(ctx)
	if err != nil {
		mr.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", userID).Msg("Failed to add participant to meeting")
		return err
	}
	return nil
}

func (mr *meetingRepository) RemoveParticipantFromMeeting(ctx context.Context, meetingID, userID int) error {
	_, err := mr.db.NewDelete().
		Model(&models.MeetingMember{}).
		Where("meeting_id = ?", meetingID).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		mr.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", userID).Msg("Failed to remove participant from meeting")
		return err
	}
	return nil
}

func (mr *meetingRepository) IsParticipantInMeeting(ctx context.Context, meetingID, userID int) (bool, error) {
	count, err := mr.db.NewSelect().
		Model((*models.MeetingMember)(nil)).
		Where("meeting_id = ?", meetingID).
		Where("user_id = ?", userID).
		Count(ctx)
	if err != nil {
		mr.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", userID).Msg("Failed to check if user is participant in meeting")
		return false, err
	}
	return count > 0, nil
}
