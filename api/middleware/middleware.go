package middleware

func IsInternalServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

func IsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
