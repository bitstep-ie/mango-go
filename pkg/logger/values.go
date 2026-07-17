package logger

// RFC3339NanoMC is the desired timestamp output format
const RFC3339NanoMC = "2006-01-02T15:04:05.999Z0700"

// All the different Log Types
const (
	BusinessType    = "Business"
	SecurityType    = "Security"
	PerformanceType = "Performance"
)

// ALLOWED_TYPES are the allowed values for TYPE
var ALLOWED_TYPES = []string{
	BusinessType,
	SecurityType,
	PerformanceType,
}
