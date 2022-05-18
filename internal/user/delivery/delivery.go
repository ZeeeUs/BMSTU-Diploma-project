package delivery

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/hasher"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type Manager interface {
}

type UserHandler struct {
	UserUseCase    usecase.UserUsecase
	SessionUseCase usecase.SessionUsecase
	logger         *logrus.Logger
}

func SetUserRouting(router *mux.Router, log *logrus.Logger, uu usecase.UserUsecase, su usecase.SessionUsecase, m *middleware.Middleware) {
	userHandler := &UserHandler{
		UserUseCase:    uu,
		SessionUseCase: su,
		logger:         log,
	}

	router.HandleFunc("/user/login", m.SetCSRF(userHandler.UserLogin)).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/login", m.CheckCSRFAndGetUser(userHandler.UpdatePassword)).Methods("PUT", "OPTIONS")
}

func (uh *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var creds models.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		uh.logger.Errorf("UserLogin: failed read json with error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, status, err := uh.UserUseCase.UserLogin(r.Context(), creds)
	if err != nil || status != http.StatusOK {
		uh.logger.Errorf("UserLogin: failed user verification with [error: %s] [status: %d]", err, status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := uh.newUserCookie(user.Email)
	if err != nil {
		uh.logger.Errorf("UserDelivery.UserLoginPost: failed create cookie for user with error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sess := models.Session{
		Cookie: cookie.Value,
		Id:     user.Id,
	}

	err = uh.SessionUseCase.AddSession(r.Context(), sess)
	if err != nil {
		uh.logger.Errorf("UserLoginPost: failed add session too redis for user with error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	res, _ := json.Marshal(user)
	_, err = w.Write(res)
	if err != nil {
		uh.logger.Errorf("UserLoginPost: faild to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// TODO реализовать ручку
	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		uh.logger.Errorf("Problem with get value from cookie %v", ok)
	}

	uh.logger.Infof("Current user %v", curUser)

	jsnUsr, _ := json.Marshal(curUser)
	w.Write(jsnUsr)
}

func (uh *UserHandler) newUserCookie(email string) (http.Cookie, error) {
	expiration := time.Now().Add(12 * time.Hour)

	hashedEmail, err := hasher.HashAndSalt(email)
	if err != nil {
		return http.Cookie{}, err
	}
	data := hashedEmail + time.Now().String()
	md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    md5CookieValue,
		Expires:  expiration,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	}

	return cookie, nil
}
