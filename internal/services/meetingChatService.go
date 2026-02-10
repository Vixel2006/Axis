package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"sync"
	"time"

	"axis/internal/models"
	"axis/internal/repositories"
	"axis/internal/utils"
)

type MeetingChatService interface {
	JoinMeetingChat(ctx context.Context, meetingID, userID int) error
	LeaveMeetingChat(ctx context.Context, meetingID, userID int) error
	SendMessage(ctx context.Context, meetingID, senderID int, parentMessageID *int, content string, messageType models.MessageType, attachments []models.SendAttachmentDetails) (*models.Message, error)
	AddReaction(ctx context.Context, messageID, userID int, emoji string) (*models.ReactionBroadcastPayload, error)
	RemoveReaction(ctx context.Context, messageID, userID int, emoji string) (*models.ReactionBroadcastPayload, error)
	GetMeetingMessages(ctx context.Context, meetingID int, limit, offset int) ([]models.Message, error)
	GetOrCreateHubForMeeting(meetingID int) *utils.Hub
	BroadcastMessage(meetingID int, message []byte)
	RegisterClient(meetingID int, client *utils.Client)
	UnregisterClient(meetingID int, client *utils.Client)
}

type meetingChatService struct {
	meetingRepo    repositories.MeetingRepo
	messageRepo    repositories.MessageRepo
	userRepo       repositories.UserRepo
	attachmentRepo repositories.AttachmentRepo
	reactionRepo   repositories.ReactionRepo
	log            zerolog.Logger
	hubs           map[int]*utils.Hub
	mu             sync.Mutex
}

func NewMeetingChatService(mr repositories.MeetingRepo, msgRepo repositories.MessageRepo, ur repositories.UserRepo, ar repositories.AttachmentRepo, rr repositories.ReactionRepo, logger zerolog.Logger) MeetingChatService {
	return &meetingChatService{
		meetingRepo:    mr,
		messageRepo:    msgRepo,
		userRepo:       ur,
		attachmentRepo: ar,
		reactionRepo:   rr,
		log:            logger,
		hubs:           make(map[int]*utils.Hub),
	}
}

func (s *meetingChatService) JoinMeetingChat(ctx context.Context, meetingID, userID int) error {
	meeting, err := s.meetingRepo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to get meeting by ID for joining chat")
		if err == sql.ErrNoRows {
			return NewNotFoundError(fmt.Sprintf("Meeting with ID %d not found", meetingID))
		}
		return fmt.Errorf("database error: %w", err)
	}
	if meeting == nil {
		s.log.Warn().Int("meeting_id", meetingID).Msg("Meeting not found when trying to join chat")
		return NewNotFoundError(fmt.Sprintf("Meeting with ID %d not found", meetingID))
	}

	isParticipant, err := s.meetingRepo.IsParticipantInMeeting(ctx, meetingID, userID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", userID).Msg("Failed to check if user is participant")
		return fmt.Errorf("failed to check participant status: %w", err)
	}
	if isParticipant {
		s.log.Info().Int("meeting_id", meetingID).Int("user_id", userID).Msg("User is already a participant in meeting chat")
		return nil
	}

	err = s.meetingRepo.AddParticipantToMeeting(ctx, meetingID, userID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", userID).Msg("Failed to add participant to meeting chat")
		return fmt.Errorf("failed to add participant: %w", err)
	}

	s.log.Info().Int("meeting_id", meetingID).Int("user_id", userID).Msg("User joined meeting chat successfully")
	return nil
}

func (s *meetingChatService) LeaveMeetingChat(ctx context.Context, meetingID, userID int) error {
	meeting, err := s.meetingRepo.GetMeetingByID(ctx, meetingID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to get meeting by ID for leaving chat")
		if err == sql.ErrNoRows {
			return NewNotFoundError(fmt.Sprintf("Meeting with ID %d not found", meetingID))
		}
		return fmt.Errorf("database error: %w", err)
	}
	if meeting == nil {
		s.log.Warn().Int("meeting_id", meetingID).Msg("Meeting not found when trying to leave chat")
		return NewNotFoundError(fmt.Sprintf("Meeting with ID %d not found", meetingID))
	}

	err = s.meetingRepo.RemoveParticipantFromMeeting(ctx, meetingID, userID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Int("user_id", userID).Msg("Failed to remove participant from meeting chat")
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	s.log.Info().Int("meeting_id", meetingID).Int("user_id", userID).Msg("User left meeting chat successfully")
	return nil
}

func (s *meetingChatService) SendMessage(ctx context.Context, meetingID, senderID int, parentMessageID *int, content string, messageType models.MessageType, attachments []models.SendAttachmentDetails) (*models.Message, error) {
	isParticipant, err := s.meetingRepo.IsParticipantInMeeting(ctx, meetingID, senderID)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Int("sender_id", senderID).Msg("Failed to check if sender is participant before sending message")
		return nil, fmt.Errorf("failed to verify participant status: %w", err)
	}
	if !isParticipant {
		s.log.Warn().Int("meeting_id", meetingID).Int("sender_id", senderID).Msg("Sender is not a participant in the meeting chat")
		return nil, NewUnauthorizedError("sender is not a participant in this chat")
	}

	message := &models.Message{
		MeetingID:       meetingID,
		SenderID:        senderID,
		Content:         content,
		MessageType:     messageType,
		ParentMessageID: parentMessageID,
		CreatedAt:       time.Now(),
	}

	err = s.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Int("sender_id", senderID).Msg("Failed to create message in database")
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	for _, attachDetail := range attachments {
		attachment := &models.Attachment{
			MessageID: message.ID,
			UserID:    senderID,
			FileName:  attachDetail.FileName,
			FileType:  attachDetail.FileType,
			FileSize:  attachDetail.FileSize,
			URL:       attachDetail.URL,
			CreatedAt: time.Now(),
		}
		err := s.attachmentRepo.CreateAttachment(ctx, attachment)
		if err != nil {
			s.log.Error().Err(err).Int("message_id", message.ID).Str("file_name", attachDetail.FileName).Msg("Failed to create attachment in database")
			continue
		}
	}

	if message.ID != 0 {
		if len(attachments) > 0 {
			fetchedAttachments, err := s.attachmentRepo.GetAttachmentsByMessageID(ctx, message.ID)
			if err != nil {
				s.log.Error().Err(err).Int("message_id", message.ID).Msg("Failed to fetch attachments for message after creation")
			} else {
				message.Attachments = make([]*models.Attachment, len(fetchedAttachments))
				for i := range fetchedAttachments {
					message.Attachments[i] = &fetchedAttachments[i]
				}
			}
		}
	}

	s.log.Info().Int("message_id", message.ID).Int("meeting_id", meetingID).Int("sender_id", senderID).Msg("Message sent and saved successfully")
	return message, nil
}

