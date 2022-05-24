package delivery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

const maxFileSize = 20 * 1024 * 1025

func SetStudentRouting(router *mux.Router, log *logrus.Logger, su usecase.StudentUsecase, m *middleware.Middleware) {
	studentHandler := &StudentHandler{
		StudentUsecase: su,
		logger:         log,
	}

	router.HandleFunc("/api/v1/student", m.CheckCSRFAndGetStudent(studentHandler.GetStudent)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/table", m.CheckCSRFAndGetStudent(studentHandler.GetTable)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/group", m.CheckCSRFAndGetStudent(studentHandler.GetGroupByUserId)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/event/{id:[0-9]+}/file", m.CheckCSRFAndGetStudent(studentHandler.UploadFile)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/student/event/{id:[0-9]+}/file", m.CheckCSRFAndGetStudent(studentHandler.LoadFile)).Methods("GET", "OPTIONS")
}

func (sh *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	student, ok := r.Context().Value("student").(models.Student)
	if !ok {
		sh.logger.Errorf("Problem with get value from context %v", ok)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//group, studentId, err := sh.StudentUsecase.GetStudentGroup(r.Context(), student.Id)
	//if err != nil {
	//	sh.logger.Errorf("Problem with get student group %s", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//student := models.Student{
	//	Id: studentId,
	//}

	jsnStud, _ := json.Marshal(student)
	w.Write(jsnStud)
}

func (sh *StudentHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	student, ok := r.Context().Value("student").(models.Student)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	table, err := sh.StudentUsecase.GetTable(r.Context(), student.Id)
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

func (sh *StudentHandler) GetGroupByUserId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	student, ok := r.Context().Value("student").(models.Student)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	table, err := sh.StudentUsecase.GetGroup(r.Context(), student.Id)
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

func (sh *StudentHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		sh.logger.Errorf("can't parse file: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	studentEventId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sh.logger.Errorf("can't get studentEventId from url: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uploadedFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		sh.logger.Errorf("can't upload file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	file, err := sh.StudentUsecase.AddFile(r.Context(), uploadedFile, fileHeader.Filename, studentEventId)
	if err != nil {
		sh.logger.Errorf("%s", err)
		return
	}

	fmt.Sprintf("%v", file)

	//responses.SendData(w, photo)
	w.WriteHeader(http.StatusOK)
}

func (sh *StudentHandler) LoadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")

	fileBytes, err := ioutil.ReadFile("/home/zeus/BMSTU-Diploma-project/test.png")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
	return
}
