package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"axis/internal/models"
	"axis/internal/services"
	"axis/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now, consider securing this in production
		return true
	},
}

type ChatHandler struct {
	chatService services.MeetingChatService
	log         zerolog.Logger
}

func NewChatHandler(cs services.MeetingChatService, logger zerolog.Logger) *ChatHandler {
	return &ChatHandler{
		chatService: cs,
		log:         logger,
	}
}

func (h *ChatHandler) ServeMeetingChatWs(c *gin.Context) {
	meetingIDStr := c.Param("meeting_id")
	meetingID, err := strconv.Atoi(meetingIDStr)
	if err != nil {
		h.log.Error().Err(err).Str("meeting_id_str", meetingIDStr).Msg("Invalid meeting ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get user ID from context for WebSocket connection")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error().Err(err).Int("user_id", userID).Int("meeting_id", meetingID).Msg("Failed to upgrade connection to WebSocket")
		return
	}
	defer conn.Close()

	err = h.chatService.JoinMeetingChat(c.Request.Context(), meetingID, userID)
	if err != nil {
		h.log.Error().Err(err).Int("user_id", userID).Int("meeting_id", meetingID).Msg("Failed to join meeting chat")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to join chat"))
		return
	}

	client := &utils.Client{ID: userID, Conn: conn, Message: make(chan []byte, 256)}

	_ = h.chatService.GetOrCreateHubForMeeting(meetingID)

	h.chatService.RegisterClient(meetingID, client)
	defer h.chatService.UnregisterClient(meetingID, client)

	// Send join notification to all clients
	joinData := models.WSRoomData{
		RoomID:    meetingID,
		UserID:    userID,
		Action:    "join",
		Timestamp: time.Now(),
		User:      &models.WSUserData{ID: userID, Name: "", Username: ""}, // TODO: Get user details
	}
	h.broadcastWSMessage(meetingID, "room", joinData)

	go h.writePump(client, meetingID)
	h.readPump(c.Request.Context(), client, meetingID)
}

func (h *ChatHandler) readPump(ctx context.Context, client *utils.Client, meetingID int) {
	defer func() {
		h.chatService.UnregisterClient(meetingID, client)
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(512)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		_, messageBytes, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				h.log.Error().Err(err).Int("user_id", client.ID).Int("meeting_id", meetingID).Msg("Unexpected close error while reading message")
			} else {
				h.log.Info().Err(err).Int("user_id", client.ID).Int("meeting_id", meetingID).Msg("Client disconnected or read error")
			}
			break
		}

		// Parse using new simplified WS message format
		var wsMessage models.WSMessage
		if err := json.Unmarshal(messageBytes, &wsMessage); err != nil {
			h.sendError(client, "INVALID_FORMAT", "Invalid message format", err.Error())
			continue
		}

		switch wsMessage.Type {
		case "message":
			h.handleWSMessage(ctx, client, meetingID, wsMessage)
		case "reaction":
			h.handleWSReaction(ctx, client, meetingID, wsMessage)
		case "typing":
			h.handleWSTyping(ctx, client, meetingID, wsMessage)
		case "history":
			h.handleWSHistory(ctx, client, meetingID, wsMessage)
		default:
			h.sendError(client, "UNKNOWN_TYPE", "Unknown message type", fmt.Sprintf("Type: %s", wsMessage.Type))
		}
	}
}

