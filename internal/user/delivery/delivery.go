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
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserUseCase    usecase.UserUsecase
	SessionUseCase usecase.SessionUsecas
	logger         *logrus.Logger
}

func SetUserRouting(router *mux.Router, log *logrus.Logger, us usecase.UserUsecase) {
	userHandler := &UserHandler{
		UserUseCase: us,
		logger:      log,
	}

	router.HandleFunc("/user/login", userHandler.UserLogin).Methods("POST", "OPTIONS")
}

func (uh *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var creds models.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		uh.logger.Errorf("UserLogin: failed read json with error: %w", err)
		return
	}

	user, status, err := uh.UserUseCase.UserLogin(r.Context(), creds)
	if err != nil || status != http.StatusOK {
		uh.logger.Errorf("UserLogin: failed user verification with [error: %w] [status: %d]", err, status)
		return
	}

	cookie, err := uh.newUserCookie(user.Email)
	if err != nil {
		uh.logger.Errorf("UserDelivery.UserLoginPost: failed create cookie for user with error: %w", err)
		return
	}

	sess := models.Session{
		Cookie: cookie.Value,
		Id:     user.Id,
	}

	err = uh.SessionUseCase.AddSession(r.Context(), sess)
	if err != nil {
		uh.logger.Errorf("UserLoginPost: failed add session in tnt for user with error: %w", err)
		return
	}

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
		SameSite: http.SameSiteNoneMode,
	}

	return cookie, nil
}
