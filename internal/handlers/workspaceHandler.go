package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils"
	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	workspaceService services.WorkspaceService
}

func NewWorkspaceHandler(ws services.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceService: ws,
	}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	var workspace models.Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspace.CreatorID = int(userID) // Set the CreatorID from the context

	createdWorkspace, err := h.workspaceService.CreateWorkspace(c.Request.Context(), &workspace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workspace"})
		return
	}

	c.JSON(http.StatusCreated, createdWorkspace)
}

func (h *WorkspaceHandler) GetWorkspaceByID(c *gin.Context) {
	idStr := c.Param("workspaceID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	workspace, err := h.workspaceService.GetWorkspaceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workspace"})
		return
	}

	if workspace == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	idStr := c.Param("workspaceID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var workspace models.Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspace.ID = id // Ensure the ID from the URL is used

	// Pass userID to the service for authorization
	updatedWorkspace, err := h.workspaceService.UpdateWorkspace(c.Request.Context(), int(userID), &workspace)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workspace"})
		return
	}

	if updatedWorkspace == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workspace not found"})
		return
	}

	c.JSON(http.StatusOK, updatedWorkspace)
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	idStr := c.Param("workspaceID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	// Pass userID to the service for authorization
	err = h.workspaceService.DeleteWorkspace(c.Request.Context(), int(userID), id)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workspace"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *WorkspaceHandler) GetWorkspacesForUser(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	workspaces, err := h.workspaceService.GetWorkspacesForUser(c.Request.Context(), int(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workspaces for user"})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}
