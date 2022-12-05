package models

import (
	"time"
)

// BaseGormModel mimixks GormModel but uses uuid's for ID, generated in go
type BaseGormModel struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// AccountModel
type Account struct {
	BaseGormModel
	Email            string    `sql:"type:varchar(255);not null;unique"`
	Name             string    `sql:"type:varchar(255);not null"`
	PhoneNumber      string    `sql:"type:varchar(255);not null"`
	ConfirmAndActive bool      `sql:"default:false"`
	MemberSince      time.Time `sql:"not null;DEFAULT:current_timestamp"`
	PasswordHash     string    `sql:"type:varchar(100)"`
	PasswordSalt     string    `sql:"type:binary(60);not null;default:''"`
	PhotoUrl         string    `sql:"type:varchar(255);not null"`
	Support          bool      `sql:"not null;default:false"`
}

func (c *Account) TableName() string {
	return "accounts"
}
