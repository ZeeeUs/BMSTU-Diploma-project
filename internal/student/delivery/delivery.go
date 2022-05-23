package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type StudentHandler struct {
	StudentUsecase usecase.StudentUsecase
	logger         *logrus.Logger
}

func SetStudentRouting(router *mux.Router, log *logrus.Logger, su usecase.StudentUsecase, m *middleware.Middleware) {
	studentHandler := &StudentHandler{
		StudentUsecase: su,
		logger:         log,
	}

	router.HandleFunc("/api/v1/student", m.CheckCSRFAndGetUser(studentHandler.GetStudent)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/table", m.CheckCSRFAndGetStudent(studentHandler.GetTable)).Methods("GET", "OPTIONS")
}

func (sh *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	group, err := sh.StudentUsecase.GetStudentGroup(r.Context(), curUser.Id)
	if err != nil {
		sh.logger.Errorf("Problem with get student group %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	student := models.Student{
		Id:         curUser.Id,
		Firstname:  curUser.Firstname,
		MiddleName: curUser.MiddleName,
		Lastname:   curUser.Lastname,
		Email:      curUser.Email,
		Group:      group,
	}

	jsnUsr, _ := json.Marshal(student)
	_, err = w.Write(jsnUsr)
	if err != nil {
		sh.logger.Errorf("GetStudent: failed to write json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *StudentHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	studId, ok := r.Context().Value("studentId").(int)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	table, err := sh.StudentUsecase.GetTable(r.Context(), studId)
	if err != nil {
		sh.logger.Errorf("Can't get table for user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsnTable, err := json.Marshal(table)
	if err != nil {
		sh.logger.Errorf("Can't marshal user table: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsnTable)
}
