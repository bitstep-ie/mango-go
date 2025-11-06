package mango_logger

import "log/slog"

// StructuredLog is the structure of every log entry (output)
type StructuredLog struct {
	// Timestamp of the log entry
	Timestamp string `json:"ts"`

	// Type of the log entry. One of: Security, Business, or Performance
	Type string `json:"type"`

	// Application of which this log entry belongs. Should be corresponding with TAT
	Application string `json:"application"`

	// Operation is synonymous with the application/system's function or method.
	// Consider such examples: search, create, health, user_registration, checkout, token_issueance, case-status, etc.
	Operation string `json:"operation"`

	// Correlationid from the caller or self generated allowing to relate different systems around one
	Correlationid string `json:"correlationid"`

	// LogId is a unique identifier for each log entry - Helps in referring to logs when searching
	LogId string `json:"logId"`

	// Level of the log entry (slog.Debug, slog.Info, slog.Warn, slog.Error)
	Level slog.Level `json:"level"`

	// Message is the actual message of the log entry
	Message any `json:"message"`

	// Attributes set with slog or on the logger
	Attributes map[string]interface{} `json:"attributes"`
}

// Helper function to convert []slog.Attr to a map[string]interface{}
func ToMap(attrs []slog.Attr) map[string]interface{} {
	result := make(map[string]interface{})
	for _, attr := range attrs {
		result[attr.Key] = attr.Value.Any()
	}
	return result
}
