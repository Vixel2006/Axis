package services

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog"
)

type AttachmentService interface {
	CreateAttachment(ctx context.Context, attachment *models.Attachment) (*models.Attachment, error)
	GetAttachmentByID(ctx context.Context, id int) (*models.Attachment, error)
	GetAttachmentsForMessage(ctx context.Context, messageID int) ([]models.Attachment, error)
}

type attachmentService struct {
	attachmentRepo repositories.AttachmentRepo
	log            zerolog.Logger
}

func NewAttachmentService(ar repositories.AttachmentRepo, logger zerolog.Logger) AttachmentService {
	return &attachmentService{
		attachmentRepo: ar,
		log:            logger,
	}
}

func (s *attachmentService) CreateAttachment(ctx context.Context, attachment *models.Attachment) (*models.Attachment, error) {
	err := s.attachmentRepo.CreateAttachment(ctx, attachment)
	if err != nil {
		s.log.Error().Err(err).Str("filename", attachment.FileName).Msg("Failed to create attachment")
		return nil, err
	}
	s.log.Info().Int("attachment_id", attachment.ID).Str("filename", attachment.FileName).Msg("Attachment created successfully")
	return attachment, nil
}

func (s *attachmentService) GetAttachmentByID(ctx context.Context, id int) (*models.Attachment, error) {
	attachment, err := s.attachmentRepo.GetAttachmentByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("attachment_id", id).Msg("Attachment not found")
			return nil, nil
		}
		s.log.Error().Err(err).Int("attachment_id", id).Msg("Failed to get attachment by ID")
		return nil, err
	}
	return attachment, nil
}

func (s *attachmentService) GetAttachmentsForMessage(ctx context.Context, messageID int) ([]models.Attachment, error) {
	attachments, err := s.attachmentRepo.GetAttachmentsByMessageID(ctx, messageID)
	if err != nil {
		s.log.Error().Err(err).Int("message_id", messageID).Msg("Failed to get attachments for message")
		return nil, err
	}
	return attachments, nil
}
