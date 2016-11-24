package models

type User struct {
	Id                   string `json:"id"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email"`
	MobileNumber         string `json:"mobile_number"`
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
	StartTime     string               `json:"start_time"`
	Routes        []map[string]float64 `json:"routes"`
	EndTime       string               `json:"end_time"`
	EndLocation   string               `json:"end_location"`
	MobileNumber  string               `json:"mobile_number"`
	CreatedAt     string               `json:"created_at"`
}

type Message struct {
	User    User                   `json:"user"`
	Session Session                `json:"session"`
	Device  Device                 `json:"device"`
	Tracker map[string]interface{} `json:"tracker"`
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Error   string                 `json:"error"`
}
