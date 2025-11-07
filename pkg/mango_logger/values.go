package mango_logger

// RFC3339NanoMC is the desired timestamp output format
const RFC3339NanoMC = "2006-01-02T15:04:05.999Z0700"

// All the different Log Types
const (
	BusinessType    = "Business"
	SecurityType    = "Security"
	PerformanceType = "Performance"
)

// Define a named type based on string
type ctxKey string

// Contract fields expected in the Context to be available for logging purposes
const (
	CORRELATION_ID ctxKey = "correlationid"
	TYPE           ctxKey = "type"
	APPLICATION    ctxKey = "application"
	OPERATION      ctxKey = "operation"
)

// ALLOWED_TYPES are the allowed values for TYPE
var ALLOWED_TYPES = []string{
	BusinessType,
	SecurityType,
	PerformanceType,
}

// REQUIRED_FIELDS are the fields checked against when MangoConfig.Strict is set
var REQUIRED_FIELDS = []ctxKey{
	TYPE,
	APPLICATION,
	OPERATION,
}
