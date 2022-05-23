package usecase

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/repository"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type SupersUsecase interface {
	GetSupersCourses(ctx context.Context, id int) ([]models.Course, error)
	GetSuperId(ctx context.Context, id int) (int, error)
}

type supersUsecase struct {
	SupersRepository repository.SupersRepository
	logger           *logrus.Logger
}

func NewSupersUsecase(sr repository.SupersRepository, log *logrus.Logger) SupersUsecase {
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
