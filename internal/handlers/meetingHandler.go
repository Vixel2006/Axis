package handlers

import (
	"net/http"
	"strconv"
	"time"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type MeetingHandler struct {
	meetingService services.MeetingService
	log            zerolog.Logger
}

func NewMeetingHandler(ms services.MeetingService, logger zerolog.Logger) *MeetingHandler {
	return &MeetingHandler{
		meetingService: ms,
		log:            logger,
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
	h.log.Info().Msg("Handling CreateMeeting request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in CreateMeeting")
		return
	}

	var req CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for CreateMeeting")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Interface("request_body", req).Msg("CreateMeeting request body")

	meeting := &models.Meeting{
		Name:        req.Name,
		Description: req.Description,
		ChannelID:   req.ChannelID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	createdMeeting, err := h.meetingService.CreateMeeting(c.Request.Context(), meeting, int(userID), req.ParticipantIDs)
	if err != nil {
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("channel_id", req.ChannelID).Msg("Channel not found for meeting creation")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("creator_id", int(userID)).Str("meeting_name", req.Name).Msg("Failed to create meeting via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meeting"})
		return
	}

	h.log.Info().Int("meeting_id", createdMeeting.ID).Int("creator_id", int(userID)).Msg("Meeting created successfully")
	c.JSON(http.StatusCreated, createdMeeting)
}

func (h *MeetingHandler) GetMeetingByID(c *gin.Context) {
	h.log.Info().Msg("Handling GetMeetingByID request")
	idStr := c.Param("meetingID")
	h.log.Debug().Str("meetingID_param", idStr).Msg("Parsing meeting ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("meetingID_param", idStr).Msg("Invalid meeting ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	meeting, err := h.meetingService.GetMeetingByID(c.Request.Context(), id)
	if err != nil {
		h.log.Error().Err(err).Int("meeting_id", id).Msg("Failed to retrieve meeting via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meeting"})
		return
	}

	if meeting == nil {
		h.log.Info().Int("meeting_id", id).Msg("Meeting not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting not found"})
		return
	}

	h.log.Info().Int("meeting_id", id).Msg("Meeting retrieved successfully")
	c.JSON(http.StatusOK, meeting)
}

func (h *MeetingHandler) GetMeetingsByChannelID(c *gin.Context) {
	h.log.Info().Msg("Handling GetMeetingsByChannelID request")
	channelIDStr := c.Param("channelID")
	h.log.Debug().Str("channelID_param", channelIDStr).Msg("Parsing channel ID")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("channelID_param", channelIDStr).Msg("Invalid channel ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	meetings, err := h.meetingService.GetMeetingsByChannelID(c.Request.Context(), channelID)
	if err != nil {
		h.log.Error().Err(err).Int("channel_id", channelID).Msg("Failed to retrieve meetings for channel via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve meetings for channel"})
		return
	}

	h.log.Info().Int("channel_id", channelID).Int("meetings_count", len(meetings)).Msg("Meetings for channel retrieved successfully")
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
	h.log.Info().Msg("Handling UpdateMeeting request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in UpdateMeeting")
		return
	}

	idStr := c.Param("meetingID")
	h.log.Debug().Str("meetingID_param", idStr).Msg("Parsing meeting ID for update")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("meetingID_param", idStr).Msg("Invalid meeting ID format for update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	var req UpdateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for UpdateMeeting")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Int("meeting_id", id).Interface("request_body", req).Msg("UpdateMeeting request body")

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
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("meeting_id", id).Msg("User forbidden from updating meeting")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("meeting_id", id).Msg("Meeting not found for update")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("meeting_id", id).Msg("Failed to update meeting via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meeting"})
		return
	}

	h.log.Info().Int("meeting_id", updatedMeeting.ID).Int("user_id", int(userID)).Msg("Meeting updated successfully")
	c.JSON(http.StatusOK, updatedMeeting)
}

func (h *MeetingHandler) DeleteMeeting(c *gin.Context) {
	h.log.Info().Msg("Handling DeleteMeeting request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in DeleteMeeting")
		return
	}

	idStr := c.Param("meetingID")
	h.log.Debug().Str("meetingID_param", idStr).Msg("Parsing meeting ID for deletion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Error().Err(err).Str("meetingID_param", idStr).Msg("Invalid meeting ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	err = h.meetingService.DeleteMeeting(c.Request.Context(), id, int(userID))
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("meeting_id", id).Msg("User forbidden from deleting meeting")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("meeting_id", id).Msg("Meeting not found for deletion")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("meeting_id", id).Msg("Failed to delete meeting via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meeting"})
		return
	}

	h.log.Info().Int("meeting_id", id).Int("user_id", int(userID)).Msg("Meeting deleted successfully")
	c.JSON(http.StatusNoContent, nil)
}

type AddRemoveParticipantRequest struct {
	ParticipantID int `json:"participant_id" binding:"required"`
}

func (h *MeetingHandler) AddParticipant(c *gin.Context) {
	h.log.Info().Msg("Handling AddParticipant request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in AddParticipant")
		return
	}

	meetingIDStr := c.Param("meetingID")
	h.log.Debug().Str("meetingID_param", meetingIDStr).Msg("Parsing meeting ID for adding participant")
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("meetingID_param", meetingIDStr).Msg("Invalid meeting ID format for adding participant")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	var req AddRemoveParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for AddParticipant")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Int("meeting_id", meetingID).Int("participant_id_req", req.ParticipantID).Msg("AddParticipant request body")

	err = h.meetingService.AddParticipant(c.Request.Context(), meetingID, int(userID), req.ParticipantID)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("meeting_id", meetingID).Int("participant_id", req.ParticipantID).Msg("User forbidden from adding participant")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("meeting_id", meetingID).Int("participant_id", req.ParticipantID).Msg("Meeting or participant user not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("meeting_id", meetingID).Int("participant_id", req.ParticipantID).Msg("Failed to add participant via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add participant"})
		return
	}

	h.log.Info().Int("meeting_id", meetingID).Int("participant_id", req.ParticipantID).Msg("Participant added successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Participant added successfully"})
}

func (h *MeetingHandler) RemoveParticipant(c *gin.Context) {
	h.log.Info().Msg("Handling RemoveParticipant request")
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context in RemoveParticipant")
		return
	}

	meetingIDStr := c.Param("meetingID")
	h.log.Debug().Str("meetingID_param", meetingIDStr).Msg("Parsing meeting ID for removing participant")
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("meetingID_param", meetingIDStr).Msg("Invalid meeting ID format for removing participant")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	var req AddRemoveParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("Failed to bind JSON for RemoveParticipant")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Debug().Int("meeting_id", meetingID).Int("participant_id_req", req.ParticipantID).Msg("RemoveParticipant request body")

	err = h.meetingService.RemoveParticipant(c.Request.Context(), meetingID, int(userID), req.ParticipantID)
	if err != nil {
		if _, ok := err.(*services.ForbiddenError); ok {
			h.log.Warn().Err(err).Int("user_id", int(userID)).Int("meeting_id", meetingID).Int("participant_id", req.ParticipantID).Msg("User forbidden from removing participant")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*services.NotFoundError); ok {
			h.log.Warn().Err(err).Int("meeting_id", meetingID).Msg("Meeting not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Int("user_id", int(userID)).Int("meeting_id", meetingID).Int("participant_id", req.ParticipantID).Msg("Failed to remove participant via service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove participant"})
		return
	}

	h.log.Info().Int("meeting_id", meetingID).Int("participant_id", req.ParticipantID).Msg("Participant removed successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Participant removed successfully"})
}
