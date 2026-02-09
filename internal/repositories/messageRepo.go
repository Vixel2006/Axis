package repositories

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type MessageRepo interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessageByID(ctx context.Context, messageID int) (*models.Message, error)
	GetMessagesByMeetingID(ctx context.Context, meetingID int, limit, offset int) ([]models.Message, error)
	GetThreadedMessages(ctx context.Context, parentMessageID int) ([]models.Message, error)
	UpdateMessage(ctx context.Context, message *models.Message) error
	DeleteMessage(ctx context.Context, messageID int) error
}

type messageRepository struct {
	db *bun.DB
}

func NewMessageRepo(db *bun.DB) MessageRepo {
	return &messageRepository{
		db: db,
	}
}

func (mr *messageRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	_, err := mr.db.NewInsert().Model(message).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("meeting_id", message.MeetingID).Int("sender_id", message.SenderID).Msg("Failed to create message")
		return err
	}
	return nil
}

func (mr *messageRepository) GetMessageByID(ctx context.Context, messageID int) (*models.Message, error) {
	message := new(models.Message)
	err := mr.db.NewSelect().Model(message).Where("id = ?", messageID).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("message_id", messageID).Msg("Message not found")
			return nil, nil
		}
		log.Error().Err(err).Int("message_id", messageID).Msg("Failed to get message by ID")
		return nil, err
	}
	return message, nil
}

func (mr *messageRepository) GetMessagesByMeetingID(ctx context.Context, meetingID int, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := mr.db.NewSelect().
		Model(&messages).
		Where("meeting_id = ?", meetingID).
		Order("created_at DESC"). // Latest messages first
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Int("meeting_id", meetingID).Msg("Failed to get messages by meeting ID")
		return nil, err
	}
	return messages, nil
}

func (mr *messageRepository) GetThreadedMessages(ctx context.Context, parentMessageID int) ([]models.Message, error) {
	var messages []models.Message
	err := mr.db.NewSelect().
		Model(&messages).
		Where("parent_message_id = ?", parentMessageID).
		Order("created_at ASC"). // Order by creation for thread flow
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Int("parent_message_id", parentMessageID).Msg("Failed to get threaded messages")
		return nil, err
	}
	return messages, nil
}

func (mr *messageRepository) UpdateMessage(ctx context.Context, message *models.Message) error {
	_, err := mr.db.NewUpdate().Model(message).WherePK().Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("message_id", message.ID).Msg("Failed to update message")
		return err
	}
	return nil
}

func (mr *messageRepository) DeleteMessage(ctx context.Context, messageID int) error {
	_, err := mr.db.NewDelete().Model(&models.Message{}).Where("id = ?", messageID).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("message_id", messageID).Msg("Failed to delete message")
		return err
	}
	return nil
}
