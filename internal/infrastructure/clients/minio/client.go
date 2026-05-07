package minio

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"lms_system/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client         *minio.Client
	bucketName     string
	publicEndpoint string // просто строка
}

func NewClient(cfg *config.Config) (*Client, error) {
	minioClient, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKeyID, cfg.MinIO.SecretAccessKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.MinIO.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.MinIO.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &Client{
		client:         minioClient,
		bucketName:     cfg.MinIO.BucketName,
		publicEndpoint: cfg.MinIO.PublicEndpoint,
	}, nil
}

// UploadFile uploads a file to MinIO
func (c *Client) UploadFile(ctx context.Context, objectName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Get content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload the file
	info, err := c.client.PutObject(ctx, c.bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return info.Key, nil
}

// DownloadFile downloads a file from MinIO
func (c *Client) DownloadFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := c.client.GetObject(ctx, c.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return object, nil
}

// DeleteFile deletes a file from MinIO
func (c *Client) DeleteFile(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetFileURL returns a presigned URL for downloading a file

func (c *Client) GetFileURL(ctx context.Context, objectName string, expiry int) (string, error) {
	// Просто возвращаем прямой публичный URL
	return fmt.Sprintf("http://%s/%s/%s", c.publicEndpoint, c.bucketName, objectName), nil
}
