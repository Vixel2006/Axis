package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type WorkspaceHandler struct {
	workspaceService services.WorkspaceService
	log              zerolog.Logger
}

func NewWorkspaceHandler(ws services.WorkspaceService, logger zerolog.Logger) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceService: ws,
		log:              logger,
	}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	h.log.Info().Msg("Handling CreateWorkspace request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in CreateWorkspace")
		return
	}

	var workspace models.Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for CreateWorkspace")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspace.CreatorID = int(userID)
	h.log.Debug().Int("creator_id", workspace.CreatorID).Str("workspace_name", workspace.Name).Msg("CreateWorkspace request body")

	createdWorkspace, err := h.workspaceService.CreateWorkspace(c.Request.Context(), &workspace)
	if err != nil {
		h.log.Error().Err(err).Int("creator_id", int(userID)).Str("workspace_name", workspace.Name).Msg("Failed to create workspace via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace"})
		return
	}

	h.log.Info().Int("workspace_id", createdWorkspace.ID).Int("creator_id", int(userID)).Msg("Workspace created successfully")
	c.JSON(http.StatusCreated, createdWorkspace)
}

func (h *WorkspaceHandler) GetWorkspaceByID(c *gin.Context) {
	h.log.Info().Msg("Handling GetWorkspaceByID request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in GetWorkspaceByID")
		return
	}

	idStr := c.Param("workspaceID")
	h.log.Debug().Str("workspaceID_param", idStr).Msg("Parsing workspace ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("workspaceID_param", idStr).Msg("Invalid workspace ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}
	h.log.Debug().Int("user_id", int(userID)).Int("workspace_id", id).Msg("Retrieving workspace by ID")

	workspace, err := h.workspaceService.GetWorkspaceByIDAuthorized(c.Request.Context(), int(userID), id)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("workspace_id", id).Msg("User forbidden from accessing workspace")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("workspace_id", id).Msg("Failed to retrieve workspace via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workspace"})
		return
	}

	if workspace == nil {
		h.log.Info().Int("workspace_id", id).Msg("Workspace not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	h.log.Info().Int("workspace_id", id).Int("user_id", int(userID)).Msg("Workspace retrieved successfully")
	c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	h.log.Info().Msg("Handling UpdateWorkspace request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in UpdateWorkspace")
		return
	}

	idStr := c.Param("workspaceID")
	h.log.Debug().Str("workspaceID_param", idStr).Msg("Parsing workspace ID for update")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("workspaceID_param", idStr).Msg("Invalid workspace ID format for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var workspace models.Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for UpdateWorkspace")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspace.ID = id
	h.log.Debug().Int("workspace_id", id).Int("user_id", int(userID)).Interface("request_body", workspace).Msg("UpdateWorkspace request body")

	updatedWorkspace, err := h.workspaceService.UpdateWorkspace(c.Request.Context(), int(userID), &workspace)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("workspace_id", id).Msg("User forbidden from updating workspace")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("workspace_id", id).Msg("Workspace not found for update")
			c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("workspace_id", id).Msg("Failed to update workspace via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workspace"})
		return
	}

	if updatedWorkspace == nil {
		h.log.Info().Int("workspace_id", id).Msg("Workspace not found for update (service returned nil)")
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	h.log.Info().Int("workspace_id", updatedWorkspace.ID).Int("user_id", int(userID)).Msg("Workspace updated successfully")
	c.JSON(http.StatusOK, updatedWorkspace)
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	h.log.Info().Msg("Handling DeleteWorkspace request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in DeleteWorkspace")
		return
	}

	idStr := c.Param("workspaceID")
	h.log.Debug().Str("workspaceID_param", idStr).Msg("Parsing workspace ID for deletion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("workspaceID_param", idStr).Msg("Invalid workspace ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}
	h.log.Debug().Int("workspace_id", id).Int("user_id", int(userID)).Msg("Attempting to delete workspace")

	err = h.workspaceService.DeleteWorkspace(c.Request.Context(), int(userID), id)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("workspace_id", id).Msg("User forbidden from deleting workspace")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("workspace_id", id).Msg("Workspace not found for deletion")
			c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("workspace_id", id).Msg("Failed to delete workspace via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workspace"})
		return
	}

	h.log.Info().Int("workspace_id", id).Int("user_id", int(userID)).Msg("Workspace deleted successfully")
	c.JSON(http.StatusNoContent, nil)
}

func (h *WorkspaceHandler) GetWorkspacesForUser(c *gin.Context) {
	h.log.Info().Msg("Handling GetWorkspacesForUser request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in GetWorkspacesForUser")
		return
	}
	h.log.Debug().Int("user_id_from_context", userID).Msg("Retrieving workspaces for user")

	workspaces, err := h.workspaceService.GetWorkspacesForUser(c.Request.Context(), int(userID))
	if err != nil {
		h.log.Error().Err(err).Int("user_id", int(userID)).Msg("Failed to retrieve workspaces for user via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workspaces for user"})
		return
	}

	h.log.Info().Int("user_id", int(userID)).Int("workspaces_count", len(workspaces)).Msg("Workspaces for user retrieved successfully")
	c.JSON(http.StatusOK, workspaces)
}
