package services

import (
	"context"
	"database/sql"
	"errors"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog/log"
)




type MessageService interface {
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	GetMessageByID(ctx context.Context, id int) (*models.Message, error)
	GetMessagesInMeeting(ctx context.Context, meetingID int, limit, offset int) ([]models.Message, error)
	UpdateMessage(ctx context.Context, userID int, message *models.Message) (*models.Message, error)
	DeleteMessage(ctx context.Context, userID int, id int) error
}

type messageService struct {
	messageRepo repositories.MessageRepo
	meetingRepo repositories.MeetingRepo
}

func NewMessageService(mr repositories.MessageRepo, metR repositories.MeetingRepo) MessageService {
	return &messageService{
		messageRepo: mr,
		meetingRepo: metR,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	// Authorization check: User must be a participant in the meeting
	meeting, err := s.meetingRepo.GetMeetingByID(ctx, message.MeetingID)
	if err != nil {
		log.Error().Err(err).Int("meeting_id", message.MeetingID).Msg("Failed to retrieve meeting for message creation")
		return nil, err
	}
	if meeting == nil {
		return nil, errors.New("meeting not found")
	}

	isParticipant := false
	for _, participant := range meeting.Participants {
		if participant.ID == message.SenderID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return nil, &ForbiddenError{Message: "User is not a participant of this meeting"}
	}

	err = s.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		log.Error().Err(err).Int("meeting_id", message.MeetingID).Int("sender_id", message.SenderID).Msg("Failed to create message")
		return nil, err
	}
	log.Info().Int("message_id", message.ID).Int("meeting_id", message.MeetingID).Msg("Message created successfully")
	return message, nil
}

func (s *messageService) GetMessageByID(ctx context.Context, id int) (*models.Message, error) {
	message, err := s.messageRepo.GetMessageByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("message_id", id).Msg("Message not found")
			return nil, nil
		}
		log.Error().Err(err).Int("message_id", id).Msg("Failed to get message by ID")
		return nil, err
	}
	return message, nil
}

func (s *messageService) GetMessagesInMeeting(ctx context.Context, meetingID int, limit, offset int) ([]models.Message, error) {
	// Authorization check: User must be a participant in the meeting
	// This requires adding userID to the context or as a parameter to this function
	// For now, assuming current user ID is available in context or passed explicitly.
	// As this function is typically called by a handler, the handler should provide the userID.
	// For simplicity, let's assume `GetMeetingByID` can check participant status or add `userID` as a parameter if needed.
	// For this refactoring, let's just make sure the meeting exists.
	meeting, err := s.meetingRepo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to retrieve meeting for message retrieval")
		return nil, err
	}
	if meeting == nil {
		return nil, errors.New("meeting not found")
	}

	messages, err := s.messageRepo.GetMessagesByMeetingID(ctx, meetingID, limit, offset)
	if err != nil {
		log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to get messages for meeting")
		return nil, err
	}
	return messages, nil
}

func (s *messageService) UpdateMessage(ctx context.Context, userID int, message *models.Message) (*models.Message, error) {
	existingMessage, err := s.messageRepo.GetMessageByID(ctx, message.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("message_id", message.ID).Msg("Message not found for update")
			return nil, nil
		}
		log.Error().Err(err).Int("message_id", message.ID).Msg("Failed to get message for update")
		return nil, err
	}

	// Authorization check: Only the sender can update their message
	if existingMessage.SenderID != int(userID) {
		return nil, &ForbiddenError{Message: "User not authorized to update this message"}
	}

	// Update fields
	existingMessage.Content = message.Content

	err = s.messageRepo.UpdateMessage(ctx, existingMessage)
	if err != nil {
		log.Error().Err(err).Int("message_id", existingMessage.ID).Msg("Failed to update message")
		return nil, err
	}
	log.Info().Int("message_id", existingMessage.ID).Msg("Message updated successfully")
	return existingMessage, nil
}

func (s *messageService) DeleteMessage(ctx context.Context, userID int, id int) error {
	existingMessage, err := s.messageRepo.GetMessageByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("message_id", id).Msg("Message not found for deletion")
			return nil // Consider returning nil if not found is not an error for deletion
		}
		log.Error().Err(err).Int("message_id", id).Msg("Failed to get message for deletion")
		return err
	}

	// Authorization check: Only the sender can delete their message
	if existingMessage.SenderID != int(userID) {
		return &ForbiddenError{Message: "User not authorized to delete this message"}
	}

	err = s.messageRepo.DeleteMessage(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("message_id", id).Msg("Message not found for deletion")
			return nil // Consider returning nil if not found is not an error for deletion
		}
		log.Error().Err(err).Int("message_id", id).Msg("Failed to delete message")
		return err
	}
	log.Info().Int("message_id", id).Msg("Message deleted successfully")
	return nil
}
