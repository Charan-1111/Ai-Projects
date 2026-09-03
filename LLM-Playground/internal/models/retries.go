package models

type Retries struct {
	InitialDelayMs int `json:"initial_delay_ms"`
	MaximumDelayMs int `json:"maximum_delay_ms"`
	MaxAttempts    int `json:"max_attempts"`
}