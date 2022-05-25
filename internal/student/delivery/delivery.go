package delivery

import (
	"encoding/json"
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

const sourcePath = "/usr/src/app/upload_files/"

//const sourcePath = "/home/zeus/BMSTU-Diploma-project/"

func SetStudentRouting(router *mux.Router, log *logrus.Logger, su usecase.StudentUsecase, m *middleware.Middleware) {
	studentHandler := &StudentHandler{
		StudentUsecase: su,
		logger:         log,
	}

	router.HandleFunc("/api/v1/student", m.CheckCSRFAndGetStudent(studentHandler.GetStudent)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/table", m.CheckCSRFAndGetStudent(studentHandler.GetTable)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/group", m.CheckCSRFAndGetStudent(studentHandler.GetGroupByUserId)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/event/{id:[0-9]+}/file", m.CheckCSRFAndGetStudent(studentHandler.UploadFile)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/student/file/{fileName}", m.CheckCSRFAndGetStudent(studentHandler.LoadFile)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/event/{id:[0-9]+}/status", m.CheckCSRFAndGetStudent(studentHandler.ChangeEventStatus)).Methods("PUT", "OPTIONS")
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsnFile, _ := json.Marshal(file)
	w.WriteHeader(http.StatusOK)
	w.Write(jsnFile)
}

// Скачивание файла
func (sh *StudentHandler) LoadFile(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/octet-stream")

	nameFromUrl := mux.Vars(r)["fileName"]

	student, ok := r.Context().Value("student").(models.Student)
	if !ok {
		sh.logger.Errorf("Problem with get value from cookie %v", ok)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	path := sourcePath + strconv.Itoa(student.Id) + "/" + nameFromUrl
	http.ServeFile(w, r, path)
	return
}

func (sh *StudentHandler) ChangeEventStatus(w http.ResponseWriter, r *http.Request) {
	var status models.EventStatus

	studentEventId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sh.logger.Errorf("can't get studentEventId from url: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		sh.logger.Errorf("can't decode status from request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = sh.StudentUsecase.ChangeEventStatus(r.Context(), status.Status, studentEventId)
	if err != nil {
		sh.logger.Errorf("can't change event status on db: %s", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}
