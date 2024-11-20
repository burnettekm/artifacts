package api

import "time"

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Cooldown struct {
	TotalSeconds     int       `json:"total_seconds"`
	RemainingSeconds int       `json:"remaining_seconds"`
	StartedAt        time.Time `json:"started_at"`
	Expiration       time.Time `json:"expiration"`
	Reason           string    `json:"reason"`
}

type Content struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

type Item struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}
