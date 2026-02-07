package repositories

import (
	"context"
	"database/sql"

	"backend/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type WorkspaceRepo interface {
	CreateWorkspace(ctx context.Context, workspace *models.Workspace) error
	GetWorkspaceByID(ctx context.Context, workspaceID int) (*models.Workspace, error)
	UpdateWorkspace(ctx context.Context, workspace *models.Workspace) error
	DeleteWorkspace(ctx context.Context, workspaceID int) error
}

type workspaceRepository struct {
	db *bun.DB
}

func NewWorkspaceRepo(db *bun.DB) WorkspaceRepo {
	return &workspaceRepository{
		db: db,
	}
}

func (wr *workspaceRepository) CreateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	_, err := wr.db.NewInsert().Model(workspace).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Str("workspace_name", workspace.Name).Msg("Failed to create workspace")
		return err
	}
	return nil
}

func (wr *workspaceRepository) GetWorkspaceByID(ctx context.Context, workspaceID int) (*models.Workspace, error) {
	workspace := new(models.Workspace)
	err := wr.db.NewSelect().Model(workspace).Where("id = ?", workspaceID).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("workspace_id", workspaceID).Msg("Workspace not found")
			return nil, nil
		}
		log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to get workspace by ID")
		return nil, err
	}
	return workspace, nil
}

func (wr *workspaceRepository) UpdateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	_, err := wr.db.NewUpdate().Model(workspace).WherePK().Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspace.ID).Msg("Failed to update workspace")
		return err
	}
	return nil
}

func (wr *workspaceRepository) DeleteWorkspace(ctx context.Context, workspaceID int) error {
	_, err := wr.db.NewDelete().Model(&models.Workspace{}).Where("id = ?", workspaceID).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to delete workspace")
		return err
	}
	return nil
}
