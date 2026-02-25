package storage

import (
	"fmt"
	"san/internal/config"
	"san/pkg/logger"
)

type StorageType string

const (
	StorageTypeS3    StorageType = "s3"
	StorageTypeMinIO StorageType = "minio"
)

func NewStorage(cfg config.Config, log logger.Logger) (FileStorage, error) {
	var storageType StorageType

	if cfg.Environment == "production" || cfg.Environment == "staging" {
		storageType = StorageTypeS3
	} else {
		storageType = StorageTypeMinIO
	}

	switch storageType {
	case StorageTypeS3:
		return NewAWSStorage(cfg, log)
	case StorageTypeMinIO:
		return NewMinIOStorage(cfg, log)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
