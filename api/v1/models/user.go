package models

import (
	"github.com/jinzhu/gorm"
	// "github.com/lib/pq/hstore"
	// "time"
)

type User struct {
	gorm.Model
	FirstName            string `json:"first_name" validate:"required"`
	LastName             string `json:"last_name"`
	Email                string `json:"email" sql:"unique" validate:"required,email"`
	MobileNumber         string `json:"mobile_number" sql:"unique" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required"`
	Device               Device `gorm:"ForeignKey:DeviseToken"`
	DeviseToken          string `json:"devise_token" validate:"required"`
}

type Session struct {
	gorm.Model
	User         User   `gorm:"ForeignKey:UserID"`
	UserID       int    `json:"user_id"`
	MobileNumber string `json:"mobile_number"`
	Device       Device `gorm:"ForeignKey:DeviseToken"`
	DeviseToken  string `json:"devise_token"`
}

type Device struct {
	gorm.Model
	Token string `json:"token"`
}

type Tracker struct {
	gorm.Model
	StartLocation string `json:"start_location"`
	// StartTime     time.Time `json:"start_time"`
	Routes string `json:"routes"`
	// Routes map[string]interface{} `json:"route" sql:"type:jsonb`
	// EndTime       time.Time `json:"end_time"`
	EndLocation string `json:"end_location"`
	User        User   `gorm:"ForeignKey:UserID"`
	UserID      int    `json:"user_id"`
}

type Message struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string  `json:"error"`
}
