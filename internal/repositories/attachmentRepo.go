package repositories

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type AttachmentRepo interface {
	CreateAttachment(ctx context.Context, attachment *models.Attachment) error
	GetAttachmentByID(ctx context.Context, attachmentID int) (*models.Attachment, error)
	GetAttachmentsByMessageID(ctx context.Context, messageID int) ([]models.Attachment, error)
	DeleteAttachment(ctx context.Context, attachmentID int) error
}

type attachmentRepository struct {
	db  *bun.DB
	log zerolog.Logger
}

func NewAttachmentRepo(db *bun.DB, logger zerolog.Logger) AttachmentRepo {
	return &attachmentRepository{
		db:  db,
		log: logger,
	}
}

func (ar *attachmentRepository) CreateAttachment(ctx context.Context, attachment *models.Attachment) error {
	_, err := ar.db.NewInsert().Model(attachment).Exec(ctx)
	if err != nil {
		ar.log.Error().Err(err).Msg("Failed to create attachment")
		return err
	}
	return nil
}

func (ar *attachmentRepository) GetAttachmentByID(ctx context.Context, attachmentID int) (*models.Attachment, error) {
	attachment := new(models.Attachment)
	err := ar.db.NewSelect().Model(attachment).Where("id = ?", attachmentID).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ar.log.Info().Int("attachment_id", attachmentID).Msg("Attachment not found")
			return nil, nil // Return nil, nil if no attachment found
		}
		ar.log.Error().Err(err).Int("attachment_id", attachmentID).Msg("Failed to get attachment by ID")
		return nil, err
	}
	return attachment, nil
}

func (ar *attachmentRepository) GetAttachmentsByMessageID(ctx context.Context, messageID int) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := ar.db.NewSelect().Model(&attachments).Where("message_id = ?", messageID).Scan(ctx)
	if err != nil {
		ar.log.Error().Err(err).Int("message_id", messageID).Msg("Failed to get attachments by message ID")
		return nil, err
	}
	return attachments, nil
}

func (ar *attachmentRepository) DeleteAttachment(ctx context.Context, attachmentID int) error {
	_, err := ar.db.NewDelete().Model(&models.Attachment{}).Where("id = ?", attachmentID).Exec(ctx)
	if err != nil {
		ar.log.Error().Err(err).Int("attachment_id", attachmentID).Msg("Failed to delete attachment")
		return err
	}
	return nil
}
