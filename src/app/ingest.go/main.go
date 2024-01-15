package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

// HumidityStats contains information from a humidity sensor reading
type HumidityStats struct {
	Reading    int    `json:"reading"`
	Unit       string `json:"unit"`
	CustomerID int    `json:"customer_id"`
	SensorID   int    `json:"sensor_id"`
}

func getFileContents(event events.S3EventRecord) ([]byte, error) {
	bucket := event.S3.Bucket.Name
	path := event.S3.Object.Key

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "could not get AWS config")
	}

	s3Client := s3.NewFromConfig(cfg)
	response, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
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

func writeRecords(readings []HumidityStats) error {
	// db connection code, db insert code...
	return nil
}

func performIngest(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		path := record.S3.Object.Key
		fileContents, err := getFileContents(record)
		if err != nil {
			return fmt.Errorf(
				"could not fetch s3 file contents for bucket=%s path=%s: %s",
				bucket,
				path,
				err.Error(),
			)
		}
		readings := make([]HumidityStats, 0)
		err = json.Unmarshal(fileContents, &readings)
		if err != nil {
			return fmt.Errorf(
				"could not unmarshal contents of bucket=%s path=%s into []HumidityStats: %s",
				bucket,
				path, err.Error(),
			)
		}

		err = writeRecords(readings)
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
	lambda.Start(performIngest)
}
