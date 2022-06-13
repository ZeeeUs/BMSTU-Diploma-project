package storage

import (
	"context"
	"fmt"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/minio"

	"github.com/sirupsen/logrus"
)

type MinioStorage interface {
	GetFile(ctx context.Context, fileName string) error
	UploadFile(ctx context.Context, file models.FileUnit, id int) error
	DeleteFile(ctx context.Context, fileName string) error
}

type minioStorage struct {
	client *minio.Client
	logger *logrus.Logger
}

func NewMinioStorage(endpoint string, accessKeyId string, secretKeyId string, logger *logrus.Logger) (MinioStorage, error) {
	client, err := minio.NewClient(endpoint, accessKeyId, secretKeyId, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}

	return &minioStorage{
		client: client,
		logger: logger,
	}, nil
}

func (ms *minioStorage) GetFile(ctx context.Context, fileName string) error {
	return nil
}

func (ms *minioStorage) UploadFile(ctx context.Context, file models.FileUnit, id int) error {
	err := ms.client.UploadFile(ctx, "123567.png", "test.txt", "bmstu", file.PayloadSize, file.Payload)
	if err != nil {
		return err
	}
	return nil
}

func (ms *minioStorage) DeleteFile(ctx context.Context, fileName string) error {
	return nil
}
