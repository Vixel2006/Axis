package models

import (
	"time"

	"github.com/uptrace/bun"
)

type UserRole int

const (
	Admin UserRole = iota
	Member
)

func (r UserRole) String() string {
	switch r {
	case Admin:
		return "admin"
	case Member:
		return "member"
	default:
		return "unknown"
	}
}

type WorkspaceMember struct {
	bun.BaseModel `bun:"table:workspace_members,alias:wm"`

	WorkspaceID int       `bun:",pk" json:"workspace_id"`
	UserID      int       `bun:",pk" json:"user_id"`
	Role        UserRole  `bun:",notnull" json:"role"`
	CreatedAt   time.Time `bun:",nullzero,default:current_timestamp" json:"created_at"`

	// Relationships
	Workspace *Workspace `bun:"rel:belongs-to,join:workspace_id=id"`
	User      *User      `bun:"rel:belongs-to,join:user_id=id"`
}
