package models

import (
	"github.com/jinzhu/gorm"
	// "github.com/lib/pq/hstore"
	// "time"
)

type User struct {
	gorm.Model
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email" sql:"unique"`
	MobileNumber         string `json:"mobile_number" sql:"unique"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	Device               Device `gorm:"ForeignKey:DeviseToken"`
	DeviseToken          string `json:"devise_token"`
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
	StartLocation string               `json:"start_location"`
	Route         []map[string]float64 `json:"route" sql:"type:jsonb"`
	EndLocation   string               `json:"end_location"`
	User          User                 `gorm:"ForeignKey:UserID"`
	UserID        int                  `json:"user_id"`
}

// type Route struct {
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// 	Tracker   Tracker `gorm:"ForeignKey:TrackerID"`
// 	TrackerID int     `json:"tracker_id"`
// }

type Message struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
