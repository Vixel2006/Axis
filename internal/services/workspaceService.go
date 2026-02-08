package services

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog/log"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error)
	GetWorkspaceByID(ctx context.Context, id int) (*models.Workspace, error)
	UpdateWorkspace(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error)
	DeleteWorkspace(ctx context.Context, id int) error
	GetWorkspacesForUser(ctx context.Context, userID int) ([]*models.Workspace, error)
}

type workspaceService struct {
	workspaceRepo       repositories.WorkspaceRepo
	workspaceMemberRepo repositories.WorkspaceMemberRepo
}

func NewWorkspaceService(wr repositories.WorkspaceRepo, wmr repositories.WorkspaceMemberRepo) WorkspaceService {
	return &workspaceService{
		workspaceRepo:       wr,
		workspaceMemberRepo: wmr,
	}
}

func (s *workspaceService) CreateWorkspace(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error) {
	err := s.workspaceRepo.CreateWorkspace(ctx, workspace)
	if err != nil {
		log.Error().Err(err).Str("workspace_name", workspace.Name).Msg("Failed to create workspace")
		return nil, err
	}
	log.Info().Str("workspace_name", workspace.Name).Int("workspace_id", workspace.ID).Msg("Workspace created successfully")

	// Add the creator as a member of the new workspace
	err = s.workspaceMemberRepo.AddMemberToWorkspace(ctx, workspace.ID, workspace.CreatorID, models.Admin)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspace.ID).Int("user_id", workspace.CreatorID).Msg("Failed to add creator as member to workspace")
		// Decide if this error should prevent workspace creation or just be logged
		// For now, we'll return the error as it's a critical step
		return nil, err
	}
	log.Info().Int("workspace_id", workspace.ID).Int("user_id", workspace.CreatorID).Msg("Creator added as admin to workspace")

	return workspace, nil
}

func (s *workspaceService) GetWorkspaceByID(ctx context.Context, id int) (*models.Workspace, error) {
	workspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("workspace_id", id).Msg("Workspace not found")
			return nil, nil // Return nil, nil for not found case
		}
		log.Error().Err(err).Int("workspace_id", id).Msg("Failed to get workspace by ID")
		return nil, err
	}
	return workspace, nil
}

func (s *workspaceService) UpdateWorkspace(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error) {
	existingWorkspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, workspace.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("workspace_id", workspace.ID).Msg("Workspace not found for update")
			return nil, nil
		}
		log.Error().Err(err).Int("workspace_id", workspace.ID).Msg("Failed to get workspace for update")
		return nil, err
	}

	// Update fields
	existingWorkspace.Name = workspace.Name // Assuming Name is always provided for update

	err = s.workspaceRepo.UpdateWorkspace(ctx, existingWorkspace)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", existingWorkspace.ID).Msg("Failed to update workspace")
		return nil, err
	}
	log.Info().Int("workspace_id", existingWorkspace.ID).Msg("Workspace updated successfully")
	return existingWorkspace, nil
}

func (s *workspaceService) DeleteWorkspace(ctx context.Context, id int) error {
	err := s.workspaceRepo.DeleteWorkspace(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int("workspace_id", id).Msg("Workspace not found for deletion")
			return nil // Consider returning nil if not found is not an error for deletion
		}
		log.Error().Err(err).Int("workspace_id", id).Msg("Failed to delete workspace")
		return err
	}
	log.Info().Int("workspace_id", id).Msg("Workspace deleted successfully")
	return nil
}

func (s *workspaceService) GetWorkspacesForUser(ctx context.Context, userID int) ([]*models.Workspace, error) {
	log.Debug().Int("user_id", userID).Msg("Calling WorkspaceMemberRepo.GetWorkspacesForUser")
	memberships, err := s.workspaceMemberRepo.GetWorkspacesForUser(ctx, userID)
	if err != nil {
		log.Error().Err(err).Int("user_id", userID).Msg("Failed to get workspace memberships for user")
		return nil, err
	}
	log.Debug().Int("user_id", userID).Int("memberships_count", len(memberships)).Msg("Received workspace memberships from repo")

	workspaces := make([]*models.Workspace, 0, len(memberships))
	for i := range memberships {
		if memberships[i].Workspace != nil {
			workspaces = append(workspaces, memberships[i].Workspace)
		} else {
			log.Warn().Int("user_id", userID).Int("membership_index", i).Msg("Workspace relation is nil for a membership")
		}
	}

	log.Info().Int("user_id", userID).Int("workspace_count", len(workspaces)).Msg("Retrieved workspaces for user successfully")
	return workspaces, nil
}
