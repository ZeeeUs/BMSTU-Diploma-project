package usecase

import (
	"context"
	"fmt"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/repository"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type StudentUsecase interface {
	GetStudentGroup(ctx context.Context, id int) (models.Group, int, error)
	GetTable(ctx context.Context, id int) (models.Table, error)
	GetGroup(ctx context.Context, id int) (models.Group, error)
}

type studentUsecase struct {
	StudentRepository repository.StudentRepository
	logger            *logrus.Logger
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
