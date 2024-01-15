package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dleonard203/go-live-mocking/src/domain"
	"github.com/stretchr/testify/require"
)

type s3Mock struct {
	getObjectBytes []byte
	getObjectErr   error
}

// GetObjectContents mock
func (s *s3Mock) GetObjectContents(bucket, path string) ([]byte, error) {
	return s.getObjectBytes, s.getObjectErr
}

type dbMock struct {
	writeHumdityErr error
}

// WriteHumidityReadings mock
func (d *dbMock) WriteHumidityReadings(readings []domain.HumidityStats) error {
	return d.writeHumdityErr
}

func getMockS3Event(numRecords int) events.S3Event {
	toReturn := events.S3Event{}
	for index := 0; index < numRecords; index++ {
		toReturn.Records = append(toReturn.Records, events.S3EventRecord{
			S3: events.S3Entity{
				Bucket: events.S3Bucket{
					Name: "test-bucket",
				},
				Object: events.S3Object{
					Key: "test-path",
				},
			},
		})
	}
	return toReturn
}

func TestPerformIngest(t *testing.T) {

	reading := domain.HumidityStats{
		Reading:    50,
		Unit:       "percent",
		CustomerID: 5,
		SensorID:   8,
	}
	validHumidityBytes, err := json.Marshal([]domain.HumidityStats{reading})
	require.Nil(t, err)

	type temperatureReading struct {
		DegreesCelsius int
		CustomerID     int
		SensorID       int
	}

	tempReading := temperatureReading{
		DegreesCelsius: 5,
		CustomerID:     10,
		SensorID:       1,
	}

	invalidHumidityBytes, err := json.Marshal([]temperatureReading{tempReading})
	require.Nil(t, err)

	testCases := []struct {
		name               string
		s3Event            events.S3Event
		getObjectBytes     []byte
		getObjectErr       error
		writeHumidityError error
		shouldError        bool
	}{
		{
			name:           "unexpected data shape fails",
			s3Event:        getMockS3Event(1),
			getObjectBytes: invalidHumidityBytes,
			shouldError:    true,
		},

		{
			name:           "valid reading",
			s3Event:        getMockS3Event(1),
			getObjectBytes: validHumidityBytes,
		},
		{
			name:         "s3 error bubbles up",
			s3Event:      getMockS3Event(1),
			getObjectErr: fmt.Errorf("s3 connection problem"),
			shouldError:  true,
		},
		{
			name:               "db error bubbles up",
			s3Event:            getMockS3Event(1),
			getObjectBytes:     validHumidityBytes,
			writeHumidityError: fmt.Errorf("db rejected the records"),
			shouldError:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s3Reader := &s3Mock{
				getObjectBytes: tc.getObjectBytes,
				getObjectErr:   tc.getObjectErr,
			}
			dbWriter := &dbMock{
				writeHumdityErr: tc.writeHumidityError,
			}

			processor := NewS3SensorProcessor(s3Reader, dbWriter)

			err := processor.performIngest(context.Background(), tc.s3Event)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

		})
	}
}
