package entities

import "time"

type BeerLogs struct {
	ID          string    `json:"id"`
	Method      string    `json:"method"`
	RequestData string    `json:"request_data"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
}
