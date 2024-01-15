package domain

// IS3Reader is the set of functions needed to read from an AWS S3 Bucket
type IS3Reader interface {
	GetObjectContents(bucket, path string) ([]byte, error)
}

// ISensorIngetsor is the set of functions needed to ingest sensor readings
type ISensorIngetsor interface {
	WriteHumidityReadings(readings []HumidityStats) error
}
