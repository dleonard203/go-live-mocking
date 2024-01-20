package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dleonard203/go-live-mocking/src/domain"
	"github.com/dleonard203/go-live-mocking/src/domain/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

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
		hasJSONErr         bool
		shouldError        bool
	}{
		{
			name:           "unexpected data shape fails",
			s3Event:        getMockS3Event(1),
			getObjectBytes: invalidHumidityBytes,
			hasJSONErr:     true,
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s3Reader := mocks.NewMockIS3Reader(ctrl)
			dbWriter := mocks.NewMockISensorIngetsor(ctrl)

			s3Reader.EXPECT().GetObjectContents(
				tc.s3Event.Records[0].S3.Bucket.Name,
				tc.s3Event.Records[0].S3.Object.Key,
			).Times(1).Return(tc.getObjectBytes, tc.getObjectErr)

			// db write only gets called if we get a successful s3 read for a valid payload
			if tc.getObjectErr == nil && !tc.hasJSONErr {
				dbWriter.EXPECT().WriteHumidityReadings(gomock.Any()).Times(1).Return(tc.writeHumidityError)
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
