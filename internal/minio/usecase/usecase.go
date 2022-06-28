package usecase

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/minio/storage"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/sirupsen/logrus"
)

type MinioUseCase interface {
	GetFile(ctx context.Context, fileName string) error
	UploadFile(ctx context.Context, file models.CreateFileUnit, id int) error
	DeleteFile(ctx context.Context, fileName string) error
}

type minioUseCase struct {
	MinioStorage storage.MinioStorage
	logger       *logrus.Logger
}

func NewMinioUseCase(ms storage.MinioStorage, log *logrus.Logger) MinioUseCase {
	return &minioUseCase{
		MinioStorage: ms,
		logger:       log,
	}
}

func (mu *minioUseCase) UploadFile(ctx context.Context, dto models.CreateFileUnit, id int) error {
	dto.NormalizeName()
	file, err := models.NewFile(dto)
	if err != nil {
		return err
	}
	err = mu.MinioStorage.UploadFile(ctx, file, id)
	if err != nil {
		return err
	}
	return nil
}

func (mu *minioUseCase) GetFile(ctx context.Context, fileName string) error {
	return nil
}

func (mu *minioUseCase) DeleteFile(ctx context.Context, fileName string) error {
	return nil
}
