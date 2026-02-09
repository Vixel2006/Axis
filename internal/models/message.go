package models

import (
	"github.com/uptrace/bun"
	"time"
)

type MessageType int

const (
	MessageTypeMessage MessageType = iota
	MessageTypeFileShare
	MessageTypeSystem
)

type Message struct {
	bun.BaseModel `bun:"table:messages,alias:m"`

	ID              int         `bun:",pk,autoincrement" json:"id"`
	ParentMessageID *int        `bun:"" json:"parent_message_id"` // For message threading/replies
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
}
