package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Meeting struct {
	bun.BaseModel `bun:"table:meetings,alias:m"`

	ID          int         `bun:",pk,autoincrement" json:"id"`
	Name        string      `bun:",notnull" json:"name"`
	Description *string     `bun:"" json:"description"`
	ChannelID   int         `bun:",notnull" json:"channel_id"` // The channel where the meeting is held
	CreatorID   int         `bun:",notnull" json:"creator_id"`
	StartTime   time.Time   `bun:",notnull" json:"start_time"`
	EndTime     time.Time   `bun:",notnull" json:"end_time"`
	CreatedAt   time.Time   `bun:",nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time   `bun:",nullzero,default:current_timestamp" json:"updated_at"`

	Creator     *User      `bun:"rel:belongs-to,join:creator_id=id" json:"-"`
	Channel     *Channel   `bun:"rel:belongs-to,join:channel_id=id" json:"-"`
	Participants []*User    `bun:"m2m:meeting_members,join:Meeting=User" json:"-"` // Many-to-many through MeetingMember
	Messages    []*Message `bun:"rel:has-many,join:id=meeting_id" json:"-"` // Messages belonging to this meeting
}

type MeetingMember struct {
	bun.BaseModel `bun:"table:meeting_members,alias:mm"`

	MeetingID int       `bun:",pk" json:"meeting_id"`
	UserID    int       `bun:",pk" json:"user_id"`
	JoinedAt  time.Time `bun:",nullzero,default:current_timestamp" json:"joined_at"`

	Meeting *Meeting `bun:"rel:belongs-to,join:meeting_id=id" json:"-"`
	User    *User    `bun:"rel:belongs-to,join:user_id=id" json:"-"`
}
