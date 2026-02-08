package services

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog/log"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	GetMessageByID(ctx context.Context, id int) (*models.Message, error)
	GetMessagesInChannel(ctx context.Context, channelID int, limit, offset int) ([]models.Message, error)
	UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	DeleteMessage(ctx context.Context, id int) error
}

type messageService struct {
	messageRepo repositories.MessageRepo
}

func NewMessageService(mr repositories.MessageRepo) MessageService {
	return &messageService{
		messageRepo: mr,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	err := s.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		log.Error().Err(err).Int("channel_id", message.ChannelID).Int("sender_id", message.SenderID).Msg("Failed to create message")
		return nil, err
	}
	log.Info().Int("message_id", message.ID).Int("channel_id", message.ChannelID).Msg("Message created successfully")
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

func (s *messageService) GetMessagesInChannel(ctx context.Context, channelID int, limit, offset int) ([]models.Message, error) {
	messages, err := s.messageRepo.GetMessagesByChannelID(ctx, channelID, limit, offset)
	if err != nil {
		log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to get messages for channel")
		return nil, err
	}
	return messages, nil
}

func (s *messageService) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	existingMessage, err := s.messageRepo.GetMessageByID(ctx, message.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("message_id", message.ID).Msg("Message not found for update")
			return nil, nil
		}
		log.Error().Err(err).Int("message_id", message.ID).Msg("Failed to get message for update")
		return nil, err
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

func (s *messageService) DeleteMessage(ctx context.Context, id int) error {
	err := s.messageRepo.DeleteMessage(ctx, id)
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
