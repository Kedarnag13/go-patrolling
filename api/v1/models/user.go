package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email"`
	MobileNumber         string `json:"mobile_number"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type Session struct {
	gorm.Model
	User        User   `gorm:"ForeignKey:UserID"`
	UserID      uint   `json:"user_id"`
	Device      Device `gorm:"ForeignKey:DeviseToken"`
	DeviseToken string `json:"devise_token"`
}

type Device struct {
	gorm.Model
	Token string `json:"token"`
}

type Message struct {
	Success bool
	Message string
	Error   string
}
