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

	ID              int         `bun:",pk,autoincrement" json:""`
	ParentMessageID *int        `bun:"" json:"parent_message_id"` // For message threading/replies
	Content         string      `bun:",notnull"`
	MessageType     MessageType `bun:",notnull" json:""`
	ChannelID       int         `bun:",notnull" json:""`
	SenderID        int         `bun:",notnull" json:""`
	IsEdited        bool        `bun:",notnull,default:false" json:""`
	EditedAt        time.Time   `bun:",nullzero,default:current_timestamp" json:""`
	CreatedAt       time.Time   `bun:",nullzero,default:current_timestamp" json:""`

	Sender        *User    `bun:"rel:belongs-to,join:sender_id=id"`
	Channel       *Channel `bun:"rel:belongs-to,join:channel_id=id"`
	ParentMessage *Message `bun:"rel:belongs-to,join:parent_message_id=id"`
}
