package dto

import ()

// These structs define the expected JSON request bodies for authentication-related endpoints 
// such as registration, login, and token refresh. 
// They are used to decode incoming JSON data into 
// Go structs that can be easily processed by the handlers and service layer.
type RegisterRequest struct {
	Email  string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
