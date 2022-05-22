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
	IsSuper    bool   `json:"is_super"`
}

type UpdateUser struct {
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
}

//type Supervisors struct {
//	Id         int      `json:"id"`
//	Email      string   `json:"email"`
//	Firstname  string   `json:"firstname"`
//	MiddleName string   `json:"middle_name"`
//	Lastname   string   `json:"lastname"`
//	Courses    []Course `json:"courses"`
//}

type Course struct {
	Id         int    `json:"id"`
	Semester   int    `json:"semester"`
	CourseName string `json:"course_name"`
}

type Group struct {
	Id        int    `json:"id"`
	GroupCode string `json:"group_code"`
}

type Student struct {
	Id         int    `json:"id"`
	Firstname  string `json:"firstname"`
	MiddleName string `json:"middleName"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Group      Group  `json:"group"`
}
