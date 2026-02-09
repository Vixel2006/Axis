package services

import (
	"context"

	"axis/internal/models"
	"axis/internal/repositories"
	"github.com/rs/zerolog/log"
)

type WorkspaceMemberService interface {
	AddMemberToWorkspace(ctx context.Context, workspaceID, userID int, role models.UserRole) (*models.WorkspaceMember, error)
	RemoveMemberFromWorkspace(ctx context.Context, workspaceID, userID int) error
	GetWorkspaceMembers(ctx context.Context, workspaceID int) ([]models.WorkspaceMember, error)
	JoinWorkspace(ctx context.Context, workspaceID, userID int) (*models.WorkspaceMember, error)
}

type workspaceMemberService struct {
	workspaceMemberRepo repositories.WorkspaceMemberRepo
	workspaceRepo       repositories.WorkspaceRepo // Added to check if workspace exists
}

func NewWorkspaceMemberService(wmr repositories.WorkspaceMemberRepo, wr repositories.WorkspaceRepo) WorkspaceMemberService {
	return &workspaceMemberService{
		workspaceMemberRepo: wmr,
		workspaceRepo:       wr,
	}
}

func (s *workspaceMemberService) AddMemberToWorkspace(ctx context.Context, workspaceID, userID int, role models.UserRole) (*models.WorkspaceMember, error) {
	err := s.workspaceMemberRepo.AddMemberToWorkspace(ctx, workspaceID, userID, role)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).Str("role", role.String()).Msg("Failed to add member to workspace")
		return nil, err
	}
	workspaceMember := &models.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
	}
	log.Info().Int("workspace_id", workspaceID).Int("user_id", userID).Str("role", role.String()).Msg("Member added to workspace successfully")
	return workspaceMember, nil
}

func (s *workspaceMemberService) RemoveMemberFromWorkspace(ctx context.Context, workspaceID, userID int) error {
	err := s.workspaceMemberRepo.RemoveMemberFromWorkspace(ctx, workspaceID, userID)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Failed to remove member from workspace")
		return err
	}
	log.Info().Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Member removed from workspace successfully")
	return nil
}

func (s *workspaceMemberService) GetWorkspaceMembers(ctx context.Context, workspaceID int) ([]models.WorkspaceMember, error) {
	members, err := s.workspaceMemberRepo.GetWorkspaceMembers(ctx, workspaceID)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to get workspace members")
		return nil, err
	}
	return members, nil
}

func (s *workspaceMemberService) JoinWorkspace(ctx context.Context, workspaceID, userID int) (*models.WorkspaceMember, error) {
	// 1. Check if the workspace exists
	workspace, err := s.workspaceRepo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to get workspace by ID")
		return nil, err
	}
	if workspace == nil {
		return nil, &NotFoundError{Message: "Workspace not found"}
	}

	// 2. Check if the user is already a member
	isMember, err := s.workspaceMemberRepo.IsMemberOfWorkspace(ctx, workspaceID, userID)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Failed to check if user is already a member")
		return nil, err
	}
	if isMember {
		return nil, &ConflictError{Message: "User is already a member of this workspace"}
	}

	// 3. Add the user as a member with the default 'Member' role
	err = s.workspaceMemberRepo.AddMemberToWorkspace(ctx, workspaceID, userID, models.Member)
	if err != nil {
		log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Failed to add user to workspace")
		return nil, err
	}

	workspaceMember := &models.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        models.Member,
	}
	log.Info().Int("workspace_id", workspaceID).Int("user_id", userID).Msg("User joined workspace successfully")
	return workspaceMember, nil
}
