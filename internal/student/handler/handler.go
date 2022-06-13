package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	muc "github.com/ZeeeUs/BMSTU-Diploma-project/internal/minio/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	suc "github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type StudentHandler struct {
	StudentUseCase suc.StudentUseCase
	MinioUseCase   muc.MinioUseCase
	logger         *logrus.Logger
}

const (
	maxFileSize = 20 * 1024 * 1025
	sourcePath  = "/usr/src/app/upload_files/"
)

func SetStudentRouting(router *mux.Router, log *logrus.Logger, su suc.StudentUseCase, mu muc.MinioUseCase, m *middleware.Middleware) {
	studentHandler := &StudentHandler{
		StudentUseCase: su,
		MinioUseCase:   mu,
		logger:         log,
	}

	router.HandleFunc("/api/v1/student", m.CheckCSRFAndGetStudent(studentHandler.GetStudent)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/table", m.CheckCSRFAndGetStudent(studentHandler.GetTable)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/group", m.CheckCSRFAndGetStudent(studentHandler.GetGroupByUserId)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/file/{fileName}", m.CheckCSRFAndGetStudent(studentHandler.LoadFile)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/student/event/{id:[0-9]+}/file", m.CheckCSRFAndGetStudent(studentHandler.UploadFile)).Methods("POST", "OPTIONS")
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

	table, err := sh.StudentUseCase.GetTable(r.Context(), student.Id)
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

	table, err := sh.StudentUseCase.GetGroup(r.Context(), student.Id)
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

	file := models.FileUnit{
		PayloadSize: fileHeader.Size,
		Payload:     uploadedFile,
	}
	defer uploadedFile.Close()

	//fmt.Sprintf("%v%v", file, studentEventId)
	//sh.logger.Info(file.PayloadSize)
	//sh.logger.Info(file, studentEventId)
	err = sh.MinioUseCase.UploadFile(r.Context(), file, studentEventId)
	if err != nil {
		sh.logger.Errorf("can't upload file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Скачивание файла
func (sh *StudentHandler) LoadFile(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/octet-stream")

	nameFromUrl := mux.Vars(r)["fileName"]

	student, ok := r.Context().Value("student").(models.Student)
	if !ok {
		sh.logger.Errorf("Problem with get value from contetxt %v", ok)
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

	err = sh.StudentUseCase.ChangeEventStatus(r.Context(), status.Status, studentEventId)
	if err != nil {
		sh.logger.Errorf("can't change event status on db: %s", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}
