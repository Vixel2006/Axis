package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ChannelMember struct {
	bun.BaseModel `bun:"table:channel_members,alias:cm"`

	ChannelID         int       `bun:",pk" json:"channel_id"`
	UserID            int       `bun:",pk" json:"user_id"`
	JoinedAt          time.Time `bun:",nullzero,default:current_timestamp" json:"joined_at"`
	LastReadMessageID *int      `bun:"" json:"last_read_message_id"`

	Channel *Channel `bun:"rel:belongs-to,join:channel_id=id"`
	User    *User    `bun:"rel:belongs-to,join:user_id=id"`
}
