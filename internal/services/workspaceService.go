package services

import (
	"context"
	"database/sql"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog"
)

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error)
	GetWorkspaceByID(ctx context.Context, id int) (*models.Workspace, error)
	GetWorkspaceByIDAuthorized(ctx context.Context, userID, workspaceID int) (*models.Workspace, error)
	UpdateWorkspace(ctx context.Context, userID int, workspace *models.Workspace) (*models.Workspace, error)
	DeleteWorkspace(ctx context.Context, userID int, id int) error
	GetWorkspacesForUser(ctx context.Context, userID int) ([]*models.Workspace, error)
}

type workspaceService struct {
	workspaceRepo       repositories.WorkspaceRepo
	workspaceMemberRepo repositories.WorkspaceMemberRepo
	log                 zerolog.Logger
}

func NewWorkspaceService(wr repositories.WorkspaceRepo, wmr repositories.WorkspaceMemberRepo, logger zerolog.Logger) WorkspaceService {
	return &workspaceService{
		workspaceRepo:       wr,
		workspaceMemberRepo: wmr,
		log:                 logger,
	}
}

func (s *workspaceService) CreateWorkspace(ctx context.Context, workspace *models.Workspace) (*models.Workspace, error) {
	err := s.workspaceRepo.CreateWorkspace(ctx, workspace)
	if err != nil {
		s.log.Error().Err(err).Str("workspace_name", workspace.Name).Msg("Failed to create workspace")
		return nil, err
	}
	s.log.Info().Str("workspace_name", workspace.Name).Int("workspace_id", workspace.ID).Msg("Workspace created successfully")

	err = s.workspaceMemberRepo.AddMemberToWorkspace(ctx, workspace.ID, workspace.CreatorID, models.Admin)
	if err != nil {
		s.log.Error().Err(err).Int("workspace_id", workspace.ID).Int("user_id", workspace.CreatorID).Msg("Failed to add creator as member to workspace")
		return nil, err
	}
	s.log.Info().Int("workspace_id", workspace.ID).Int("user_id", workspace.CreatorID).Msg("Creator added as admin to workspace")

	return workspace, nil
}

func (s *workspaceService) GetWorkspaceByID(ctx context.Context, id int) (*models.Workspace, error) {
	workspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("workspace_id", id).Msg("Workspace not found")
			return nil, nil
		}
		s.log.Error().Err(err).Int("workspace_id", id).Msg("Failed to get workspace by ID")
		return nil, err
	}
	return workspace, nil
}

func (s *workspaceService) GetWorkspaceByIDAuthorized(ctx context.Context, userID, workspaceID int) (*models.Workspace, error) {
	isMember, err := s.workspaceMemberRepo.IsMemberOfWorkspace(ctx, workspaceID, userID)
	if err != nil {
		s.log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Failed to check workspace membership")
		return nil, err
	}
	if !isMember {
		s.log.Warn().Int("workspace_id", workspaceID).Int("user_id", userID).Msg("User not authorized to access this workspace")
		return nil, &ForbiddenError{Message: "User not authorized to access this workspace"}
	}

	workspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("workspace_id", workspaceID).Msg("Workspace not found")
			return nil, nil
		}
		s.log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to get workspace by ID after authorization")
		return nil, err
	}
	return workspace, nil
}

