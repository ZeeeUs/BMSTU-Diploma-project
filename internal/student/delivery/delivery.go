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

func SetStudentRouting(router *mux.Router, log *logrus.Logger, m *middleware.Middleware) {
	studentHandler := &StudentHandler{
		logger: log,
	}

	router.HandleFunc("/api/v1/student", m.CheckCSRFAndGetUser(studentHandler.GetStudent))
}

func (sh *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	curUser, ok := r.Context().Value("user").(models.User)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		return
	}

	group, err := sh.StudentUsecase.GetStudentGroup(r.Context(), curUser.Id)

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
