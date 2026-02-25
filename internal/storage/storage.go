package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	UploadFile(ctx context.Context, file io.Reader, size int64, contentType string, fileName string) (string, error)
	GetFileURL(ctx context.Context, fileName string) (string, error)
	DeleteFile(ctx context.Context, fileName string) error
}
