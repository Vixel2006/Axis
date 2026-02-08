package repositories

import (
	"context"

	"axis/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type WorkspaceMemberRepo interface {
	AddMemberToWorkspace(ctx context.Context, workspaceID, userID int, role models.UserRole) error
	RemoveMemberFromWorkspace(ctx context.Context, workspaceID, userID int) error
	GetWorkspaceMembers(ctx context.Context, workspaceID int) ([]models.WorkspaceMember, error)
	GetWorkspacesForUser(ctx context.Context, userID int) ([]models.WorkspaceMember, error)
	UpdateWorkspaceMemberRole(ctx context.Context, workspaceID, userID int, role models.UserRole) error
}

type workspaceMemberRepository struct {
	db *bun.DB
}

func NewWorkspaceMemberRepo(db *bun.DB) WorkspaceMemberRepo {
	return &workspaceMemberRepository{
		db: db,
	}
}

func (wmr *workspaceMemberRepository) AddMemberToWorkspace(ctx context.Context, workspaceID, userID int, role models.UserRole) error {
	workspaceMember := &models.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
	}
	_, err := wmr.db.NewInsert().Model(workspaceMember).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).
			Str("role", role.String()).Msg("Failed to add member to workspace")
		return err
	}
	return nil
}

func (wmr *workspaceMemberRepository) RemoveMemberFromWorkspace(ctx context.Context, workspaceID, userID int) error {
	_, err := wmr.db.NewDelete().
		Model(&models.WorkspaceMember{}).
		Where("workspace_id = ?", workspaceID).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Failed to remove member from workspace")
		return err
	}
	return nil
}

func (wmr *workspaceMemberRepository) GetWorkspaceMembers(ctx context.Context, workspaceID int) ([]models.WorkspaceMember, error) {
	var members []models.WorkspaceMember
	err := wmr.db.NewSelect().
		Model(&members).
		Where("workspace_id = ?", workspaceID).
		Relation("User"). // Eager load user details
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to get workspace members")
		return nil, err
	}
	return members, nil
}

func (wmr *workspaceMemberRepository) GetWorkspacesForUser(ctx context.Context, userID int) ([]models.WorkspaceMember, error) {
	log.Debug().Int("user_id", userID).Msg("Fetching workspaces for user from database")
	var memberships []models.WorkspaceMember
	err := wmr.db.NewSelect().
		Model(&memberships).
		Where("user_id = ?", userID).
		Relation("Workspace"). // Eager load workspace details
		Scan(ctx)
	if err != nil {
		log.Error().Err(err).Int("user_id", userID).Msg("Failed to get workspaces for user")
		return nil, err
	}
	log.Debug().Int("user_id", userID).Int("memberships_count", len(memberships)).Msg("Successfully fetched workspace memberships")
	return memberships, nil
}

func (wmr *workspaceMemberRepository) UpdateWorkspaceMemberRole(ctx context.Context, workspaceID, userID int, role models.UserRole) error {
	_, err := wmr.db.NewUpdate().
		Model(&models.WorkspaceMember{Role: role}).
		Where("workspace_id = ?", workspaceID).
		Where("user_id = ?", userID).
		Set("role = ?", role).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).
			Str("role", role.String()).Msg("Failed to update workspace member role")
		return err
	}
	return nil
}
