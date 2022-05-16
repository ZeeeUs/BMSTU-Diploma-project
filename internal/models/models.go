package models

import "time"

type session struct {
	User   User
	Expire time.Time
}

type User struct {
	Firstname  string `json:"firstname"`
	MiddleName string `json:"middleName"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
