package s3

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dleonard203/go-live-mocking/src/awsutils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Writer implements the domain.IS3Writer interface
type Writer struct {
	bucket  string
	timeout time.Duration
	client  *s3.Client
}

// NewWriter creates a new Writer for AWS S3
func NewWriter(bucket, region string, timeout time.Duration) (*Writer, error) {
	if timeout <= 0 {
		return nil, errors.New("timeout must be positive duration")
	}

	if !awsutils.IsValidRegion(region) {
		return nil, fmt.Errorf("region %s is not a supported region", region)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return &Writer{
		client:  client,
		bucket:  bucket,
		timeout: timeout,
	}, nil
}

func (w *Writer) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), w.timeout)
}

// PutObject writes the contents to the S3 bucket at the objectPath
func (w *Writer) PutObject(contents, objectPath string) error {
	ctx, cancel := w.getContext()
	defer cancel()

	reader := strings.NewReader(contents)
	_, err := w.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(w.bucket),
		Key:    aws.String(objectPath),
		Body:   reader,
	})

	return err
}
