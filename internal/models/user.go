package models

import (
	"time"

	"github.com/uptrace/bun"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID       int64          `bun:",pk,autoincrement" json:"id"`
	Name     string         `bun:",notnull" json:"name"`
	Email    string         `bun:",unique,notnull" json:"email"`
	Settings map[string]any `bun:",type:jsonb" json:"settings"`
	Status   UserStatus     `bun:",type:varchar(20),default:'active'" json:"status"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`
}
