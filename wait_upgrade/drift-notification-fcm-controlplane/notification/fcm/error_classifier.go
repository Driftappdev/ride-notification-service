package fcm

import "strings"

type ClassifiedError struct {
	Code        string
	Retryable   bool
	RemoveToken bool
}

func ClassifyError(err error) ClassifiedError {
	if err == nil { return ClassifiedError{} }

	msg := strings.ToLower(err.Error())
	out := ClassifiedError{Code: "FCM_UNKNOWN"}

	switch {
	case strings.Contains(msg, "registration-token-not-registered"),
		strings.Contains(msg, "unregistered"),
		strings.Contains(msg, "not registered"):
		out.Code = "FCM_UNREGISTERED"
		out.RemoveToken = true
	case strings.Contains(msg, "invalid-registration-token"),
		strings.Contains(msg, "invalid argument"):
		out.Code = "FCM_INVALID_ARGUMENT"
		out.RemoveToken = true
	case strings.Contains(msg, "quota"),
		strings.Contains(msg, "too many requests"),
		strings.Contains(msg, "resource exhausted"):
		out.Code = "FCM_QUOTA_EXCEEDED"
		out.Retryable = true
	case strings.Contains(msg, "unavailable"),
		strings.Contains(msg, "deadline"),
		strings.Contains(msg, "timeout"),
		strings.Contains(msg, "internal"):
		out.Code = "FCM_TEMPORARY_FAILURE"
		out.Retryable = true
	case strings.Contains(msg, "sender id mismatch"):
		out.Code = "FCM_SENDER_ID_MISMATCH"
	default:
		out.Code = "FCM_UNKNOWN"
	}
	return out
}
