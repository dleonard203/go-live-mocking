package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/dleonard203/go-live-mocking/src/domain"
	"github.com/pkg/errors"
)

// SensorWriter writes sensor readings to postgres
type SensorWriter struct {
	connection   *sql.DB
	writeTimeout time.Duration
}

// NewSensorWriter creates a new SensorWriter
func NewSensorWriter(
	dsn string, writeTimeout time.Duration,
) (*SensorWriter, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to the database")
	}

	return &SensorWriter{
		writeTimeout: writeTimeout,
		connection:   db,
	}, nil
}

func (s *SensorWriter) getCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.writeTimeout)
}

// WriteHumidityReadings writes the readings to the database
func (s *SensorWriter) WriteHumidityReadings(readings []domain.HumidityStats) error {
	for _, reading := range readings {

		err := func() error {
			// in-line func so cancel is called at the end of each write
			ctx, cancel := s.getCtx()
			defer cancel()
			_, err := s.connection.ExecContext(
				ctx,
				`INSERT INTO humidity (reading, unit, customer_id, sensor_id) 
				    VALUES (?, ?, ?, ?)`,
				reading.Reading, reading.Unit, reading.CustomerID, reading.SensorID,
			)
			return err
		}()

		if err != nil {
			return err
		}
	}

	return nil
}
