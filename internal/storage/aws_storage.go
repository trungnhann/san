package storage

import (
	"context"
	"fmt"
	"io"
	"san/internal/config"
	"san/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type AWSStorage struct {
	client     *minio.Client
	bucketName string
	region     string
	log        logger.Logger
}

func NewAWSStorage(cfg config.Config, log logger.Logger) (*AWSStorage, error) {
	var creds *credentials.Credentials
	if cfg.StorageAccessKey != "" && cfg.StorageSecretKey != "" {
		creds = credentials.NewStaticV4(cfg.StorageAccessKey, cfg.StorageSecretKey, "")
	} else {
		creds = credentials.NewIAM("")
	}

	options := &minio.Options{
		Creds:  creds,
		Secure: true,
		Region: cfg.StorageRegion,
	}

	client, err := minio.New("", options)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS storage client: %v", err)
	}

	return &AWSStorage{
		client:     client,
		bucketName: cfg.StorageBucket,
		region:     cfg.StorageRegion,
		log:        log,
	}, nil
}

func (s *AWSStorage) UploadFile(ctx context.Context, file io.Reader, size int64, contentType string, fileName string) (string, error) {
	info, err := s.client.PutObject(ctx, s.bucketName, fileName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to AWS S3: %v", err)
	}

	s.log.Infof("Successfully uploaded %s of size %d to AWS S3", fileName, info.Size)

	// Return the URL
	return s.GetFileURL(ctx, fileName)
}

func (s *AWSStorage) GetFileURL(ctx context.Context, fileName string) (string, error) {
	region := s.region
	if region == "" {
		region = "us-east-1"
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, region, fileName), nil
}

func (s *AWSStorage) DeleteFile(ctx context.Context, fileName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from AWS S3: %v", err)
	}
	s.log.Infof("Successfully deleted %s from AWS S3", fileName)
	return nil
}
