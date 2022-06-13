package storage

import (
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type MinioStorage interface {
}

type minioClient struct {
	client *minio.Client
	logger *logrus.Logger
}

func NewMinioStorage(endpoint string, accessKeyId string, secretKeyId string, logger *logrus.Logger) MinioStorage {
	//client, err := minio.New(endpoint, accessKeyId, secretKeyId, logger)
	return nil
}