func (s *meetingChatService) GetMeetingMessages(ctx context.Context, meetingID int, limit, offset int) ([]models.Message, error) {
	messages, err := s.messageRepo.GetMessagesByMeetingID(ctx, meetingID, limit, offset)
	if err != nil {
		s.log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to retrieve messages for meeting")
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	s.log.Debug().Int("meeting_id", meetingID).Int("count", len(messages)).Msg("Retrieved meeting messages")
	return messages, nil
}

func (s *meetingChatService) AddReaction(ctx context.Context, messageID, userID int, emoji string) (*models.ReactionBroadcastPayload, error) {
	_, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		s.log.Error().Err(err).Int("message_id", messageID).Msg("Message not found for reaction")
		if err == sql.ErrNoRows {
			return nil, NewNotFoundError(fmt.Sprintf("Message with ID %d not found", messageID))
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	reaction := &models.Reaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}

	err = s.reactionRepo.CreateReaction(ctx, reaction)
	if err != nil {
		s.log.Error().Err(err).Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Failed to add reaction to message")
		return nil, fmt.Errorf("failed to add reaction: %w", err)
	}

	s.log.Info().Int("reaction_id", reaction.ID).Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Reaction added successfully")

	return &models.ReactionBroadcastPayload{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
		Action:    "added",
	}, nil
}

func (s *meetingChatService) RemoveReaction(ctx context.Context, messageID, userID int, emoji string) (*models.ReactionBroadcastPayload, error) {
	// Check if message exists (optional, could just try to delete)
	_, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		s.log.Error().Err(err).Int("message_id", messageID).Msg("Message not found for removing reaction")
		if err == sql.ErrNoRows {
			// If message not found, we can consider it successfully "removed" if that's the desired behavior,
			// or return an error if we strictly need the message to exist. For now, return nil,nil payload.
			return &models.ReactionBroadcastPayload{
				MessageID: messageID,
				UserID:    userID,
				Emoji:     emoji,
				Action:    "removed",
			}, nil
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	err = s.reactionRepo.DeleteReaction(ctx, messageID, userID, emoji)
	if err != nil {
		s.log.Error().Err(err).Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Failed to remove reaction from message")
		return nil, fmt.Errorf("failed to remove reaction: %w", err)
	}

	s.log.Info().Int("message_id", messageID).Int("user_id", userID).Str("emoji", emoji).Msg("Reaction removed successfully")
	return &models.ReactionBroadcastPayload{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
		Action:    "removed",
	}, nil
}

func (s *meetingChatService) GetOrCreateHubForMeeting(meetingID int) *utils.Hub {
	s.mu.Lock()
	defer s.mu.Unlock()

	if hub, ok := s.hubs[meetingID]; ok {
		return hub
	}

	hub := utils.NewHub(s.log)
	s.hubs[meetingID] = hub
	go hub.Run()
	s.log.Info().Int("meeting_id", meetingID).Msg("Created and started new WebSocket hub for meeting")
	return hub
}

func (s *meetingChatService) BroadcastMessage(meetingID int, message []byte) {
	s.mu.Lock()
	hub, ok := s.hubs[meetingID]
	s.mu.Unlock()

	if !ok {
		s.log.Warn().Int("meeting_id", meetingID).Msg("Attempted to broadcast to a non-existent hub")
		return
	}
	s.log.Debug().Int("meeting_id", meetingID).Msg("Broadcasting message to meeting hub")
	hub.Broadcast <- message
}

func (s *meetingChatService) RegisterClient(meetingID int, client *utils.Client) {
	s.mu.Lock()
	hub, ok := s.hubs[meetingID]
	s.mu.Unlock()

	if !ok {
		s.log.Warn().Int("meeting_id", meetingID).Msg("Attempted to register client to a non-existent hub")
		return
	}
	s.log.Debug().Int("meeting_id", meetingID).Int("user_id", client.ID).Msg("Registering client to meeting hub")
	hub.Register <- client
}

func (s *meetingChatService) UnregisterClient(meetingID int, client *utils.Client) {
	s.mu.Lock()
	hub, ok := s.hubs[meetingID]
	s.mu.Unlock()

	if !ok {
		s.log.Warn().Int("meeting_id", meetingID).Msg("Attempted to unregister client from a non-existent hub")
		return
	}
	s.log.Debug().Int("meeting_id", meetingID).Int("user_id", client.ID).Msg("Unregistering client from meeting hub")
	hub.UnRegister <- client
}
