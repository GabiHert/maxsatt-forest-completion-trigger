package mock

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Object struct {
	ContentType *string
	Data        []byte
}

type S3Client struct {
	Objects          map[string]map[string]S3Object // bucket -> key -> S3Object
	multipartUploads map[string]*multipartUpload    // uploadId -> upload state
	mu               sync.Mutex
}

type multipartUpload struct {
	parts  map[int32][]byte
	bucket string
	key    string
}

var s3Instance *S3Client
var s3Init sync.Once

func NewS3Client() *S3Client {
	s3Init.Do(func() {
		s3Instance = &S3Client{
			Objects:          make(map[string]map[string]S3Object),
			multipartUploads: make(map[string]*multipartUpload),
		}
	})
	return s3Instance
}

func (m *S3Client) PutObject(_ context.Context, params *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	bucket := *params.Bucket
	key := *params.Key

	if _, ok := m.Objects[bucket]; !ok {
		m.Objects[bucket] = make(map[string]S3Object)
	}

	var buf bytes.Buffer
	if params.Body != nil {
		_, _ = io.Copy(&buf, params.Body)
	}

	m.Objects[bucket][key] = S3Object{
		Data:        buf.Bytes(),
		ContentType: params.ContentType,
	}

	return &s3.PutObjectOutput{}, nil
}

func (m *S3Client) GetObject(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	bucket := *params.Bucket
	key := *params.Key

	if obj, ok := m.Objects[bucket][key]; ok {
		contentLength := int64(len(obj.Data))

		return &s3.GetObjectOutput{
			Body:          io.NopCloser(bytes.NewReader(obj.Data)),
			ContentLength: &contentLength,
			ContentType:   obj.ContentType,
		}, nil
	}

	return nil, errors.New("object not found")
}

func (m *S3Client) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Objects = make(map[string]map[string]S3Object)
	m.multipartUploads = make(map[string]*multipartUpload)
}

// CreateMultipartUpload initiates a multipart upload
func (m *S3Client) CreateMultipartUpload(_ context.Context, params *s3.CreateMultipartUploadInput, _ ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	uploadID := aws.String("mock-upload-" + *params.Key)

	m.multipartUploads[*uploadID] = &multipartUpload{
		bucket: *params.Bucket,
		key:    *params.Key,
		parts:  make(map[int32][]byte),
	}

	return &s3.CreateMultipartUploadOutput{
		UploadId: uploadID,
	}, nil
}

// UploadPart uploads a part of a multipart upload
func (m *S3Client) UploadPart(_ context.Context, params *s3.UploadPartInput, _ ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	upload, ok := m.multipartUploads[*params.UploadId]
	if !ok {
		return nil, errors.New("upload not found")
	}

	var buf bytes.Buffer
	if params.Body != nil {
		_, _ = io.Copy(&buf, params.Body)
	}

	upload.parts[*params.PartNumber] = buf.Bytes()

	etag := aws.String("mock-etag")
	return &s3.UploadPartOutput{
		ETag: etag,
	}, nil
}

// CompleteMultipartUpload completes a multipart upload
func (m *S3Client) CompleteMultipartUpload(_ context.Context, params *s3.CompleteMultipartUploadInput, _ ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	upload, ok := m.multipartUploads[*params.UploadId]
	if !ok {
		return nil, errors.New("upload not found")
	}

	// Concatenate all parts in order
	var completeData bytes.Buffer
	for i := int32(1); i <= int32(len(upload.parts)); i++ {
		if part, exists := upload.parts[i]; exists {
			completeData.Write(part)
		}
	}

	// Store the complete object
	bucket := upload.bucket
	key := upload.key

	if _, ok := m.Objects[bucket]; !ok {
		m.Objects[bucket] = make(map[string]S3Object)
	}

	m.Objects[bucket][key] = S3Object{
		Data: completeData.Bytes(),
	}

	// Clean up multipart upload state
	delete(m.multipartUploads, *params.UploadId)

	return &s3.CompleteMultipartUploadOutput{}, nil
}

// AbortMultipartUpload aborts a multipart upload
func (m *S3Client) AbortMultipartUpload(_ context.Context, params *s3.AbortMultipartUploadInput, _ ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.multipartUploads, *params.UploadId)

	return &s3.AbortMultipartUploadOutput{}, nil
}

// HeadObject retrieves metadata about an object
func (m *S3Client) HeadObject(_ context.Context, params *s3.HeadObjectInput, _ ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	bucket := *params.Bucket
	key := *params.Key

	if obj, ok := m.Objects[bucket][key]; ok {
		contentLength := int64(len(obj.Data))

		return &s3.HeadObjectOutput{
			ContentLength: &contentLength,
			ContentType:   obj.ContentType,
		}, nil
	}

	return nil, &types.NoSuchKey{
		Message: aws.String("object not found"),
	}
}