func (s *workspaceService) UpdateWorkspace(ctx context.Context, userID int, workspace *models.Workspace) (*models.Workspace, error) {
	existingWorkspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, workspace.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("workspace_id", workspace.ID).Msg("Workspace not found for update")
			return nil, nil
		}
		s.log.Error().Err(err).Int("workspace_id", workspace.ID).Msg("Failed to get workspace for update")
		return nil, err
	}

	if existingWorkspace.CreatorID != int(userID) {
		isMember, err := s.workspaceMemberRepo.IsMemberOfWorkspace(ctx, existingWorkspace.ID, int(userID))
		if err != nil {
			s.log.Error().Err(err).Int("workspace_id", existingWorkspace.ID).Int("user_id", int(userID)).Msg("Failed to check workspace membership for update")
			return nil, err
		}
		if !isMember {
			s.log.Warn().Int("workspace_id", existingWorkspace.ID).Int("user_id", int(userID)).Msg("User not authorized to update this workspace")
			return nil, &ForbiddenError{Message: "User not authorized to update this workspace"}
		}
		member, err := s.workspaceMemberRepo.GetWorkspaceMember(ctx, existingWorkspace.ID, int(userID))
		if err != nil && err != sql.ErrNoRows {
			s.log.Error().Err(err).Int("workspace_id", existingWorkspace.ID).Int("user_id", int(userID)).Msg("Failed to get workspace member role for update")
			return nil, err
		}
		if member == nil || member.Role != models.Admin {
			s.log.Warn().Int("workspace_id", existingWorkspace.ID).Int("user_id", int(userID)).Msg("User does not have admin role to update this workspace")
			return nil, &ForbiddenError{Message: "User not authorized to update this workspace"}
		}
	}

	existingWorkspace.Name = workspace.Name 

	err = s.workspaceRepo.UpdateWorkspace(ctx, existingWorkspace)
	if err != nil {
		s.log.Error().Err(err).Int("workspace_id", existingWorkspace.ID).Msg("Failed to update workspace")
		return nil, err
	}
	s.log.Info().Int("workspace_id", existingWorkspace.ID).Msg("Workspace updated successfully")
	return existingWorkspace, nil
}

func (s *workspaceService) DeleteWorkspace(ctx context.Context, userID int, id int) error {
	existingWorkspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("workspace_id", id).Msg("Workspace not found for deletion")
			return nil
		}
		s.log.Error().Err(err).Int("workspace_id", id).Msg("Failed to get workspace for deletion")
		return err
	}

	if existingWorkspace.CreatorID != int(userID) {
		isMember, err := s.workspaceMemberRepo.IsMemberOfWorkspace(ctx, existingWorkspace.ID, int(userID))
		if err != nil {
			s.log.Error().Err(err).Int("workspace_id", existingWorkspace.ID).Int("user_id", int(userID)).Msg("Failed to check workspace membership for deletion")
			return err
		}
		if !isMember {
			s.log.Warn().Int("workspace_id", existingWorkspace.ID).Int("user_id", int(userID)).Msg("User not authorized to delete this workspace")
			return &ForbiddenError{Message: "User not authorized to delete this workspace"}
		}
		member, err := s.workspaceMemberRepo.GetWorkspaceMember(ctx, existingWorkspace.ID, int(userID))
		if err != nil && err != sql.ErrNoRows {
			s.log.Error().Err(err).Int("workspace_id", existingWorkspace.ID).Int("user_id", int(userID)).Msg("Failed to get workspace member role for deletion")
			return err
		}
		if member == nil || member.Role != models.Admin {
			s.log.Warn().Int("workspace_id", existingWorkspace.ID).Int("user_id", int(userID)).Msg("User does not have admin role to delete this workspace")
			return &ForbiddenError{Message: "User not authorized to delete this workspace"}
		}
	}

	err = s.workspaceRepo.DeleteWorkspace(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Info().Int("workspace_id", id).Msg("Workspace not found for deletion")
			return nil
		}
		s.log.Error().Err(err).Int("workspace_id", id).Msg("Failed to delete workspace")
		return err
	}
	s.log.Info().Int("workspace_id", id).Msg("Workspace deleted successfully")
	return nil
}

func (s *workspaceService) GetWorkspacesForUser(ctx context.Context, userID int) ([]*models.Workspace, error) {
	s.log.Debug().Int("user_id", int(userID)).Msg("Calling WorkspaceMemberRepo.GetWorkspacesForUser")
	memberships, err := s.workspaceMemberRepo.GetWorkspacesForUser(ctx, int(userID))
	if err != nil {
		s.log.Error().Err(err).Int("user_id", int(userID)).Msg("Failed to get workspace memberships for user")
		return nil, err
	}
	s.log.Debug().Int("user_id", int(userID)).Int("memberships_count", len(memberships)).Msg("Received workspace memberships from repo")

	workspaces := make([]*models.Workspace, 0, len(memberships))
	for i := range memberships {
		if memberships[i].Workspace != nil {
			workspaces = append(workspaces, memberships[i].Workspace)
		} else {
			s.log.Warn().Int("user_id", int(userID)).Int("membership_index", i).Msg("Workspace relation is nil for a membership")
		}
	}

	s.log.Info().Int("user_id", int(userID)).Int("workspace_count", len(workspaces)).Msg("Retrieved workspaces for user successfully")
	return workspaces, nil
}