func (h *ChatHandler) handleWSMessage(ctx context.Context, client *utils.Client, meetingID int, wsMessage models.WSMessage) {
	// Parse message data
	dataBytes, err := json.Marshal(wsMessage.Data)
	if err != nil {
		h.sendError(client, "INVALID_MESSAGE_DATA", "Invalid message data", err.Error())
		return
	}

	var msgData models.WSMessageData
	if err := json.Unmarshal(dataBytes, &msgData); err != nil {
		h.sendError(client, "INVALID_MESSAGE_DATA", "Invalid message data", err.Error())
		return
	}

	// Validate room ID matches
	if msgData.RoomID != meetingID {
		h.sendError(client, "ROOM_MISMATCH", "Room ID mismatch", fmt.Sprintf("Expected: %d, Got: %d", meetingID, msgData.RoomID))
		return
	}

	// Convert to internal message type
	messageType := models.MessageTypeMessage
	if msgData.Type == "file" {
		messageType = models.MessageTypeFileShare
	} else if msgData.Type == "system" {
		messageType = models.MessageTypeSystem
	}

	// Convert attachments
	var attachments []models.SendAttachmentDetails
	for _, file := range msgData.Files {
		attachments = append(attachments, models.SendAttachmentDetails{
			FileName: file.Name,
			FileType: file.Type,
			FileSize: file.Size,
			URL:      file.URL,
		})
	}

	savedMessage, err := h.chatService.SendMessage(ctx, meetingID, client.ID, msgData.ReplyTo, msgData.Content, messageType, attachments)
	if err != nil {
		h.sendError(client, "SEND_FAILED", "Failed to send message", err.Error())
		return
	}

	// Create simplified response message
	responseData := models.WSMessageData{
		ID:        savedMessage.ID,
		Content:   savedMessage.Content,
		RoomID:    meetingID,
		UserID:    client.ID,
		Timestamp: savedMessage.CreatedAt,
		Type:      string(msgData.Type),
		ReplyTo:   msgData.ReplyTo,
		Files:     msgData.Files,
		User:      &models.WSUserData{ID: client.ID, Name: "", Username: ""}, // TODO: Get user details
	}

	// Broadcast simplified message
	h.broadcastWSMessage(meetingID, "message", responseData)
}

func (h *ChatHandler) handleWSReaction(ctx context.Context, client *utils.Client, meetingID int, wsMessage models.WSMessage) {
	dataBytes, err := json.Marshal(wsMessage.Data)
	if err != nil {
		h.sendError(client, "INVALID_REACTION_DATA", "Invalid reaction data", err.Error())
		return
	}

	var reactionData models.WSReactionData
	if err := json.Unmarshal(dataBytes, &reactionData); err != nil {
		h.sendError(client, "INVALID_REACTION_DATA", "Invalid reaction data", err.Error())
		return
	}

	var errReaction error
	if reactionData.Action == "add" {
		_, errReaction = h.chatService.AddReaction(ctx, reactionData.MessageID, client.ID, reactionData.Emoji)
	} else if reactionData.Action == "remove" {
		_, errReaction = h.chatService.RemoveReaction(ctx, reactionData.MessageID, client.ID, reactionData.Emoji)
	} else {
		h.sendError(client, "INVALID_REACTION_ACTION", "Invalid reaction action", fmt.Sprintf("Action: %s", reactionData.Action))
		return
	}

	if errReaction != nil {
		h.sendError(client, "REACTION_FAILED", "Failed to handle reaction", errReaction.Error())
		return
	}

	// Create simplified response
	responseData := models.WSReactionData{
		MessageID: reactionData.MessageID,
		UserID:    client.ID,
		Emoji:     reactionData.Emoji,
		Action:    reactionData.Action,
		Timestamp: time.Now(),
		User:      &models.WSUserData{ID: client.ID, Name: "", Username: ""}, // TODO: Get user details
	}

	h.broadcastWSMessage(meetingID, "reaction", responseData)
}

func (h *ChatHandler) handleWSTyping(ctx context.Context, client *utils.Client, meetingID int, wsMessage models.WSMessage) {
	dataBytes, err := json.Marshal(wsMessage.Data)
	if err != nil {
		h.sendError(client, "INVALID_TYPING_DATA", "Invalid typing data", err.Error())
		return
	}

	var typingData models.WSTypingData
	if err := json.Unmarshal(dataBytes, &typingData); err != nil {
		h.sendError(client, "INVALID_TYPING_DATA", "Invalid typing data", err.Error())
		return
	}

	// Validate room ID
	if typingData.RoomID != meetingID {
		h.sendError(client, "ROOM_MISMATCH", "Room ID mismatch", fmt.Sprintf("Expected: %d, Got: %d", meetingID, typingData.RoomID))
		return
	}

	// Set user ID and broadcast
	typingData.UserID = client.ID
	typingData.User = &models.WSUserData{ID: client.ID, Name: "", Username: ""} // TODO: Get user details

	h.broadcastWSMessage(meetingID, "typing", typingData)
}

