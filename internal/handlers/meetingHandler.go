package handlers

import (
	"net/http"
	"strconv"
	"time"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils"
	"github.com/gin-gonic/gin"
)

type MeetingHandler struct {
	meetingService services.MeetingService
}

func NewMeetingHandler(ms services.MeetingService) *MeetingHandler {
	return &MeetingHandler{
		meetingService: ms,
	}
}

type CreateMeetingRequest struct {
	Name          string    `json:"name" binding:"required"`
	Description   *string   `json:"description"`
	ChannelID     int       `json:"channel_id" binding:"required"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
	ParticipantIDs []int     `json:"participant_ids"`
}

func (h *MeetingHandler) CreateMeeting(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	var req CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meeting := &models.Meeting{
		Name:        req.Name,
		Description: req.Description,
		ChannelID:   req.ChannelID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	createdMeeting, err := h.meetingService.CreateMeeting(c.Request.Context(), meeting, int(userID), req.ParticipantIDs)
	if err != nil {
		if err.Error() == "channel not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meeting"})
		return
	}

	c.JSON(http.StatusCreated, createdMeeting)
}

func (h *MeetingHandler) GetMeetingByID(c *gin.Context) {
	idStr := c.Param("meetingID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	meeting, err := h.meetingService.GetMeetingByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meeting"})
		return
	}

	if meeting == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	c.JSON(http.StatusOK, meeting)
}

func (h *MeetingHandler) GetMeetingsByChannelID(c *gin.Context) {
	channelIDStr := c.Param("channelID")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	meetings, err := h.meetingService.GetMeetingsByChannelID(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meetings for channel"})
		return
	}

	c.JSON(http.StatusOK, meetings)
}

type UpdateMeetingRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description *string   `json:"description"`
	ChannelID   int       `json:"channel_id" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
}

func (h *MeetingHandler) UpdateMeeting(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	idStr := c.Param("meetingID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	var req UpdateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meeting := &models.Meeting{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		ChannelID:   req.ChannelID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	updatedMeeting, err := h.meetingService.UpdateMeeting(c.Request.Context(), meeting, int(userID))
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "meeting not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meeting"})
		return
	}

	c.JSON(http.StatusOK, updatedMeeting)
}

func (h *MeetingHandler) DeleteMeeting(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	idStr := c.Param("meetingID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	err = h.meetingService.DeleteMeeting(c.Request.Context(), id, int(userID))
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "meeting not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meeting"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type AddRemoveParticipantRequest struct {
	ParticipantID int `json:"participant_id" binding:"required"`
}

func (h *MeetingHandler) AddParticipant(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	meetingIDStr := c.Param("meetingID")
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	var req AddRemoveParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.meetingService.AddParticipant(c.Request.Context(), meetingID, int(userID), req.ParticipantID)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "meeting not found" || err.Error() == "participant user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant added successfully"})
}

func (h *MeetingHandler) RemoveParticipant(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	meetingIDStr := c.Param("meetingID")
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	var req AddRemoveParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.meetingService.RemoveParticipant(c.Request.Context(), meetingID, int(userID), req.ParticipantID)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "meeting not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove participant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant removed successfully"})
}
