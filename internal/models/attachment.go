package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Attachment struct {
	bun.BaseModel `bun:"table:attachments,alias:a"`

	ID        int       `bun:",pk,autoincrement" json:"id"`
	MessageID int       `bun:",notnull" json:"message_id"`
	UserID    int       `bun:",notnull" json:"user_id"`
	FileName  string    `bun:",notnull" json:"file_name"`
	FileType  string    `bun:",notnull" json:"file_type"`
	FileSize  int64     `bun:",notnull" json:"file_size"`
	URL       string    `bun:",notnull" json:"url"`
	CreatedAt time.Time `bun:",nullzero,default:current_timestamp" json:"created_at"`

	Message *Message `bun:"rel:belongs-to,join:message_id=id"`
	User    *User    `bun:"rel:belongs-to,join:user_id=id"`
}
