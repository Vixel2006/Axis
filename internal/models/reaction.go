package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Reaction struct {
	bun.BaseModel `bun:"table:reactions,alias:r"`

	ID        int       `bun:",pk,autoincrement" json:"id"`
	MessageID int       `bun:",notnull" json:"message_id"`
	UserID    int       `bun:",notnull" json:"user_id"`
	Emoji     string    `bun:",notnull" json:"emoji"`
	CreatedAt time.Time `bun:",nullzero,default:current_timestamp" json:"created_at"`

	Message *Message `bun:"rel:belongs-to,join:message_id=id"`
	User    *User    `bun:"rel:belongs-to,join:user_id=id"`
}
