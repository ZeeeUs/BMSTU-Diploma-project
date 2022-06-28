package usecase

import (
	"context"
	"fmt"
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
	ChangeEventStatus(ctx context.Context, status int, studentEventId int) error
}

type studentUseCase struct {
	StudentStorage storage.StudentStorage
	logger         *logrus.Logger
	contextTimeout time.Duration
}

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

func (su *studentUseCase) ChangeEventStatus(ctx context.Context, status int, studEvent int) error {
	err := su.StudentStorage.ChangeEventStatus(ctx, status, studEvent)
	if err != nil {
		return err
	}

	return nil
}