func (h *ChatHandler) handleWSHistory(ctx context.Context, client *utils.Client, meetingID int, wsMessage models.WSMessage) {
	dataBytes, err := json.Marshal(wsMessage.Data)
	if err != nil {
		h.sendError(client, "INVALID_HISTORY_DATA", "Invalid history data", err.Error())
		return
	}

	var historyData models.WSHistoryData
	if err := json.Unmarshal(dataBytes, &historyData); err != nil {
		h.sendError(client, "INVALID_HISTORY_DATA", "Invalid history data", err.Error())
		return
	}

	// Validate room ID
	if historyData.RoomID != meetingID {
		h.sendError(client, "ROOM_MISMATCH", "Room ID mismatch", fmt.Sprintf("Expected: %d, Got: %d", meetingID, historyData.RoomID))
		return
	}

	// Get messages with pagination (default: 50, starting from historyData.Offset)
	limit := 50
	offset := historyData.Offset
	if offset < 0 {
		offset = 0
	}

	messages, err := h.chatService.GetMeetingMessages(ctx, meetingID, limit, offset)
	if err != nil {
		h.sendError(client, "HISTORY_FAILED", "Failed to get message history", err.Error())
		return
	}

	// Convert to WS format
	var wsMessages []models.WSMessageData
	for _, msg := range messages {
		var files []models.WSAttachmentData
		for _, att := range msg.Attachments {
			files = append(files, models.WSAttachmentData{
				ID:   att.ID,
				Name: att.FileName,
				Type: att.FileType,
				Size: att.FileSize,
				URL:  att.URL,
			})
		}

		wsMessages = append(wsMessages, models.WSMessageData{
			ID:        msg.ID,
			Content:   msg.Content,
			RoomID:    meetingID,
			UserID:    msg.SenderID,
			Timestamp: msg.CreatedAt,
			Type:      "text", // TODO: Convert from MessageType
			ReplyTo:   msg.ParentMessageID,
			Files:     files,
			User:      &models.WSUserData{ID: msg.SenderID, Name: "", Username: ""}, // TODO: Get user details
		})
	}

	// Create response
	responseData := models.WSHistoryData{
		RoomID:   meetingID,
		Messages: wsMessages,
		HasMore:  len(messages) == limit, // If we got full limit, there might be more
		Offset:   offset + len(messages),
	}

	// Send only to requesting client
	h.sendWSMessage(client, "history", responseData)
}

func (h *ChatHandler) writePump(client *utils.Client, meetingID int) {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Message:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				h.log.Error().Err(err).Int("user_id", client.ID).Int("meeting_id", meetingID).Msg("Failed to get next writer for WebSocket")
				return
			}
			w.Write(message)

			n := len(client.Message)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-client.Message)
			}

			if err := w.Close(); err != nil {
				h.log.Error().Err(err).Int("user_id", client.ID).Int("meeting_id", meetingID).Msg("Failed to close writer for WebSocket")
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				h.log.Error().Err(err).Int("user_id", client.ID).Int("meeting_id", meetingID).Msg("Failed to send ping message")
				return
			}
		}
	}
}

func (h *ChatHandler) sendError(client *utils.Client, code, message, details string) {
	errorData := models.WSErrorData{
		Code:    code,
		Message: message,
		Details: details,
	}
	h.sendWSMessage(client, "error", errorData)
}

func (h *ChatHandler) sendWSMessage(client *utils.Client, msgType string, data interface{}) {
	wsMessage := models.WSMessage{
		Type: msgType,
		Data: data,
	}
	if err := client.Conn.WriteJSON(wsMessage); err != nil {
		h.log.Error().Err(err).Int("user_id", client.ID).Msg("Failed to send WebSocket message")
	}
}

func (h *ChatHandler) broadcastWSMessage(meetingID int, msgType string, data interface{}) {
	wsMessage := models.WSMessage{
		Type: msgType,
		Data: data,
	}
	if messageBytes, err := json.Marshal(wsMessage); err != nil {
		h.log.Error().Err(err).Str("message_type", msgType).Int("meeting_id", meetingID).Msg("Failed to marshal broadcast message")
	} else {
		h.chatService.BroadcastMessage(meetingID, messageBytes)
	}
}
