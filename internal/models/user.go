package models

import (
	"github.com/uptrace/bun"
	"time"
)

type UserStatus int

const (
	Active = iota
	Buzy
	AFK
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          int        `bun:",pk,autoincrement" json:"id"`
	Name        string     `bun:",notnull" json:"name"`
	Username    string     `bun:",notnull,unique" json:"username"`
	Email       string     `bun:",notnull,unique" json:"email"`
	Password    string     `bun:",notnull" json:"-"`
	Status      UserStatus `bun:",notnull" json:"status"`
	Timezone    string     `bun:"" json:"timezone"`
	Locale      string     `bun:",notnull" json:"locale"`
	IsVerified  bool       `bun:",notnull,default:false" json:"is_verified"`
	LastLoginAt *time.Time `bun:",nullzero" json:"last_login_at"`
	CreatedAt   time.Time  `bun:",nullzero,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time  `bun:",nullzero,default:current_timestamp" json:"updated_at"`
}

type RegisterModel struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Timezone string `json:"timezone"`
	Locale   string `json:"locale"`
}

type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUser struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Timezone *string `json:"timezone"`
	Locale   *string `json:"locale"`
}
