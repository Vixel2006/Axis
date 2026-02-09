package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils" // Import the utils package
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type WorkspaceMemberHandler struct {
	workspaceMemberService services.WorkspaceMemberService
	log                    zerolog.Logger
}

func NewWorkspaceMemberHandler(wms services.WorkspaceMemberService, logger zerolog.Logger) *WorkspaceMemberHandler {
	return &WorkspaceMemberHandler{
		workspaceMemberService: wms,
		log:                    logger,
	}
}

func (h *WorkspaceMemberHandler) AddMemberToWorkspace(c *gin.Context) {
	h.log.Info().Msg("Handling AddMemberToWorkspace request")
	workspaceIDStr := c.Param("workspaceID")
	h.log.Debug().Str("workspaceID_param", workspaceIDStr).Msg("Parsing workspace ID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("workspaceID_param", workspaceIDStr).Msg("Invalid workspace ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var reqBody struct {
		UserID int             `json:"user_id"`
		Role   models.UserRole `json:"role"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for AddMemberToWorkspace")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Int("workspace_id", workspaceID).Int("user_id", reqBody.UserID).Str("role", reqBody.Role.String()).Msg("AddMemberToWorkspace request body")

	workspaceMember, err := h.workspaceMemberService.AddMemberToWorkspace(c.Request.Context(), workspaceID, reqBody.UserID, reqBody.Role)
	if err != nil {
		h.log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", reqBody.UserID).Msg("Failed to add member to workspace via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member to workspace"})
		return
	}

	h.log.Info().Int("workspace_id", workspaceID).Int("user_id", reqBody.UserID).Msg("Member added to workspace successfully")
	c.JSON(http.StatusCreated, workspaceMember)
}

func (h *WorkspaceMemberHandler) RemoveMemberFromWorkspace(c *gin.Context) {
	h.log.Info().Msg("Handling RemoveMemberFromWorkspace request")
	workspaceIDStr := c.Param("workspaceID")
	h.log.Debug().Str("workspaceID_param", workspaceIDStr).Msg("Parsing workspace ID for removal")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("workspaceID_param", workspaceIDStr).Msg("Invalid workspace ID format for removal")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	userIDStr := c.Param("userID")
	h.log.Debug().Str("userID_param", userIDStr).Msg("Parsing user ID for removal")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("userID_param", userIDStr).Msg("Invalid user ID format for removal")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	h.log.Debug().Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Attempting to remove member from workspace")

	err = h.workspaceMemberService.RemoveMemberFromWorkspace(c.Request.Context(), workspaceID, userID)
	if err != nil {
		h.log.Error().Err(err).Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Failed to remove member from workspace via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member from workspace"})
		return
	}

	h.log.Info().Int("workspace_id", workspaceID).Int("user_id", userID).Msg("Member removed from workspace successfully")
	c.JSON(http.StatusNoContent, nil)
}

func (h *WorkspaceMemberHandler) GetWorkspaceMembers(c *gin.Context) {
	h.log.Info().Msg("Handling GetWorkspaceMembers request")
	workspaceIDStr := c.Param("workspaceID")
	h.log.Debug().Str("workspaceID_param", workspaceIDStr).Msg("Parsing workspace ID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("workspaceID_param", workspaceIDStr).Msg("Invalid workspace ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}
	h.log.Debug().Int("workspace_id", workspaceID).Msg("Retrieving workspace members")

	members, err := h.workspaceMemberService.GetWorkspaceMembers(c.Request.Context(), workspaceID)
	if err != nil {
		h.log.Error().Err(err).Int("workspace_id", workspaceID).Msg("Failed to retrieve workspace members via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workspace members"})
		return
	}

	h.log.Info().Int("workspace_id", workspaceID).Int("members_count", len(members)).Msg("Workspace members retrieved successfully")
	c.JSON(http.StatusOK, members)
}

func (h *WorkspaceMemberHandler) JoinWorkspace(c *gin.Context) {
	h.log.Info().Msg("Handling JoinWorkspace request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in JoinWorkspace")
		return // GetUserIDFromContext already handles the error response
	}

	workspaceIDStr := c.Param("workspaceID")
	h.log.Debug().Str("workspaceID_param", workspaceIDStr).Msg("Parsing workspace ID for joining")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("workspaceID_param", workspaceIDStr).Msg("Invalid workspace ID format for joining")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}
	h.log.Debug().Int("user_id", int(userID)).Int("workspace_id", workspaceID).Msg("Attempting to join workspace")

	workspaceMember, err := h.workspaceMemberService.JoinWorkspace(c.Request.Context(), workspaceID, int(userID))
	if err != nil {
		// Handle specific errors from service layer
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("workspace_id", workspaceID).Msg("Workspace not found for joining")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.ConflictError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("workspace_id", workspaceID).Msg("User already a member of workspace")
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("workspace_id", workspaceID).Msg("Failed to join workspace via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join workspace"})
		return
	}

	h.log.Info().Int("workspace_id", workspaceID).Int("user_id", int(userID)).Msg("User joined workspace successfully")
	c.JSON(http.StatusCreated, workspaceMember)
}
