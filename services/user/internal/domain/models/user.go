package models

type User struct {
	UID       int64  `json:"uid"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=255"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type Session struct {
	UID   int64  `json:"uid"`
	Token string `json:"token"`
}
