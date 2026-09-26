package constants

import "time"

const (
	ClientMaxIdleConns        = 100
	ClientMaxIdleConnsPerHost = 20
	ClientMaxConnsPerHost     = 50
	ClientIdleConnTimeout     = 90 * time.Second
	ClientTimeout             = 10 * time.Second
)