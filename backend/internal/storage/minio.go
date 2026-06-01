package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"billsoftware/backend/internal/config"
)

type ObjectStorage struct {
	client         *minio.Client
	bucket         string
	presignExpires time.Duration
}

func NewObjectStorage(cfg config.MinIOConfig) (*ObjectStorage, error) {
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("minio endpoint is required")
	}

	bucket := strings.TrimSpace(cfg.Bucket)
	if bucket == "" {
		return nil, fmt.Errorf("minio bucket is required")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(strings.TrimSpace(cfg.AccessKey), strings.TrimSpace(cfg.SecretKey), ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	presignMinutes := cfg.PresignMinutes
	if presignMinutes <= 0 {
		presignMinutes = 5
	}

	return &ObjectStorage{
		client:         client,
		bucket:         bucket,
		presignExpires: time.Duration(presignMinutes) * time.Minute,
	}, nil
}

func (s *ObjectStorage) Put(ctx context.Context, objectKey string, data []byte, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, objectKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (s *ObjectStorage) Delete(ctx context.Context, objectKeys ...string) error {
	for _, objectKey := range objectKeys {
		objectKey = strings.TrimSpace(objectKey)
		if objectKey == "" {
			continue
		}
		if err := s.client.RemoveObject(ctx, s.bucket, objectKey, minio.RemoveObjectOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (s *ObjectStorage) PresignedGetURL(ctx context.Context, objectKey string) (string, error) {
	objectKey = strings.TrimSpace(objectKey)
	if objectKey == "" {
		return "", nil
	}
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, objectKey, s.presignExpires, url.Values{})
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
