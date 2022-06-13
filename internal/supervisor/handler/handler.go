package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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

const sourcePath = "/usr/src/app/upload_files/"

func SetSupersRouting(router *mux.Router, log *logrus.Logger, su usecase.SupersUsecase, m *middleware.Middleware) {
	supersHandler := &SupersHandler{
		SupersUseCase: su,
		logger:        log,
	}

	router.HandleFunc("/api/v1/supervisor/courses", m.CheckCSRFAndGetSuper(supersHandler.GetSupersCourses)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/supervisor", m.CheckCSRFAndGetSuper(supersHandler.GetSupers)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/course/{id:[0-9]+}/group", m.CheckCSRFAndGetSuper(supersHandler.GetGroupsByCourseId)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/group/{id:[0-9]+}/students", m.CheckCSRFAndGetSuper(supersHandler.GetStudentsByGroup)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/course/{id:[0-9]+}/events", m.CheckCSRFAndGetSuper(supersHandler.GetEventsByCourse)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/supervisor/student/{id:[0-9]+}/file/{fileName}", m.CheckCSRFAndGetSuper(supersHandler.DownloadFile)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/supervisor/student/{studentId}/course/{courseId}", m.CheckCSRFAndGetSuper(supersHandler.GetStudentEventsByCourse)).Methods("GET", "OPTIONS")
}

func (sh *SupersHandler) GetSupers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	super, ok := r.Context().Value("supervisor").(models.Supervisor)
	if !ok {
		sh.logger.Errorf("Problem with get value from context %v", ok)
		return
	}

	jsnUsr, _ := json.Marshal(super)
	_, err := w.Write(jsnUsr)
	if err != nil {
		sh.logger.Errorf("GetSupers: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *SupersHandler) GetSupersCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	super, ok := r.Context().Value("supervisor").(models.Supervisor)
	if !ok {
		sh.logger.Errorf("Problem with get value from context %v", ok)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	courses, err := sh.SupersUseCase.GetSupersCourses(r.Context(), super.UserId)
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

func (sh *SupersHandler) GetEventsByCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	courseId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sh.logger.Errorf("can't get course id from url: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events, err := sh.SupersUseCase.GetEventsByCourseId(r.Context(), courseId)
	if err != nil {
		sh.logger.Errorf("can't get gruops by set id: [id: %d]: %s", courseId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsnGroups, _ := json.Marshal(events)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsnGroups)
	if err != nil {
		sh.logger.Errorf("GetGroupsByCourseId: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *SupersHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {

	nameFromUrl := mux.Vars(r)["fileName"]
	studId := mux.Vars(r)["id"]
	if nameFromUrl == "" || studId == "" {
		sh.logger.Infof("can't get vals from url: fileName - %s, id - %s", nameFromUrl, studId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := sourcePath + studId + "/" + nameFromUrl
	http.ServeFile(w, r, path)
	return
}

func (sh *SupersHandler) GetStudentEventsByCourse(w http.ResponseWriter, r *http.Request) {
	studentId, err := strconv.Atoi(mux.Vars(r)["studentId"])
	if err != nil {
		sh.logger.Errorf("can't get student id from url: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	courseId, err := strconv.Atoi(mux.Vars(r)["courseId"])
	if err != nil {
		sh.logger.Errorf("can't get course id from url: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events, err := sh.SupersUseCase.GetStudentEvents(r.Context(), studentId, courseId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsnEvents, _ := json.Marshal(events)
	w.Write(jsnEvents)

}
