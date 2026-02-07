package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Workspace struct {
	bun.BaseModel `bun:"table:workspaces,alias:w"`

	ID          int       `bun:",pk,autoincrement" json:"id"`
	Name        string    `bun:",notnull" json:"name"`
	Description *string   `bun:"" json:"description"`
	MaxMembers  *int      `bun:"" json:"max_members"`
	IsActive    bool      `bun:",notnull,default:false" json:"is_active"`
	CreatorID   int       `bun:",notnull" json:"creator_id"`
	CreatedAt   time.Time `bun:",nullzero,default:current_timestamp" json:"created_at"`

	Creator *User `bun:"rel:belongs-to,join:creator_id=id" json:"-"`
}
