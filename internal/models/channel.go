package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ChannelType int

const (
	ChannelTypePrivate ChannelType = iota
	ChannelTypePublic
	ChannelTypeDM
)

type Channel struct {
	bun.BaseModel `bun:"table:channels,alias:c"`

	ID          int         `bun:",pk,autoincrement" json:"id"`
	Name        string      `bun:",notnull" json:"name"`
	Description *string     `bun:"" json:"description"`
	ChannelType ChannelType `bun:",notnull" json:"channel_type"`
	WorkspaceID int         `bun:",notnull" json:"workspace_id"`
	IsArchieved bool        `bun:",notnull,default:false" json:"is_archived"`
	CreatorID   int         `bun:",notnull" json:"creator_id"`
	CreatedAt   time.Time   `bun:",nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time   `bun:",nullzero,default:current_timestamp" json:"updated_at"`

	Creator   *User      `bun:"rel:belongs-to,join:creator_id=id" json:"-"`
	Workspace *Workspace `bun:"rel:belongs-to,join:workspace_id=id" json:"-"`
}
