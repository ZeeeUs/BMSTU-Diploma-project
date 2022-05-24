package delivery

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/hasher"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

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

	router.HandleFunc("/api/v1/user/login", m.SetCSRF(userHandler.UserLogin)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/user/login", m.CheckCSRFAndGetUser(userHandler.UpdateUser)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/user", m.CheckCSRFAndGetUser(userHandler.GetUser)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/user/logout", m.CheckCSRFAndGetUser(userHandler.UserLogout)).Methods("GET", "OPTIONS")
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

	user, err := uh.UserUseCase.UserLogin(r.Context(), creds)
	if err != nil {
		uh.logger.Errorf("UserLogin: failed user verification with [error: %s]", err)
		w.WriteHeader(http.StatusUnauthorized)
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
		uh.logger.Errorf("UserLoginPost: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		uh.logger.Errorf("Problem with get value from cookie %v", ok)
		return
	}

	jsnUsr, _ := json.Marshal(curUser)
	_, err := w.Write(jsnUsr)
	if err != nil {
		uh.logger.Errorf("GetUser: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var (
		updateData models.UpdateUser
		//validOldPass bool
	)
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		uh.logger.Errorf("UpdateUser: failed read json with error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		uh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//uh.logger.Info(curUser.PassStatus, curUser.Id)

	uh.logger.Info(curUser.Email, "", curUser.Id)

	if curUser.PassStatus {
		_, err = hasher.ComparePasswords(curUser.Password, updateData.OldPass)
		if err != nil {
			uh.logger.Errorf("compare password in UserUpdate return error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		if strings.Compare(curUser.Password, updateData.OldPass) != 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	//if !validOldPass {
	//	w.Write([]byte("passwords doesn't matched"))
	//	return
	//}

	newPass, err := hasher.HashAndSalt(updateData.NewPass)
	if err != nil {
		uh.logger.Errorf("UpdateUser: password don't be hashed and salt: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := uh.UserUseCase.UpdateUser(r.Context(), newPass, curUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(strconv.Itoa(id)))
}

func (uh *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//err := uh.SessionUseCase.DeleteSession(r.Context())
	//if err != nil {
	//	uh.logger.Errorf("UserDelivery.UserLogout: failed delete session in redis for user with error: %s", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	authCookie := &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, authCookie)

	csrfCookie := &http.Cookie{
		Name:     "csrf",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, csrfCookie)

	w.WriteHeader(http.StatusOK)
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
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}

	return cookie, nil
}
