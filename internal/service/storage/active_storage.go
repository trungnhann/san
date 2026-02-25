package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"time"

	dbsqlc "san/internal/db/sqlc"
	"san/internal/storage"
	"san/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type ActiveStorageService struct {
	repo    dbsqlc.Querier
	storage storage.FileStorage
	log     logger.Logger
}

func NewActiveStorageService(repo dbsqlc.Querier, storage storage.FileStorage, log logger.Logger) *ActiveStorageService {
	return &ActiveStorageService{
		repo:    repo,
		storage: storage,
		log:     log,
	}
}

func (s *ActiveStorageService) AttachFile(ctx context.Context, recordType string, recordID string, attachmentName string, file io.Reader, size int64, contentType string, filename string, replaceOld bool) (*dbsqlc.StorageAttachment, error) {
	const maxFileSize = 100 * 1024 * 1024 // 100MB

	if size > maxFileSize {
		return nil, fmt.Errorf("file size exceeds limit of 100MB")
	}

	limitReader := io.LimitReader(file, maxFileSize+1)

	data, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if int64(len(data)) > maxFileSize {
		return nil, fmt.Errorf("file size exceeds limit of 100MB")
	}

	if replaceOld {
		s.log.Infof("Checking for existing attachment to purge: %s for %s/%s", attachmentName, recordType, recordID)
		oldAttachment, err := s.repo.GetAttachmentByRecord(ctx, dbsqlc.GetAttachmentByRecordParams{
			RecordType: recordType,
			RecordID:   recordID,
			Name:       attachmentName,
		})
		if err == nil {
			s.log.Infof("Purging old attachment %s for %s/%s (BlobID: %s)", attachmentName, recordType, recordID, oldAttachment.BlobID)

			oldBlob, blobErr := s.repo.GetBlob(ctx, oldAttachment.BlobID)
			if blobErr == nil {
				if delErr := s.storage.DeleteFile(ctx, oldBlob.Key); delErr != nil {
					s.log.Errorf("Failed to delete old file %s: %v", oldBlob.Key, delErr)
				}
				if delErr := s.repo.DeleteBlob(ctx, oldAttachment.BlobID); delErr != nil {
					s.log.Errorf("Failed to delete old blob record: %v", delErr)
				} else {
					s.log.Infof("Deleted old blob record: %s", oldAttachment.BlobID)
				}
			} else {
				s.log.Warnf("Old blob %s not found for attachment, skipping file/blob deletion: %v", oldAttachment.BlobID, blobErr)
			}

			if delErr := s.repo.DeleteAttachment(ctx, oldAttachment.ID); delErr != nil {
				s.log.Errorf("Failed to delete old attachment record: %v", delErr)
				return nil, delErr
			} else {
				s.log.Infof("Deleted old attachment record: %s", oldAttachment.ID)
			}
		} else {
			s.log.Infof("No existing attachment found to purge (or error): %v", err)
		}
	}

	hash := md5.Sum(data)
	checksum := base64.StdEncoding.EncodeToString(hash[:])
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("%s/%s/%s/%d%s", recordType, recordID, attachmentName, time.Now().UnixNano(), ext)

	url, err := s.storage.UploadFile(ctx, bytes.NewReader(data), size, contentType, key)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	s.log.Infof("File uploaded to %s with checksum %s", url, checksum)

	blobID := uuid.New()
	blob, err := s.repo.CreateBlob(ctx, dbsqlc.CreateBlobParams{
		ID:          blobID,
		Key:         key,
		Filename:    filename,
		ContentType: &contentType,
		ByteSize:    size,
		Checksum:    &checksum,
		Metadata:    pgtype.JSONB{Bytes: []byte("{}"), Status: pgtype.Present},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create blob record: %w", err)
	}

	attachment, err := s.repo.CreateAttachment(ctx, dbsqlc.CreateAttachmentParams{
		Name:       attachmentName,
		RecordType: recordType,
		RecordID:   recordID,
		BlobID:     blob.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create attachment record: %w", err)
	}

	return attachment, nil
}

func (s *ActiveStorageService) GetAttachmentURL(ctx context.Context, recordType string, recordID string, attachmentName string) (string, error) {
	row, err := s.repo.GetAttachment(ctx, dbsqlc.GetAttachmentParams{
		RecordType: recordType,
		RecordID:   recordID,
		Name:       attachmentName,
	})
	if err != nil {
		return "", err
	}

	return s.storage.GetFileURL(ctx, row.Key)
}
