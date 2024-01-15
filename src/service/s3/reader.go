package s3

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

// Reader reads contents from s3 buckets
type Reader struct {
	timeout time.Duration
	client  *s3.Client
}

// NewReader creates a new S3 reader
func NewReader(timeout time.Duration) (*Reader, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "could not get AWS config")
	}

	s3Client := s3.NewFromConfig(cfg)
	return &Reader{
		timeout: timeout,
		client:  s3Client,
	}, nil
}

func (r *Reader) getCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.timeout)
}

// GetObjectContents fetches the contents of s3://<bucket>/<path>
func (r *Reader) GetObjectContents(bucket, path string) ([]byte, error) {
	ctx, cancel := r.getCtx()
	defer cancel()

	response, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not download s3 object")
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read object body")
	}

	return contents, nil
}
