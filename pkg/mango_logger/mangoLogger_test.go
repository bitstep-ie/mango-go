package mango_logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestLogger(cliEnabled, fileEnabled bool, strict bool, autoGenCorr bool) *MangoLogger {
	tmpFile, _ := os.CreateTemp("", "test-*.log")
	config := &LogConfig{
		Out: &OutConfig{
			Enabled: true,
			File: &FileOutputConfig{
				Enabled:    fileEnabled,
				Path:       tmpFile.Name(),
				MaxSize:    1,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   false,
			},
			Cli: &CliConfig{
				Enabled:        cliEnabled,
				Verbose:        true,
				Friendly:       true,
				VerboseFormat:  ".",
				FriendlyFormat: ".",
			},
			Syslog: &SyslogConfig{},
		},
		MangoConfig: &MangoConfig{
			Strict: strict,
			CorrelationId: &CorrelationIdConfig{
				Strict:       strict,
				AutoGenerate: autoGenCorr,
			},
		},
	}
	return NewMangoLogger(config)
}

func TestMangoLogger_AllLevels(t *testing.T) {
	levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}
	for _, lvl := range levels {
		t.Run(lvl.String(), func(t *testing.T) {
			logger := newTestLogger(true, false, false, true)

			record := slog.Record{
				Time:    time.Now(),
				Level:   lvl,
				Message: "Message " + lvl.String(),
			}

			ctx := context.Background()

			// Capture stdout/stderr depending on level
			oldOut := os.Stdout
			oldErr := os.Stderr
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			err := logger.Handle(ctx, record)
			assert.NoError(t, err)

			_ = wOut.Close()
			_ = wErr.Close()
			var bufOut, bufErr bytes.Buffer
			_, _ = bufOut.ReadFrom(rOut)
			_, _ = bufErr.ReadFrom(rErr)
			os.Stdout = oldOut
			os.Stderr = oldErr

			if lvl == slog.LevelDebug || lvl == slog.LevelInfo {
				assert.Contains(t, bufOut.String()+bufErr.String(), "Message "+lvl.String())
			} else {
				assert.Contains(t, bufErr.String(), "Message "+lvl.String())
			}
		})
	}
}

func TestMangoLogger_StrictMode_MissingRequiredField(t *testing.T) {
	logger := newTestLogger(true, false, true, false) // strict, no auto-gen

	record := slog.Record{
		Time:    time.Now(),
		Level:   slog.LevelInfo,
		Message: "Test strict",
	}

	ctx := context.Background() // no correlation ID

	_, err := logger.buildLog(ctx, record)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[STRICT_MODE ON]")
}

func TestMangoLogger_AutoGenerateCorrelation(t *testing.T) {
	logger := newTestLogger(true, false, true, true) // auto-generate correlation

	record := slog.Record{
		Time:    time.Now(),
		Level:   slog.LevelInfo,
		Message: "Test auto-gen correlation",
	}

	ctxWithVal := context.WithValue(context.Background(), OPERATION, "op1")
	ctxWithVal = context.WithValue(ctxWithVal, TYPE, "Business")
	ctxWithVal = context.WithValue(ctxWithVal, APPLICATION, "app")

	logOutput, err := logger.buildLog(ctxWithVal, record)
	assert.NoError(t, err)
	assert.NotEmpty(t, logOutput.Correlationid)
}

func TestMangoLogger_WrongType(t *testing.T) {
	logger := newTestLogger(true, false, true, true) // auto-generate correlation

	record := slog.Record{
		Time:    time.Now(),
		Level:   slog.LevelInfo,
		Message: "Test auto-gen correlation",
	}

	ctxWithVal := context.WithValue(context.Background(), OPERATION, "op1")
	ctxWithVal = context.WithValue(ctxWithVal, TYPE, "type")
	ctxWithVal = context.WithValue(ctxWithVal, APPLICATION, "app")

	_, err := logger.buildLog(ctxWithVal, record)
	assert.Error(t, err)
	assert.Equal(t, "[STRICT_MODE ON] without required context fields [type application operation] - [type] required in context and not present (or wrong type - expected string). Current value [type] is not in the allowed list: [\"Business\" \"Security\" \"Performance\"]", err.Error())
}

func TestMangoLogger_MergeAttrs(t *testing.T) {
	a1 := []slog.Attr{{Key: "a", Value: slog.StringValue("1")}}
	a2 := []slog.Attr{{Key: "b", Value: slog.StringValue("2")}, {Key: "a", Value: slog.StringValue("3")}}

	merged := mergeAttrs(a1, a2)
	assert.Len(t, merged, 2)

	m := make(map[string]string)
	for _, attr := range merged {
		m[attr.Key] = attr.Value.String()
	}

	assert.Equal(t, "3", m["a"]) // list2 takes precedence
	assert.Equal(t, "2", m["b"])
}

