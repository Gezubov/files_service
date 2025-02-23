package service

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	client *minio.Client
	bucket string
}

func NewMinioService(endpoint, accessKeyID, secretAccessKey, bucket string) (*MinioService, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinioService{client: minioClient, bucket: bucket}, nil
}

func (m *MinioService) UploadFile(filePath string, objectName string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Upload the file
	_, err = m.client.PutObject(context.Background(), m.bucket, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
