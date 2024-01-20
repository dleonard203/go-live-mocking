package main

import (
	"testing"

	"github.com/dleonard203/go-live-mocking/src/domain/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestDisplay(t *testing.T) {
	testCases := []struct {
		name        string
		getObjBytes []byte
		getObjErr   error
		panicReason any
		shouldError bool
	}{
		{
			name:        "validate panic handler",
			panicReason: "feeling blue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s3Reader := mocks.NewMockIS3Reader(ctrl)
			panicHandler := mocks.NewMockIPanicHandler(ctrl)

			// set up the mocks depending on if we panic or not
			if tc.panicReason != nil {
				s3Reader.EXPECT().GetObjectContents(gomock.Any(), gomock.Any()).Times(1).Do(func(bucket, key string) {
					panic(tc.panicReason)
				})
				panicHandler.EXPECT().Notify(tc.panicReason).Times(1)
			} else {
				s3Reader.EXPECT().GetObjectContents(gomock.Any(), gomock.Any()).Times(1).Return(tc.getObjBytes, tc.getObjErr)
			}

			reader := NewReader(s3Reader, panicHandler)

			_, err := reader.Display("bucket", "key")

			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

		})
	}
}
