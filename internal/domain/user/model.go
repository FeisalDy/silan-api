package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Username  *string   `gorm:"type:varchar(100);unique" json:"username"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  *string   `gorm:"type:varchar(255)" json:"-"` // Never expose password in JSON
	AvatarURL *string   `gorm:"type:varchar(500)" json:"avatar_url"`
	Bio       *string   `gorm:"type:text" json:"bio"`
	Status    string    `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// CreateUserDTO is used for user creation
type CreateUserDTO struct {
	Username *string `json:"username"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Bio      *string `json:"bio"`
}

// UpdateUserDTO is used for user updates
type UpdateUserDTO struct {
	Username  *string `json:"username"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
	Status    *string `json:"status"`
}

// UserResponse is used for API responses (without sensitive data)
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

// ToResponse converts User to UserResponse
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
