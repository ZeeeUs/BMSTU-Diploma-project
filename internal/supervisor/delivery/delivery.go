package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type SupersHandler struct {
	SupersUseCase usecase.SupersUsecase
	logger        *logrus.Logger
}

func SetSupersRouting(router *mux.Router, log *logrus.Logger, su usecase.SupersUsecase, m *middleware.Middleware) {
	supersHandler := &SupersHandler{
		SupersUseCase: su,
		logger:        log,
	}

	router.HandleFunc("/api/v1/supervisor/courses", m.CheckCSRFAndGetUser(supersHandler.GetSupersCourses)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/supervisor", m.CheckCSRFAndGetUser(supersHandler.GetSupers)).Methods("GET", "OPTIONS")
}

func (sh *SupersHandler) GetSupers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		return
	}

	jsnUsr, _ := json.Marshal(curUser)
	_, err := w.Write(jsnUsr)
	if err != nil {
		sh.logger.Errorf("GetSupers: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *SupersHandler) GetSupersCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !curUser.IsSuper {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	courses, err := sh.SupersUseCase.GetSupersCourses(r.Context(), curUser.Id)
	if err != nil {
		sh.logger.Errorf("Problem with get supervisor")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsnCourses, _ := json.Marshal(courses)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsnCourses)
	if err != nil {
		sh.logger.Errorf("GetSupersCourses: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
