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

type MinIOStorage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
	region     string
	log        logger.Logger
}

func NewMinIOStorage(cfg config.Config, log logger.Logger) (*MinIOStorage, error) {
	endpoint := cfg.StorageEndpoint
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	options := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.StorageAccessKey, cfg.StorageSecretKey, ""),
		Secure: cfg.StorageUseSSL,
		Region: cfg.StorageRegion,
	}

	client, err := minio.New(endpoint, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO storage client: %v", err)
	}

	return &MinIOStorage{
		client:     client,
		bucketName: cfg.StorageBucket,
		endpoint:   endpoint,
		useSSL:     options.Secure,
		region:     cfg.StorageRegion,
		log:        log,
	}, nil
}

func (s *MinIOStorage) UploadFile(ctx context.Context, file io.Reader, size int64, contentType string, fileName string) (string, error) {

	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		s.log.Warnf("failed to check bucket existence: %v", err)
	} else if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{Region: s.region})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %v", err)
		}

		policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"]}]}`, s.bucketName)
		_ = s.client.SetBucketPolicy(ctx, s.bucketName, policy)
	}

	info, err := s.client.PutObject(ctx, s.bucketName, fileName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	s.log.Infof("Successfully uploaded %s of size %d to MinIO", fileName, info.Size)

	return s.GetFileURL(ctx, fileName)
}

func (s *MinIOStorage) GetFileURL(ctx context.Context, fileName string) (string, error) {
	protocol := "http"
	if s.useSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, s.endpoint, s.bucketName, fileName), nil
}

func (s *MinIOStorage) DeleteFile(ctx context.Context, fileName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %v", err)
	}
	s.log.Infof("Successfully deleted %s from MinIO", fileName)
	return nil
}
