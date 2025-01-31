package repository

import "github.com/gofrs/uuid"

type User struct {
	UUID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Role     string    `gorm:"type:user_role;not null" json:"role"`
}

func (User) TableName() string {
	return "users"
}

type UserRepository interface {
	Login(username, password string) (bool, error)
	CreateUser(*User) error
	GetAll() ([]User, error)
	GetUser(identifier string) (*User, error)
	UpdateUser(uuid.UUID, *User) error
	DeleteUser(identifier string) error
}
