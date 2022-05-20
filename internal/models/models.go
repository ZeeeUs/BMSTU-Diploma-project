package models

type Session struct {
	Cookie string
	Id     int
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id         int    `json:"id"`
	Firstname  string `json:"firstname"`
	MiddleName string `json:"middleName"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	PassStatus bool   `json:"pass_status"`
}

type UpdateUser struct {
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
}
