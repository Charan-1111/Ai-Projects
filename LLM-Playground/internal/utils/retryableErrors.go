package utils

func IsRetryableError(errCode int) bool {
	switch errCode {
	case 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}
