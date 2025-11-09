//go:build !windows

package mango_logger

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"log/syslog"
	"testing"
)

// mockSyslogWriter implements io.Writer and simulates a syslog.Writer
type mockSyslogWriter struct {
	writeErr error
	closeErr error
	written  []byte
}

func (m *mockSyslogWriter) Write(p []byte) (int, error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	m.written = append(m.written, p...)
	return len(p), nil
}

func (m *mockSyslogWriter) Close() error {
	return m.closeErr
}

// overrideable function for syslog.New to make testing easier
var syslogNew = func(p syslog.Priority, tag string) (*syslog.Writer, error) {
	return syslog.New(p, tag)
}

// createTestLogger constructs a MangoLogger with given facility
func createTestLogger(facility SyslogFacility) MangoLogger {
	return MangoLogger{
		Config: &LogConfig{
			Out: &OutConfig{
				Syslog: &SyslogConfig{
					Facility: facility,
				},
			},
		},
	}
}

func TestHandleSyslogOutput_ValidLevels(t *testing.T) {
	facilities := []SyslogFacility{
		SyslogFacilityKern,
		SyslogFacilityUser,
		SyslogFacilityMail,
		SyslogFacilityDaemon,
		SyslogFacilityAuth,
		SyslogFacilitySyslog,
		SyslogFacilityNews,
		SyslogFacilityUucp,
		SyslogFacilityCron,
		SyslogFacilityAuthpriv,
		SyslogFacilityFtp,
		SyslogFacilityLocal0,
		SyslogFacilityLocal1,
		SyslogFacilityLocal2,
		SyslogFacilityLocal3,
		SyslogFacilityLocal4,
		SyslogFacilityLocal5,
		SyslogFacilityLocal6,
		SyslogFacilityLocal7,
	}

	levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}

	for _, f := range facilities {
		for _, lvl := range levels {
			logger := createTestLogger(f)
			log := &StructuredLog{
				Level:       lvl,
				Application: "testApp",
			}
			err := logger.handleSyslogOutput(log, []byte(`{"msg":"hello"}`))
			assert.NoError(t, err, "facility %v level %v", f, lvl)
		}
	}
}

func TestHandleSyslogOutput_InvalidFacility(t *testing.T) {
	logger := createTestLogger("invalid_facility")
	log := &StructuredLog{
		Level:       slog.LevelInfo,
		Application: "testApp",
	}

	err := logger.handleSyslogOutput(log, []byte(`{"msg":"oops"}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "facility level not valid")
}

func TestHandleSyslogOutput_InvalidLevel(t *testing.T) {
	logger := createTestLogger(SyslogFacilityUser)
	log := &StructuredLog{
		Level:       slog.Level(999), // invalid
		Application: "testApp",
	}

	err := logger.handleSyslogOutput(log, []byte(`{"msg":"bad level"}`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record level not one of")
}

func TestHandleSyslogOutput_SyslogWriterCloseError(t *testing.T) {
	logger := createTestLogger(SyslogFacilityUser)
	log := &StructuredLog{
		Level:       slog.LevelInfo,
		Application: "testApp",
	}

	origSyslogNew := syslogNew
	syslogNew = func(p syslog.Priority, t string) (*syslog.Writer, error) {
		_ = &mockSyslogWriter{closeErr: errors.New("close failed")}
		return (*syslog.Writer)(nil), nil
	}
	defer func() { syslogNew = origSyslogNew }()

	err := logger.handleSyslogOutput(log, []byte(`{"msg":"close test"}`))
	assert.NoError(t, err)
}
