package services

import (
	"context"
	"database/sql"
	"errors"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog"
)

type MeetingService interface {
	CreateMeeting(ctx context.Context, meeting *models.Meeting, creatorID int, participantIDs []int) (*models.Meeting, error)
	GetMeetingByID(ctx context.Context, id int) (*models.Meeting, error)
	GetMeetingByIDAuthorized(ctx context.Context, userID, meetingID int) (*models.Meeting, error)
	GetMeetingsByChannelID(ctx context.Context, channelID int) ([]models.Meeting, error)
	UpdateMeeting(ctx context.Context, meeting *models.Meeting, userID int) (*models.Meeting, error)
	DeleteMeeting(ctx context.Context, id int, userID int) error
	AddParticipant(ctx context.Context, meetingID, userID, participantID int) error
	RemoveParticipant(ctx context.Context, meetingID, userID, participantID int) error
}

type meetingService struct {
	meetingRepo       repositories.MeetingRepo
	channelRepo       repositories.ChannelRepo
	userRepo          repositories.UserRepo
	channelMemberRepo repositories.ChannelMemberRepo
	log               zerolog.Logger
}

func NewMeetingService(mr repositories.MeetingRepo, cr repositories.ChannelRepo, ur repositories.UserRepo, cmr repositories.ChannelMemberRepo, logger zerolog.Logger) MeetingService {
	return &meetingService{
		meetingRepo:       mr,
		channelRepo:       cr,
		userRepo:          ur,
		channelMemberRepo: cmr,
		log:               logger,
	}
}

func (s *meetingService) CreateMeeting(ctx context.Context, meeting *models.Meeting, creatorID int, participantIDs []int) (*models.Meeting, error) {
	channel, err := s.channelRepo.GetChannelByID(ctx, meeting.ChannelID)
	if err != nil {
		s.log.Error().Err(err).Int("channel_id", meeting.ChannelID).Msg("Failed to retrieve channel for meeting creation")
		return nil, err
	}
	if channel == nil {
		s.log.Warn().Int("channel_id", meeting.ChannelID).Msg("Channel not found for meeting creation")
		return nil, errors.New("channel not found")
	}

	meeting.CreatorID = creatorID
	err = s.meetingRepo.CreateMeeting(ctx, meeting)
	if err != nil {
		s.log.Error().Err(err).Str("meeting_name", meeting.Name).Msg("Failed to create meeting")
		return nil, err
	}
	s.log.Info().Int("meeting_id", meeting.ID).Str("meeting_name", meeting.Name).Msg("Meeting created successfully")

	err = s.meetingRepo.AddParticipantToMeeting(ctx, meeting.ID, creatorID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meeting.ID).Int("user_id", creatorID).Msg("Failed to add creator as participant to meeting")
		return nil, err
	}

	for _, pID := range participantIDs {
		user, err := s.userRepo.GetUserByID(ctx, pID)
		if err != nil {
			s.log.Error().Err(err).Int("user_id", pID).Msg("Participant user not found, skipping")
			continue
		}
		if user == nil {
			s.log.Warn().Int("user_id", pID).Msg("Participant user does not exist, skipping")
			continue
		}
		err = s.meetingRepo.AddParticipantToMeeting(ctx, meeting.ID, pID)
		if err != nil {
			s.log.Error().Err(err).Int("meeting_id", meeting.ID).Int("user_id", pID).Msg("Failed to add participant to meeting")
		}
	}

	return meeting, nil
}

func (s *meetingService) GetMeetingByID(ctx context.Context, id int) (*models.Meeting, error) {
	meeting, err := s.meetingRepo.GetMeetingByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("meeting_id", id).Msg("Meeting not found")
			return nil, nil
		}
		s.log.Error().Err(err).Int("meeting_id", id).Msg("Failed to get meeting by ID")
		return nil, err
	}
	return meeting, nil
}

func (s *meetingService) GetMeetingByIDAuthorized(ctx context.Context, userID, meetingID int) (*models.Meeting, error) {
	meeting, err := s.meetingRepo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("meeting_id", meetingID).Msg("Meeting not found")
			return nil, nil
		}
		s.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to get meeting by ID for authorization check")
		return nil, err
	}
	if meeting == nil {
		s.log.Info().Int("meeting_id", meetingID).Msg("Meeting not found during authorization check")
		return nil, nil
	}

	isParticipant, err := s.meetingRepo.IsParticipantInMeeting(ctx, meetingID, userID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", userID).Msg("Failed to check if user is participant in meeting")
		return nil, err
	}
	if isParticipant {
		s.log.Debug().Int("meeting_id", meetingID).Int("user_id", userID).Msg("User is a direct participant of the meeting")
		return meeting, nil
	}

	isChannelMember, err := s.channelMemberRepo.IsMemberOfChannel(ctx, meeting.ChannelID, userID)
	if err != nil {
		s.log.Error().Err(err).Int("channel_id", meeting.ChannelID).Int("user_id", userID).Msg("Failed to check if user is member of meeting's channel")
		return nil, err
	}
	if isChannelMember {
		s.log.Debug().Int("meeting_id", meetingID).Int("user_id", userID).Msg("User is a member of the meeting's channel")
		return meeting, nil
	}

	s.log.Warn().Int("meeting_id", meetingID).Int("user_id", userID).Msg("User not authorized to access this meeting")
	return nil, &ForbiddenError{Message: "User not authorized to access this meeting"}
}

