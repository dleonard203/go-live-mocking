package main

import (
	"encoding/json"

	"github.com/dleonard203/go-live-mocking/src/domain"
)

// Reader displays s3 contents in a human friendly manner
type Reader struct {
	s3Reader     domain.IS3Reader
	panicHandler domain.IPanicHandler
}

// NewReader creates a new Reader
func NewReader(s3Reader domain.IS3Reader, panicHandler domain.IPanicHandler) *Reader {
	return &Reader{
		s3Reader:     s3Reader,
		panicHandler: panicHandler,
	}
}

// Display returns a pretty-printed JSON string based on the s3 file's contents
func (r *Reader) Display(bucket, key string) (*string, error) {
	defer func() {
		reason := recover()
		if reason != nil {
			r.panicHandler.Notify(reason)
		}
	}()

	contents, err := r.s3Reader.GetObjectContents(bucket, key)
	if err != nil {
		return nil, err
	}

	dest := make(map[string]any)
	err = json.Unmarshal(contents, &dest)

	if err != nil {
		return nil, err
	}

	bytes, err := json.MarshalIndent(dest, "", "  ")
	if err != nil {
		return nil, err
	}

	toString := string(bytes)
	return &toString, nil
}

func main() {

}
