package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/usecase"
	check "github.com/ZeeeUs/BMSTU-Diploma-project/pkg/checker"
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

	router.HandleFunc("/api/v1/supervisor/courses", m.CheckCSRFAndGetUser(check.Supervisor(supersHandler.GetSupersCourses))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/supervisor", m.CheckCSRFAndGetUser(check.Supervisor(supersHandler.GetSupers))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/course/{id:[0-9]+}/group", m.CheckCSRFAndGetUser(check.Supervisor(supersHandler.GetGroupsByCourseId))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/group/{id:[0-9]+}/students", m.CheckCSRFAndGetUser(check.Supervisor(supersHandler.GetStudentsByGroup))).Methods("GET", "OPTIONS")
}

func (sh *SupersHandler) GetSupers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		return
	}

	superId, err := sh.SupersUseCase.GetSuperId(r.Context(), curUser.Id)
	if err != nil {
		sh.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsnUsr, _ := json.Marshal(superId)
	_, err = w.Write(jsnUsr)
	if err != nil {
		sh.logger.Errorf("GetSupers: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *SupersHandler) GetSupersCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	curUser, ok := r.Context().Value("superId").(models.User)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	courses, err := sh.SupersUseCase.GetSupersCourses(r.Context(), curUser.Id)
	if err != nil {
		sh.logger.Errorf("Problem with get courses")
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

func (sh *SupersHandler) GetGroupsByCourseId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	courseId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sh.logger.Errorf("can't get course id from url: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !curUser.IsSuper {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groups, err := sh.SupersUseCase.GetGroupsByCourseId(r.Context(), courseId)
	if err != nil {
		sh.logger.Errorf("can't get gruops by set id: [id: %d]: %s", courseId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsnGroups, _ := json.Marshal(groups)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsnGroups)
	if err != nil {
		sh.logger.Errorf("GetGroupsByCourseId: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *SupersHandler) GetStudentsByGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	groupId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sh.logger.Errorf("can't get course id from url: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	students, err := sh.SupersUseCase.GetStudentsByGroup(r.Context(), groupId)
	if err != nil {
		sh.logger.Errorf("can't get students by set group id: [id: %d]: %s", groupId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsnStudents, _ := json.Marshal(students)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsnStudents)
	if err != nil {
		sh.logger.Errorf("GetGroupsByCourseId: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
