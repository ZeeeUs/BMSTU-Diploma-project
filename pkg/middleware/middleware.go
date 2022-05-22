package middleware

import (
	"context"

	"net/http"
	"time"

	userRepository "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

type Middleware struct {
	permission Permission
}

func NewMiddleware(ur userRepository.UserRepository, sr userRepository.SessionRepository) *Middleware {
	return &Middleware{
		permission: Permission{
			ur,
			sr,
		},
	}
}

func (m Middleware) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return m.SetCSRF(m.permission.CheckAuth(next))
}

func (m Middleware) CheckCSRFAndAuth(next http.HandlerFunc) http.HandlerFunc {
	return m.CheckCSRF(m.permission.CheckAuth(next))
}

func (m Middleware) GetUser(next http.HandlerFunc) http.HandlerFunc {
	return m.SetCSRF(m.permission.GetCurrentUser(next))
}

func (m Middleware) CheckCSRFAndGetUser(next http.HandlerFunc) http.HandlerFunc {
	return m.CheckCSRF(m.permission.GetCurrentUser(next))
}

func (m Middleware) SetCSRF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			generateCsrfLogic(w)
			next.ServeHTTP(w, r)
		})
}

func (m Middleware) CheckCSRF(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			csrf := r.Header.Get("x-csrf-Token")
			csrfCookie, err := r.Cookie("csrf")

			if err != nil || csrf == "" || csrfCookie.Value == "" || csrfCookie.Value != csrf {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			generateCsrfLogic(w)
			next.ServeHTTP(w, r)
		})
}

type Permission struct {
	ur userRepository.UserRepository
	sr userRepository.SessionRepository
}

func (perm *Permission) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("sessionId")
		if err != nil {
			log.Errorf("Permissions.CheckAuth: no cookie: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := perm.sr.GetSessionByToken(r.Context(), session.Value)
		if err != nil {
			log.Errorf("Permissions.CheckAuth: failed GetSessionByCookie with error: %s", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "userID", id))
		next.ServeHTTP(w, r)
	})
}

func (perm *Permission) GetCurrentUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("sessionId")
		if err != nil {
			log.Errorf("Permissions.CheckAuth: no cookie: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			//w.Write([]byte(""))
			return
		}

		id, err := perm.sr.GetSessionByToken(r.Context(), session.Value)
		if err != nil {
			log.Errorf("Permissions.CheckAuth: failed GetSessionByCookie with error: %s", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		currentUser, err := perm.ur.GetUserById(r.Context(), id)
		if err != nil {
			log.Errorf("Permissions.GetCurrentUser: failed GetUserById with [error: %s]", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", currentUser))
		next.ServeHTTP(w, r)
	})
}

func generateCsrfLogic(w http.ResponseWriter) {
	csrf, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	timeDelta := time.Now().Add(time.Hour * 12)
	csrfCookie := &http.Cookie{Name: "csrf", Value: csrf.String(), Path: "/", HttpOnly: true, Expires: timeDelta}

	http.SetCookie(w, csrfCookie)
	w.Header().Set("csrf", csrf.String())
}
