package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

var (
	conn  *pgx.Conn
	pgErr pgx.PgError
)

type session struct {
	Username string
	Expire   time.Time
}

//type Credentials struct {
//	Password string `json:"password"`
//	Username string `json:"username"`
//}

type User struct {
	Firstname  string `json:"firstname"`
	MiddleName string `json:"middleName"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var sessions = map[string]session{}

func (s session) isExpired() bool {
	return s.Expire.Before(time.Now())
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//expectedPassword, ok := users[user.Username]
	//if !ok || expectedPassword != cred.Password {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(3 * time.Hour)

	//sessions[sessionToken] = session{
	//	Username: cred.Username,
	//	Expire:   expiresAt,
	//}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionToken",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	//log.Println(sessionToken)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	user := &User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = conn.Query("insert into dashboard.users (password, pass_status, firstname, middle_name, lastname, email)"+
		" values($1,$2,$3,$4,$5,$6)", string(hashPassword), 0, user.Firstname, user.Lastname, user.MiddleName, user.Email)
	if err != nil {
		if errors.As(err, &pgErr) {
			pgErr = err.(pgx.PgError)
			newErr := fmt.Errorf("SQL Error: %s,"+
				" Detail: %s,"+
				" Where: %s,"+
				" Code: %s,"+
				" SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			log.Println(newErr)
		} else {
			log.Println(err)
		}
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	log.Println(sessions)

	c, err := r.Cookie("sessionToken")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value
	log.Printf("sessionToken - %v", sessionToken)
	userSession, exists := sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome, %s!", userSession.Username)))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("sessionToken")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := c.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionToken",
		Value:   "",
		Expires: time.Now(),
	})
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err = pgx.Connect(pgx.ConnConfig{
		Host:     config.DbConfig.DbHostName,
		Port:     uint16(config.DbConfig.DbPort),
		Database: config.DbConfig.DbName,
		User:     config.DbConfig.DbUser,
		Password: config.DbConfig.DbPassword,
	})
	defer conn.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", Welcome)
	r.HandleFunc("/user/login", Login)
	r.HandleFunc("/user/signup", Signup)
	r.HandleFunc("/user/logout", Logout)

	log.Println("Start service")
	log.Fatal(http.ListenAndServe(":8080", r))
}
