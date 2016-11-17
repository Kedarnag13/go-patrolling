package models

import (
	"time"
)

type User struct {
	Id                   string `json:"id"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email" sql:"unique"`
	MobileNumber         string `json:"mobile_number" sql:"unique"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	DeviseToken          string `json:"devise_token"`
}

type Session struct {
	Id           string `json:"id"`
	MobileNumber string `json:"mobile_number"`
	Password     string `json:"password"`
	UserID       string `json:"user_id"`
	DeviseToken  string `json:"devise_token"`
}

type Device struct {
	Id        string
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
}

type Tracker struct {
	Id            string               `json:"id"`
	StartLocation string               `json:"start_location"`
	StartTime     time.Time            `json:"start_time"`
	Routes        []map[string]float64 `json:"routes"`
	EndTime       time.Time            `json:"end_time"`
	EndLocation   string               `json:"end_location"`
	MobileNumber  string               `json:"mobile_number"`
}

type Message struct {
	User    User    `json:"user"`
	Session Session `json:"session"`
	Device  Device  `json:device`
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Error   string  `json:"error"`
}
