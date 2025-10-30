package user

import (
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/role"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string        `gorm:"type:uuid;primaryKey"`
	Username  *string       `gorm:"type:varchar(100);unique"`
	Email     string        `gorm:"type:varchar(255);unique;not null"`
	Password  *string       `gorm:"type:varchar(255)"`
	AvatarURL *string       `gorm:"type:varchar(500)"`
	Bio       *string       `gorm:"type:text"`
	Status    string        `gorm:"type:varchar(50);default:'active'"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
	Roles     []role.Role   `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID"`
	Novels    []novel.Novel `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}

type CreateUserDTO struct {
	Username *string `json:"username"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Bio      *string `json:"bio"`
}

type UpdateUserDTO struct {
	Username  *string `json:"username"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
	Status    *string `json:"status"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  *string   `json:"username"`
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url"`
	Bio       *string   `json:"bio"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Bio:       u.Bio,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
