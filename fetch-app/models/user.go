package models

import (
	"strings"

	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey" json:"id"`
	FullName  string `json:"fullName"`
	Email     string `gorm:"index;unique;type:varchar(255)" json:"email"`
	CreatedAt uint   `gorm:"autoCreateTime;<-:create" json:"createdAt"`
	UpdatedAt uint   `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) setAttr(tx *gorm.DB) {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	if u.Email == "" {
		tx.Statement.Omits = append(tx.Statement.Omits, "email")
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.setAttr(tx)
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.setAttr(tx)
	return
}
