package usecase

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/storage"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type StudentUseCase interface {
	GetStudentGroup(ctx context.Context, id int) (models.Group, int, error)
	GetTable(ctx context.Context, id int) (models.Table, error)
	GetGroup(ctx context.Context, id int) (models.Group, error)
	//AddFile(c context.Context, file io.Reader, fileName string, studentEventId int) (models.File, error)
	LoadFile(ctx context.Context, student models.Student, fileName string) ([]byte, error)
	ChangeEventStatus(ctx context.Context, status int, studentEventId int) error
}

type studentUseCase struct {
	StudentStorage storage.StudentStorage
	logger         *logrus.Logger
	contextTimeout time.Duration
}

const sourcePath = "/usr/src/app/upload_files/"

//const sourcePath = "/home/zeus/BMSTU-Diploma-project/"

func NewStudentUseCase(ss storage.StudentStorage, log *logrus.Logger) StudentUseCase {
	return &studentUseCase{
		StudentStorage: ss,
		logger:         log,
	}
}

func (su *studentUseCase) GetStudentGroup(ctx context.Context, id int) (models.Group, int, error) {
	group, studentId, err := su.StudentStorage.GetUserGroup(ctx, id)
	if err == pgx.ErrNoRows {
		return models.Group{}, 0, fmt.Errorf("user with id %d is not found", id)
	} else if err != nil {
		return models.Group{}, 0, err
	}

	return group, studentId, nil
}

func (su *studentUseCase) GetTable(ctx context.Context, id int) (models.Table, error) {
	table, err := su.StudentStorage.GetTable(ctx, id)
	if err == pgx.ErrNoRows {
		return models.Table{}, fmt.Errorf("can't get table for student with id = %d: err %s", id, err)
	}
	return table, nil
}

func (su *studentUseCase) GetGroup(ctx context.Context, id int) (models.Group, error) {
	group, err := su.StudentStorage.GetGroup(ctx, id)
	if err == pgx.ErrNoRows {
		return models.Group{}, fmt.Errorf("can't get table for student with id = %d: err %s", id, err)
	}
	return group, nil
}

//func (su *studentUseCase) AddFile(c context.Context, file io.Reader, fileName string, studentEventId int) (models.File, error) {
//	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
//	defer cancel()
//
//	currStudent, ok := ctx.Value("student").(models.Student)
//	if !ok {
//		return models.File{}, errors.New("AddPhoto: can't get current student from context")
//	}
//
//	savedFileName, err := saveFile(currStudent, file, fileName, su.logger)
//	if err != nil {
//		return models.File{}, err
//	}
//
//	err = su.StudentStorage.AddFileName(ctx, fileName, studentEventId)
//	if err != nil {
//		return models.File{}, err
//	}
//
//	return models.File{File: savedFileName}, nil
//}

func (su *studentUseCase) LoadFile(ctx context.Context, student models.Student, fileName string) ([]byte, error) {
	bytesFile, err := loadFile(student.Id, fileName, su.logger)
	if err != nil {
		return nil, err
	}

	return bytesFile, nil
}

func (su *studentUseCase) ChangeEventStatus(ctx context.Context, status int, studEvent int) error {
	err := su.StudentStorage.ChangeEventStatus(ctx, status, studEvent)
	if err != nil {
		return err
	}

	return nil
}

func saveFile(student models.Student, file io.Reader, fileName string, log *logrus.Logger) (string, error) {
	path := sourcePath + strconv.Itoa(student.Id)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Errorf("can't create dir for student with userId %d: %s", student.Id, err)
			return "", err
		}
	}

	fileOnDisk, err := os.Create(path + "/" + fileName)
	if err != nil {
		return "", err
	}
	defer fileOnDisk.Close()

	io.Copy(fileOnDisk, file)

	return fileName, nil
}

func loadFile(id int, fileName string, log *logrus.Logger) ([]byte, error) {
	path := sourcePath + strconv.Itoa(id)

	fileBytes, err := ioutil.ReadFile(path + "/" + fileName)
	if err != nil {
		log.Errorf("loadFile: %s", err)
		return nil, err
	}
	return fileBytes, nil
}