func (s *meetingService) GetMeetingsByChannelID(ctx context.Context, channelID int) ([]models.Meeting, error) {
	meetings, err := s.meetingRepo.GetMeetingsByChannelID(ctx, channelID)
	if err != nil {
		s.log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to get meetings by channel ID")
		return nil, err
	}
	return meetings, nil
}

func (s *meetingService) UpdateMeeting(ctx context.Context, meeting *models.Meeting, userID int) (*models.Meeting, error) {
	existingMeeting, err := s.meetingRepo.GetMeetingByID(ctx, meeting.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("meeting_id", meeting.ID).Msg("Meeting not found for update")
			return nil, errors.New("meeting not found")
		}
		s.log.Error().Err(err).Int("meeting_id", meeting.ID).Msg("Failed to get meeting for update")
		return nil, err
	}

	if existingMeeting.CreatorID != userID {
		s.log.Warn().Int("meeting_id", meeting.ID).Int("user_id", userID).Msg("User not authorized to update this meeting")
		return nil, &ForbiddenError{Message: "User not authorized to update this meeting"}
	}

	existingMeeting.Name = meeting.Name
	existingMeeting.Description = meeting.Description
	existingMeeting.StartTime = meeting.StartTime
	existingMeeting.EndTime = meeting.EndTime
	existingMeeting.ChannelID = meeting.ChannelID

	err = s.meetingRepo.UpdateMeeting(ctx, existingMeeting)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", existingMeeting.ID).Msg("Failed to update meeting")
		return nil, err
	}
	s.log.Info().Int("meeting_id", existingMeeting.ID).Msg("Meeting updated successfully")
	return existingMeeting, nil
}

func (s *meetingService) DeleteMeeting(ctx context.Context, id int, userID int) error {
	existingMeeting, err := s.meetingRepo.GetMeetingByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("meeting_id", id).Msg("Meeting not found for deletion")
			return errors.New("meeting not found")
		}
		s.log.Error().Err(err).Int("meeting_id", id).Msg("Failed to get meeting for deletion")
		return err
	}

	if existingMeeting.CreatorID != userID {
		s.log.Warn().Int("meeting_id", id).Int("user_id", userID).Msg("User not authorized to delete this meeting")
		return &ForbiddenError{Message: "User not authorized to delete this meeting"}
	}

	err = s.meetingRepo.DeleteMeeting(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", id).Msg("Failed to delete meeting")
		return err
	}
	s.log.Info().Int("meeting_id", id).Msg("Meeting deleted successfully")
	return nil
}

func (s *meetingService) AddParticipant(ctx context.Context, meetingID, userID, participantID int) error {
	existingMeeting, err := s.meetingRepo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Warn().Int("meeting_id", meetingID).Msg("Meeting not found for adding participant")
			return errors.New("meeting not found")
		}
		s.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to get meeting for adding participant")
		return err
	}
	if existingMeeting.CreatorID != userID {
		s.log.Warn().Int("meeting_id", meetingID).Int("user_id", userID).Int("participant_id", participantID).Msg("User not authorized to add participants to this meeting")
		return &ForbiddenError{Message: "User not authorized to add participants to this meeting"}
	}

	participant, err := s.userRepo.GetUserByID(ctx, participantID)
	if err != nil {
		s.log.Error().Err(err).Int("user_id", participantID).Msg("Failed to retrieve participant user")
		return err
	}
	if participant == nil {
		s.log.Warn().Int("user_id", participantID).Msg("Participant user not found for adding to meeting")
		return errors.New("participant user not found")
	}

	err = s.meetingRepo.AddParticipantToMeeting(ctx, meetingID, participantID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", participantID).Msg("Failed to add participant to meeting")
	} else {
		s.log.Info().Int("meeting_id", meetingID).Int("user_id", participantID).Msg("Participant added to meeting successfully")
	}
	return err
}

func (s *meetingService) RemoveParticipant(ctx context.Context, meetingID, userID, participantID int) error {
	existingMeeting, err := s.meetingRepo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Warn().Int("meeting_id", meetingID).Msg("Meeting not found for removing participant")
			return errors.New("meeting not found")
		}
		s.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to get meeting for removing participant")
		return err
	}
	if existingMeeting.CreatorID != userID {
		s.log.Warn().Int("meeting_id", meetingID).Int("user_id", userID).Int("participant_id", participantID).Msg("User not authorized to remove participants from this meeting")
		return &ForbiddenError{Message: "User not authorized to remove participants from this meeting"}
	}

	err = s.meetingRepo.RemoveParticipantFromMeeting(ctx, meetingID, participantID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", participantID).Msg("Failed to remove participant from meeting")
	} else {
		s.log.Info().Int("meeting_id", meetingID).Int("user_id", participantID).Msg("Participant removed from meeting successfully")
	}
	return err
}
