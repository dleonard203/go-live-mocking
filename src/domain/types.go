package domain

// HumidityStats contains information from a humidity sensor reading
type HumidityStats struct {
	Reading    int    `json:"reading"`
	Unit       string `json:"unit"`
	CustomerID int    `json:"customer_id"`
	SensorID   int    `json:"sensor_id"`
}
