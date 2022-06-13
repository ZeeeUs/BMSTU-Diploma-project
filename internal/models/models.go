package models

import (
	"io"
	"time"

	"github.com/jackc/pgx/pgtype"
)

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
	GroupCode string `json:"groupCode"`
}

type Student struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	GroupId int `json:"groupId"`
}

type Supervisor struct {
	Id     int `json:"id"`
	UserId int `json:"userId"`
}

type Table struct {
	Courses []CCourse `json:"courses"`
}

type CCourse struct {
	CourseId   int     `json:"courseId"`
	CourseName string  `json:"courseName"`
	Events     []Event `json:"events"`
}

type Event struct {
	EventId     int         `json:"eventId"`
	EventDate   time.Time   `json:"eventDate"`
	Deadline    time.Time   `json:"deadline"`
	Status      int         `json:"status"`
	EventName   string      `json:"eventName"`
	Description pgtype.Text `json:"description"`
	Files       []string    `json:"files"`
	Comments    []string    `json:"comment"`
}

type GroupByCourse struct {
	GroupId    int    `json:"groupId"`
	CourseId   int    `json:"courseId"`
	GroupCode  string `json:"groupCode"`
	Semester   int    `json:"semester"`
	CourseName string `json:"courseName"`
}

type StudentByGroup struct {
	StudentId  int    `json:"studentId"`
	UserId     int    `json:"userId"`
	GroupId    int    `json:"groupId"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
}

type File struct {
	File string `json:"file"`
}

type EventStatus struct {
	Status int `json:"status"`
}

type StudentEvent struct {
	Id          int      `json:"id"`
	EventId     int      `json:"eventId"`
	StudentId   int      `json:"studentId"`
	Status      int      `json:"status"`
	UploadFiles []string `json:"uploadFiles"`
	Comments    []string `json:"comment"`
}

type FileUnit struct {
	Payload     io.Reader
	PayloadName string
	PayloadSize int64
}
