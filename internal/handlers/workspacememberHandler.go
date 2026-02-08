package handlers

import (
	"net/http"
	"strconv"

	"axis/internal/models"
	"axis/internal/services"
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
