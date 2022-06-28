package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/minio"

	"github.com/sirupsen/logrus"
)

//const (
//	bucket = "bmstu"
//)

var bucket = "bmstu"

type MinioStorage interface {
	GetFile(ctx context.Context, fileName string) error
	UploadFile(ctx context.Context, file *models.File, id int) error
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

func (ms *minioStorage) UploadFile(ctx context.Context, file *models.File, id int) error {
	err := ms.client.UploadFile(ctx, file.Name, file.Name, bucket, file.Size, bytes.NewBuffer(file.Bytes))
	if err != nil {
		return err
	}
	return nil
}

func (ms *minioStorage) DeleteFile(ctx context.Context, fileName string) error {
	return nil
}
