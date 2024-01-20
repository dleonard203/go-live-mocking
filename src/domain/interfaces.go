package domain

//go:generate mockgen -destination=./mocks/mock_interfaces.go -package=mocks -source=interfaces.go

// IS3Reader is the set of functions needed to read from an AWS S3 Bucket
type IS3Reader interface {
	GetObjectContents(bucket, path string) ([]byte, error)
}

// ISensorIngetsor is the set of functions needed to ingest sensor readings
type ISensorIngetsor interface {
	WriteHumidityReadings(readings []HumidityStats) error
}

// IPanicHandler takes action when a function has panic'd
type IPanicHandler interface {
	Notify(reason any) error
}
