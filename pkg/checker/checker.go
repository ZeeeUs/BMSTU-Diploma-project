package check

import (
	"net/http"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"

	log "github.com/sirupsen/logrus"
)

func Supervisor(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		curUser, ok := r.Context().Value("user").(models.User)
		if !ok {
			log.Errorf("check supervisor: no ctx user")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if curUser.IsSuper {
			next.ServeHTTP(w, r)
			return
		}

		log.Errorf("check supervisor: user isn't admin")
		w.WriteHeader(http.StatusForbidden)
		return
	})
}