func TestFormatWithGoJQ_ErrorCases(t *testing.T) {
	// invalid JSON
	_, err := formatWithGoJQ("{invalid}", ".")
	assert.Error(t, err)

	// invalid jq
	_, err = formatWithGoJQ(`{"a":1}`, "???")
	assert.Error(t, err)
}

func TestHandleEachField_ExistingAndMissing(t *testing.T) {
	logger := newTestLogger(true, false, true, true)
	record := &StructuredLog{}

	ctxWithVal := context.WithValue(context.Background(), OPERATION, "op1")
	ctxMissing := context.Background()

	// Existing value
	err := handleEachField(ctxWithVal, record, OPERATION, *logger)
	assert.NoError(t, err)
	assert.Equal(t, "op1", record.Operation)

	// Missing value with auto-generate correlation
	err = handleEachField(ctxMissing, record, CORRELATION_ID, *logger)
	assert.NoError(t, err)
	assert.NotEmpty(t, record.Correlationid)
}

func TestWriteStringToLogFile_Disabled(t *testing.T) {
	logger := newTestLogger(true, false, false, true)
	err := logger.writeStringToLogFile("hello")
	assert.NoError(t, err) // should silently skip since file is disabled
}

func TestHandlePromptOutput_DefaultFallback(t *testing.T) {
	logger := newTestLogger(true, false, false, true)
	record := &StructuredLog{
		Level:   slog.LevelInfo,
		Message: "hello",
	}

	// Capture stdout
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := logger.handlePromptOutput(record, `{"message":"hello"}`)
	assert.NoError(t, err)

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldOut

	assert.Contains(t, buf.String(), "hello")
}

func TestHandleFileOutput_AllLevels(t *testing.T) {
	// Create a temporary file for logging
	tmpFile, err := os.CreateTemp("", "mango_logger_test_*.log")
	assert.NoError(t, err)
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
	}(tmpFile)

	// Logger config with file output enabled
	config := &LogConfig{
		Out: &OutConfig{
			Enabled: true,
			File: &FileOutputConfig{
				Enabled: true,
				Debug:   true,
				Path:    tmpFile.Name(),
			},
			Cli:    &CliConfig{Enabled: false},
			Syslog: &SyslogConfig{},
		},
		MangoConfig: &MangoConfig{
			Strict: false,
			CorrelationId: &CorrelationIdConfig{
				Strict:       false,
				AutoGenerate: true,
			},
		},
	}

	logger := NewMangoLogger(config)

	// Prepare a structured log for each level
	levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}

	for _, lvl := range levels {
		t.Run(lvl.String(), func(t *testing.T) {
			logOutput := &StructuredLog{
				Level:   lvl,
				Message: "Test message " + lvl.String(),
			}

			jsonOut, err := json.Marshal(logOutput)
			assert.NoError(t, err)

			err = logger.handleFileOutput(logOutput, string(jsonOut))
			assert.NoError(t, err)

			// Read file content
			content, err := os.ReadFile(tmpFile.Name())
			assert.NoError(t, err)
			assert.Contains(t, string(content), logOutput.Message)
		})
	}
}

func TestHandleFileOutput_FileDisabled(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "should_not_exist-*.log")
	// Logger config with file output disabled
	config := &LogConfig{
		Out: &OutConfig{
			Enabled: true,
			File: &FileOutputConfig{
				Enabled: false,
				Path:    tmpFile.Name(),
			},
			Cli:    &CliConfig{Enabled: false},
			Syslog: &SyslogConfig{},
		},
		MangoConfig: &MangoConfig{
			Strict: false,
			CorrelationId: &CorrelationIdConfig{
				Strict:       false,
				AutoGenerate: true,
			},
		},
	}

	logger := NewMangoLogger(config)

	logOutput := &StructuredLog{
		Level:   slog.LevelInfo,
		Message: "Should not write",
	}

	jsonOut := `{"Level":"INFO","Message":"Should not write"}`

	// Should not error even if file is disabled
	err := logger.handleFileOutput(logOutput, jsonOut)
	assert.NoError(t, err)
}

func TestHandleFileOutput_InvalidLevel(t *testing.T) {
	// Logger config with file enabled
	tmpFile, err := os.CreateTemp("", "mango_logger_test_invalid_level_*.log")
	assert.NoError(t, err)
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
	}(tmpFile)

	config := &LogConfig{
		Out: &OutConfig{
			Enabled: true,
			File: &FileOutputConfig{
				Enabled: true,
				Path:    tmpFile.Name(),
			},
			Cli:    &CliConfig{Enabled: false},
			Syslog: &SyslogConfig{},
		},
		MangoConfig: &MangoConfig{},
	}

	logger := NewMangoLogger(config)

	// Use a non-standard level
	logOutput := &StructuredLog{
		Level:   slog.Level(999), // invalid level
		Message: "Invalid level message",
	}

	jsonOut := `{"Level":999,"Message":"Invalid level message"}`

	err = logger.handleFileOutput(logOutput, jsonOut)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record level not one of")
}
