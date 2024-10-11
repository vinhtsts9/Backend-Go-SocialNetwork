package po

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID       uuid.UUID `gorm:"column:uuid; type:int;not null;primaryKey; autoIncrement;"`
	RoleName string    `gorm:"column:role_name"`
	RoleNote string    `gorm:"column:role_note; type:text;"`
}

func (r *Role) TableName() string {
	return "go_db_user"
}
