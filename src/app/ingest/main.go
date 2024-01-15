package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dleonard203/go-live-mocking/src/domain"
	"github.com/dleonard203/go-live-mocking/src/repository/database"
	s3wrapper "github.com/dleonard203/go-live-mocking/src/service/s3"
	"github.com/pkg/errors"
)

// S3SensorProcessor reads sensor readings from the s3 bucket and writes them to the database
type S3SensorProcessor struct {
	s3Reader domain.IS3Reader
	dbWriter domain.ISensorIngetsor
}

// NewS3SensorProcessor creates a new S3SensorProcessor
func NewS3SensorProcessor(s3Reader domain.IS3Reader, dbWriter domain.ISensorIngetsor) *S3SensorProcessor {
	return &S3SensorProcessor{s3Reader: s3Reader, dbWriter: dbWriter}
}

func getProcessor() (*S3SensorProcessor, error) {
	s3Client, err := s3wrapper.NewReader(time.Second * 10)
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate s3 client")
	}

	dbClient, err := database.NewSensorWriter(os.Getenv("POSTGRES_DSN"), time.Second*10)
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate database client")
	}

	return NewS3SensorProcessor(s3Client, dbClient), nil
}

func (s *S3SensorProcessor) performIngest(ctx context.Context, s3Event events.S3Event) error {

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		path := record.S3.Object.Key
		fileContents, err := s.s3Reader.GetObjectContents(bucket, path)
		if err != nil {
			return fmt.Errorf(
				"could not fetch s3 file contents for bucket=%s path=%s: %s",
				bucket,
				path,
				err.Error(),
			)
		}
		readings := make([]domain.HumidityStats, 0)
		decoder := json.NewDecoder(strings.NewReader(string(fileContents)))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&readings)
		if err != nil {
			return fmt.Errorf(
				"could not unmarshal contents of bucket=%s path=%s into []HumidityStats: %s",
				bucket,
				path, err.Error(),
			)
		}

		err = s.dbWriter.WriteHumidityReadings(readings)
		if err != nil {
			return fmt.Errorf(
				"could not write records from bucket=%s path=%s to the database: %s",
				bucket,
				path,
				err.Error(),
			)
		}

	}

	return nil
}

func main() {
	processor, err := getProcessor()
	if err != nil {
		panic(
			fmt.Errorf("could not create processor: %s",
				err.Error()),
		)
	}

	lambda.Start(processor.performIngest)
}
