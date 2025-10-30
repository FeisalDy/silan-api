package userrole

import (
	"simple-go/internal/domain/role"
	"simple-go/internal/domain/user"
	"time"
)

type UserRole struct {
	UserID    string    `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	RoleID    string    `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	User      user.User `gorm:"foreignKey:UserID;references:ID"`
	Role      role.Role `gorm:"foreignKey:RoleID;references:ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
