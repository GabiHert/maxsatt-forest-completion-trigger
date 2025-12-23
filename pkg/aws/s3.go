package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3HelperAdapter interface {
	PutObject(ctx context.Context, bucketName, path, fileName string, data []byte) (*string, error)
	GetObject(ctx context.Context, bucket, path string) (*Object, error)
}

type S3 interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type s3Helper struct {
	logger loggerAdapter
	s3     S3
}

type Object struct {
	Reader      io.ReadCloser
	ContentType string
	Filename    string
	Size        int64
}

func S3Client(region string) *s3.Client {
	cfg := getConfig(region)
	if awsUrl := os.Getenv("AWS_URL"); awsUrl != "" {
		return s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = &awsUrl
		})
	}

	return s3.NewFromConfig(cfg)
}

func S3Helper(s3 S3, logger loggerAdapter) S3HelperAdapter {
	return &s3Helper{
		logger: logger,
		s3:     s3,
	}
}

func (s *s3Helper) PutObject(ctx context.Context, bucketName, path, fileName string, data []byte) (*string, error) {
	s.logger.Debug(ctx, "Started", map[string]any{"bucketName": bucketName, "path": path})

	if bucketName == "" {
		return nil, fmt.Errorf("bucket name is empty")
	}

	file := bytes.NewReader(data)
	s3Path := filepath.Join(strings.TrimPrefix(path, "/"), fileName)

	_, err := s.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &s3Path,
		Body:   file,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3Path)

	s.logger.Debug(ctx, "Finished", map[string]any{"bucketName": bucketName, "path": path})
	return &url, nil
}

func (s *s3Helper) GetObject(ctx context.Context, bucket, path string) (*Object, error) {
	s.logger.Debug(ctx, "Started", map[string]any{"bucketName": bucket, "path": path})

	output, err := s.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &path,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	s.logger.Debug(ctx, "Finished", map[string]any{"bucketName": bucket, "path": path})
	return &Object{
		Reader:      output.Body,
		ContentType: aws.ToString(output.ContentType),
		Size:        *output.ContentLength,
		Filename:    filepath.Base(path),
	}, nil
}
