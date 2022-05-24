package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/repository"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type StudentUsecase interface {
	GetStudentGroup(ctx context.Context, id int) (models.Group, int, error)
	GetTable(ctx context.Context, id int) (models.Table, error)
	GetGroup(ctx context.Context, id int) (models.Group, error)
	AddFile(c context.Context, file io.Reader, fileName string, studentEventId int) (models.File, error)
}

type studentUsecase struct {
	StudentRepository repository.StudentRepository
	//FileRepository    repository.FileRepository
	logger         *logrus.Logger
	contextTimeout time.Duration
}

func NewStudentUsecase(sr repository.StudentRepository, log *logrus.Logger) StudentUsecase {
	return &studentUsecase{
		StudentRepository: sr,
		logger:            log,
	}
}

func (su *studentUsecase) GetStudentGroup(ctx context.Context, id int) (models.Group, int, error) {
	group, studentId, err := su.StudentRepository.GetUserGroup(ctx, id)
	if err == pgx.ErrNoRows {
		return models.Group{}, 0, fmt.Errorf("user with id %d is not found", id)
	} else if err != nil {
		return models.Group{}, 0, err
	}

	return group, studentId, nil
}

func (su *studentUsecase) GetTable(ctx context.Context, id int) (models.Table, error) {
	table, err := su.StudentRepository.GetTable(ctx, id)
	if err == pgx.ErrNoRows {
		return models.Table{}, fmt.Errorf("can't get table for student with id = %d: err %s", id, err)
	}
	return table, nil
}

func (su *studentUsecase) GetGroup(ctx context.Context, id int) (models.Group, error) {
	group, err := su.StudentRepository.GetGroup(ctx, id)
	if err == pgx.ErrNoRows {
		return models.Group{}, fmt.Errorf("can't get table for student with id = %d: err %s", id, err)
	}
	return group, nil
}

func (su *studentUsecase) AddFile(c context.Context, file io.Reader, fileName string, studentEventId int) (models.File, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	currStudent, ok := ctx.Value("student").(models.Student)
	if !ok {
		return models.File{}, errors.New("AddPhoto: can't get current student from context")
	}

	err := saveFile(currStudent, file, fileName, su.logger)
	if err != nil {
		return models.File{}, err
	}

	// TODO запись в бд filePath - studentEventId
	//err = h.UserRepo.UpdateImgs(c, currentUser.ID, currentUser.Imgs)
	//if err != nil {
	//	return models.File{}, err
	//}

	//return models.File{File: filePath}, nil
	return models.File{}, nil
}

func saveFile(student models.Student, file io.Reader, fileName string, log *logrus.Logger) error {
	path := "/home/zeus/BMSTU-Diploma-project/" + strconv.Itoa(student.UserId)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Errorf("can't create dir for student with userId %d: %s", student.UserId, err)
			return err
		}
	}

	fileOnDisk, err := os.Create(path + "/" + fileName)
	if err != nil {
		return err
	}
	defer fileOnDisk.Close()

	io.Copy(fileOnDisk, file)

	return nil
}
