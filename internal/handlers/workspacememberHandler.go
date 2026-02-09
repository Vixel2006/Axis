package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils" // Import the utils package
	"github.com/gin-gonic/gin"
)

type WorkspaceMemberHandler struct {
	workspaceMemberService services.WorkspaceMemberService
}

func NewWorkspaceMemberHandler(wms services.WorkspaceMemberService) *WorkspaceMemberHandler {
	return &WorkspaceMemberHandler{
		workspaceMemberService: wms,
	}
}

func (h *WorkspaceMemberHandler) AddMemberToWorkspace(c *gin.Context) {
	workspaceIDStr := c.Param("workspaceID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var reqBody struct {
		UserID int             `json:"user_id"`
		Role   models.UserRole `json:"role"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workspaceMember, err := h.workspaceMemberService.AddMemberToWorkspace(c.Request.Context(), workspaceID, reqBody.UserID, reqBody.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member to workspace"})
		return
	}

	c.JSON(http.StatusCreated, workspaceMember)
}

func (h *WorkspaceMemberHandler) RemoveMemberFromWorkspace(c *gin.Context) {
	workspaceIDStr := c.Param("workspaceID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.workspaceMemberService.RemoveMemberFromWorkspace(c.Request.Context(), workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member from workspace"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *WorkspaceMemberHandler) GetWorkspaceMembers(c *gin.Context) {
	workspaceIDStr := c.Param("workspaceID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	members, err := h.workspaceMemberService.GetWorkspaceMembers(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workspace members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *WorkspaceMemberHandler) JoinWorkspace(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return // GetUserIDFromContext already handles the error response
	}

	workspaceIDStr := c.Param("workspaceID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	workspaceMember, err := h.workspaceMemberService.JoinWorkspace(c.Request.Context(), workspaceID, int(userID))
	if err != nil {
		// Handle specific errors from service layer
		if _, ok := err.(*services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join workspace"})
		return
	}

	c.JSON(http.StatusCreated, workspaceMember)
}
