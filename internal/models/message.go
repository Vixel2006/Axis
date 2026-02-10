package models

import (
	"encoding/json"
	"github.com/uptrace/bun"
	"time"
)

// ReactionBroadcastPayload is used to broadcast reaction changes to clients (used by service layer)
type ReactionBroadcastPayload struct {
	MessageID int    `json:"message_id"`
	UserID    int    `json:"user_id"`
	Emoji     string `json:"emoji"`
	Action    string `json:"action"` // "added" or "removed"
}

type BroadcastMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type MessageType int

const (
	MessageTypeMessage MessageType = iota
	MessageTypeFileShare
	MessageTypeSystem
)

type Message struct {
	bun.BaseModel `bun:"table:messages,alias:m"`

	ID              int         `bun:",pk,autoincrement" json:"id"`
	ParentMessageID *int        `bun:"" json:"parent_message_id"`
	Content         string      `bun:",notnull" json:"content"`
	MessageType     MessageType `bun:",notnull" json:"message_type"`
	MeetingID       int         `bun:",notnull" json:"meeting_id"`
	SenderID        int         `bun:",notnull" json:"sender_id"`
	IsEdited        bool        `bun:",notnull,default:false" json:"is_edited"`
	EditedAt        time.Time   `bun:",nullzero,default:current_timestamp" json:"edited_at"`
	CreatedAt       time.Time   `bun:",nullzero,default:current_timestamp" json:"created_at"`

	Sender        *User    `bun:"rel:belongs-to,join:sender_id=id"`
	Meeting       *Meeting `bun:"rel:belongs-to,join:meeting_id=id"`
	ParentMessage *Message `bun:"rel:belongs-to,join:parent_message_id=id"`

	Attachments []*Attachment `json:"attachments,omitempty" bun:"-"` // Non-persisted, populated for broadcasting
	Reactions   []*Reaction   `json:"reactions,omitempty" bun:"-"`   // Non-persisted, populated for broadcasting
}

type SendAttachmentDetails struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
	URL      string `json:"url"`
}

// WebSocket Message Models
type WSMessage struct {
	Type string      `json:"type"` // "message", "reaction", "join", "leave", "typing", "history", "error"
	Data interface{} `json:"data"`
}

type WSMessageData struct {
	ID        int                `json:"id,omitempty"`
	Content   string             `json:"content"`
	RoomID    int                `json:"room_id"`
	UserID    int                `json:"user_id"`
	Timestamp time.Time          `json:"timestamp"`
	Type      string             `json:"type,omitempty"` // "text", "file", "system"
	ReplyTo   *int               `json:"reply_to,omitempty"`
	Files     []WSAttachmentData `json:"files,omitempty"`
	User      *WSUserData        `json:"user,omitempty"`
}

type WSReactionData struct {
	MessageID int         `json:"message_id"`
	UserID    int         `json:"user_id"`
	Emoji     string      `json:"emoji"`
	Action    string      `json:"action"` // "add", "remove"
	Timestamp time.Time   `json:"timestamp"`
	User      *WSUserData `json:"user,omitempty"`
}

type WSTypingData struct {
	UserID   int         `json:"user_id"`
	RoomID   int         `json:"room_id"`
	IsTyping bool        `json:"is_typing"`
	User     *WSUserData `json:"user,omitempty"`
}

type WSRoomData struct {
	RoomID    int         `json:"room_id"`
	UserID    int         `json:"user_id"`
	Action    string      `json:"action"` // "join", "leave"
	Timestamp time.Time   `json:"timestamp"`
	User      *WSUserData `json:"user,omitempty"`
}

type WSHistoryData struct {
	RoomID   int             `json:"room_id"`
	Messages []WSMessageData `json:"messages"`
	HasMore  bool            `json:"has_more"`
	Offset   int             `json:"offset"`
}

type WSAttachmentData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

type WSUserData struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

type WSErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
