package usecase

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/storage"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type SupersUsecase interface {
	GetSupersCourses(ctx context.Context, id int) ([]models.Course, error)
	GetSuperId(ctx context.Context, id int) (int, error)
	GetGroupsByCourseId(ctx context.Context, courseId int) ([]models.GroupByCourse, error)
	GetStudentsByGroup(ctx context.Context, courseId int) ([]models.StudentByGroup, error)
	GetEventsByCourseId(ctx context.Context, groupId int) ([]models.Event, error)
	GetStudentEvents(ctx context.Context, studentId int, courseId int) ([]models.StudentEvent, error)
}

type supersUsecase struct {
	SupersRepository storage.SupersStorage
	logger           *logrus.Logger
}

func NewSupersUsecase(sr storage.SupersStorage, log *logrus.Logger) SupersUsecase {
	return &supersUsecase{
		SupersRepository: sr,
		logger:           log,
	}
}

func (su *supersUsecase) GetSupersCourses(ctx context.Context, id int) ([]models.Course, error) {
	supervisor, err := su.SupersRepository.GetSupersCourses(ctx, id)
	if err == pgx.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return supervisor, nil
}

func (su *supersUsecase) GetSuperId(ctx context.Context, id int) (int, error) {
	superId, err := su.SupersRepository.GetSuperId(ctx, id)
	if err == pgx.ErrNoRows {
		return 0, err
	} else if err != nil {
		return 0, err
	}

	return superId, nil
}

func (su *supersUsecase) GetGroupsByCourseId(ctx context.Context, courseId int) ([]models.GroupByCourse, error) {
	groups, err := su.SupersRepository.GetGroupsByCourseId(ctx, courseId)
	if err == pgx.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return groups, nil
}

func (su *supersUsecase) GetStudentsByGroup(ctx context.Context, groupId int) ([]models.StudentByGroup, error) {
	students, err := su.SupersRepository.GetStudentsByGroup(ctx, groupId)
	if err == pgx.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return students, nil
}

func (su *supersUsecase) GetEventsByCourseId(ctx context.Context, groupId int) ([]models.Event, error) {
	events, err := su.SupersRepository.GetEventsByCourseId(ctx, groupId)
	if err == pgx.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return events, nil
}

func (su *supersUsecase) GetStudentEvents(ctx context.Context, studentId int, courseId int) ([]models.StudentEvent, error) {
	events, err := su.SupersRepository.GetStudentEvents(ctx, studentId, courseId)
	if err == pgx.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	return events, nil
}
