package models

type User struct {
	UID       int64  `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Session struct {
	UID   int64  `json:"uid"`
	Token string `json:"token"`
}
