package mango_logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/bitstep-ie/mango-go/pkg/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"log/slog"
	"log/syslog"
	"os"
	"testing"
	"time"
)

const TEST_OUT_FILE_PATH = "testPath.log"

var logConfig = &LogConfig{
	MangoConfig: &MangoConfig{
		Strict: true,
		CorrelationId: &CorrelationIdConfig{
			Strict:       true,
			AutoGenerate: false,
		},
	},
	Out: &OutConfig{
		Enabled: true,
		File: &FileOutputConfig{
			Enabled: true,
			Debug:   true,
			Path:    TEST_OUT_FILE_PATH,
		},
		Cli: &CliConfig{
			Enabled:  true,
			Friendly: true,
			Verbose:  false,
		},
		Syslog: &SyslogConfig{
			Facility: SyslogFacilityLocal0,
		},
	},
}

var removeTemp = func() {
	_ = os.Remove(TEST_OUT_FILE_PATH)
}

func TestCreateMangoLoggerDefaultFormat(t *testing.T) {
	mangoLogger := NewMangoLogger(logConfig)

	assert.NotNil(t, mangoLogger)
	assert.Empty(t, mangoLogger.attrs, "No attributes added to the logger")
	assert.NotNil(t, mangoLogger.LogWriter, "The log writer shouldn't be nil")
	assert.Equal(t, logConfig, mangoLogger.Config, "The config should match")
	assert.Equal(t, DefaultFriendlyFormat, mangoLogger.Config.Out.Cli.FriendlyFormat, "Default cli friendly format is applied to configuration")
	assert.Equal(t, DefaultVerboseFormat, mangoLogger.Config.Out.Cli.VerboseFormat, "Default verbose friendly format is applied to configuration")
}

func TestCreateMangoLoggerSpecificFormat(t *testing.T) {
	const cliFriendlyFormat = "CLI FRIENDLY FORMAT"
	const verboseFormat = "VERBOSE FORMAT"
	logConfig.Out.Cli.FriendlyFormat = cliFriendlyFormat
	logConfig.Out.Cli.VerboseFormat = verboseFormat

	mangoLogger := NewMangoLogger(logConfig)

	assert.NotNil(t, mangoLogger)
	assert.Empty(t, mangoLogger.attrs, "No attributes added to the logger")
	assert.NotNil(t, mangoLogger.LogWriter, "The log writer shouldn't be nil")
	assert.Equal(t, logConfig, mangoLogger.Config, "The config should match")
	assert.Equal(t, cliFriendlyFormat, mangoLogger.Config.Out.Cli.FriendlyFormat, "Default cli friendly format is applied to configuration")
	assert.Equal(t, verboseFormat, mangoLogger.Config.Out.Cli.VerboseFormat, "Default verbose friendly format is applied to configuration")

	// restore config to default
	logConfig.Out.Cli.FriendlyFormat = DefaultFriendlyFormat
	logConfig.Out.Cli.VerboseFormat = DefaultVerboseFormat
}

func TestMangoLogger_Enabled(t *testing.T) {
	mangoLogger := NewMangoLogger(logConfig)

	debugEnabled := mangoLogger.Enabled(context.Background(), slog.LevelDebug)
	infoEnabled := mangoLogger.Enabled(context.Background(), slog.LevelInfo)
	warnEnabled := mangoLogger.Enabled(context.Background(), slog.LevelWarn)
	errorEnabled := mangoLogger.Enabled(context.Background(), slog.LevelError)
	defaultEnabled := mangoLogger.Enabled(context.Background(), -1)
	assert.True(t, debugEnabled)
	assert.True(t, infoEnabled)
	assert.True(t, warnEnabled)
	assert.True(t, errorEnabled)
	assert.False(t, defaultEnabled)
}

func TestMangoLogger_WithAttrs(t *testing.T) {
	const boolKey = "boolKey"
	const boolValue = true
	const stringKey = "stringKey"
	const stringValue = "this is the value"
	mangoLogger := NewMangoLogger(logConfig)

	mangoLoggerWithAttrs := mangoLogger.WithAttrs([]slog.Attr{slog.Bool(boolKey, boolValue), slog.String(stringKey, stringValue)}).(MangoLogger)

	assert.Len(t, mangoLoggerWithAttrs.attrs, 2, "Two keys have been added to the logger")
	assert.Equal(t, boolKey, mangoLoggerWithAttrs.attrs[0].Key, "The key of the attribute")
	assert.Equal(t, boolValue, mangoLoggerWithAttrs.attrs[0].Value.Bool(), "The value of the attribute")
	assert.Equal(t, stringKey, mangoLoggerWithAttrs.attrs[1].Key, "The key of the attribute")
	assert.Equal(t, stringValue, mangoLoggerWithAttrs.attrs[1].Value.String(), "The value of the attribute")
}

func TestMangoLogger_WithGroup(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but no panic occurred")
		}
	}()
	mangoLogger := NewMangoLogger(logConfig)

	mangoLogger.WithGroup("groupName")
}

func TestMangoLogger_Handle_StrictModeOn_Failing(t *testing.T) {
	mangoLogger := NewMangoLogger(logConfig)

	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)
	logRecord.AddAttrs(slog.Int("i", 22))

	err := mangoLogger.Handle(context.Background(), logRecord)
	if err != nil {
		if !errors.Is(err, strictModeOn) {
			t.Errorf("Error expected from handling strict mode without enough context info")
		}
	} else {
		t.Errorf("Error expected from handling strict mode without enough context info")
	}
}

func TestMangoLogger_Handle_StrictModeOn_InvalidType(t *testing.T) {
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), CORRELATION_ID, uuid.New().String())
	logContext = context.WithValue(logContext, TYPE, "Random Not allowed Type")
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		if !errors.Is(err, strictModeOn) {
			t.Errorf("Error expected from handling strict mode without enough context info")
		}
	} else {
		t.Errorf("Error expected from handling strict mode without enough context info")
	}
}

func TestMangoLogger_Handle_StrictModeOn_NoApplication(t *testing.T) {
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), CORRELATION_ID, uuid.New().String())
	logContext = context.WithValue(logContext, TYPE, BusinessType)

	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		if !errors.Is(err, strictModeOn) {
			t.Errorf("Error expected from handling strict mode without enough context info")
		}
	} else {
		t.Errorf("Error expected from handling strict mode without enough context info")
	}
}

func TestMangoLogger_Handle_StrictModeOn_NoOperation(t *testing.T) {
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), CORRELATION_ID, uuid.New().String())
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		if !errors.Is(err, strictModeOn) {
			t.Errorf("Error expected from handling strict mode without enough context info")
		}
	} else {
		t.Errorf("Error expected from handling strict mode without enough context info")
	}
}

func TestMangoLogger_Handle_StrictModeOn_NoCorrelationGenerated(t *testing.T) {
	logConfig.MangoConfig.CorrelationId.AutoGenerate = false
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		if !errors.Is(err, strictModeOn) {
			t.Errorf("Error expected from handling strict mode without enough context info")
		}
	} else {
		t.Errorf("Error expected from handling strict mode without enough context info")
	}
}

func TestMangoLogger_Handle_StrictModeOn_CorrelationGenerated_ValidInfo(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_ProvidedCorrelation_ValidInfo(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = false
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")
	logContext = context.WithValue(logContext, CORRELATION_ID, "c8a67865-29dd-4e1e-a1a3-990d3075e3bc")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_UnknownLevel(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.Level(-8), "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		_ = w.Close()
		gotStdOut := readPipe(r)
		_ = wStdErr.Close()
		gotStdErr := readPipe(rStdErr)

		assert.ErrorContains(t, err, "record level not one of: debug, info, warn or error")
		expectedStdOut := "Record level not one of: debug, info, warn or error\n"
		assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr
		assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")
		assert.NoFileExists(t, TEST_OUT_FILE_PATH, "No log file should be created")
	} else {
		t.Errorf("Error expected from handling this log record as the log level is not supported")
	}

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_UnknownLevel_NoCLI_FileOnly(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Cli.Enabled = false
	logConfig.Out.File.Enabled = true
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.Level(-8), "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		_ = w.Close()
		gotStdOut := readPipe(r)
		_ = wStdErr.Close()
		gotStdErr := readPipe(rStdErr)

		assert.ErrorContains(t, err, "record level not one of: debug, info, warn or error")
		expectedStdOut := "Record level not one of: debug, info, warn or error\n"
		assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr
		assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

		assert.NoFileExists(t, TEST_OUT_FILE_PATH, "No log file should be created")
	} else {
		t.Errorf("Error expected from handling this log record as the log level is not supported")
	}

	logConfig.Out.Cli.Enabled = true
	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_CLINonFriendly_ValidInfo(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Cli.Friendly = false
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	var logEntry StructuredLog
	err = json.Unmarshal([]byte(gotStdOut), &logEntry)
	if err != nil {
		t.Fatalf("Error unmarshalling log file: %v", err)
	}

	assert.Equal(t, logEntry.Level, slog.LevelInfo, "Level in file log entry should match")
	assert.Equal(t, logEntry.Timestamp, "0001-01-01T00:00:00Z", "Timestamp in file log entry should match")
	assert.Equal(t, logEntry.Type, "Business", "Type in file log entry should match")
	assert.Equal(t, logEntry.Application, "testApp", "Application in file log entry should match")
	assert.Equal(t, logEntry.Operation, "operation", "Operation in file log entry should match")
	testutils.AssertValidUUID(t, logEntry.Correlationid, "correlationId")
	testutils.AssertValidUUID(t, logEntry.LogId, "logId")
	assert.Equal(t, logEntry.Message, "hello from sub1", "Message in file log entry should match")
	assert.Empty(t, logEntry.Attributes, "No Attributes should be set")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	logConfig.MangoConfig.CorrelationId.AutoGenerate = false
	logConfig.Out.Cli.Friendly = true
	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_ValidDebug_NoStdOut(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelDebug, "This is a debug message", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "" // no stdout output of debug
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	validateLogFile(t, slog.Level(-4), "This is a debug message")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_ValidDebug_Verbose_VerboseFormatSame(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Cli.Verbose = true
	logConfig.Out.Cli.VerboseFormat = DefaultFriendlyFormat
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelDebug, "This is a debug message", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[DEBUG] - 0001-01-01T00:00:00Z - operation - This is a debug message - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	validateLogFile(t, slog.Level(-4), "This is a debug message")

	logConfig.Out.Cli.Verbose = false
	logConfig.Out.Cli.VerboseFormat = DefaultVerboseFormat

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_ValidDebug_Verbose_DefaultVerboseFormat(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Cli.Verbose = true

	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelDebug, "This is a debug message", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	var logEntry StructuredLog
	err = json.Unmarshal([]byte(gotStdOut), &logEntry)
	if err != nil {
		t.Fatalf("Error unmarshalling log file: %v", err)
	}

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr
	assert.Equal(t, logEntry.Level, slog.Level(-4), "Level in file log entry should match")
	assert.Equal(t, logEntry.Timestamp, "0001-01-01T00:00:00Z", "Timestamp in file log entry should match")
	assert.Equal(t, logEntry.Type, "Business", "Type in file log entry should match")
	assert.Equal(t, logEntry.Application, "testApp", "Application in file log entry should match")
	assert.Equal(t, logEntry.Operation, "operation", "Operation in file log entry should match")
	testutils.AssertValidUUID(t, logEntry.Correlationid, "correlationId")
	testutils.AssertValidUUID(t, logEntry.LogId, "logId")
	assert.Equal(t, logEntry.Message, "This is a debug message", "Message in file log entry should match")
	assert.Empty(t, logEntry.Attributes, "No Attributes should be set")

	validateLogFile(t, slog.Level(-4), "This is a debug message")

	logConfig.Out.Cli.Verbose = false
	logConfig.Out.Cli.VerboseFormat = DefaultVerboseFormat

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_ValidWarn(t *testing.T) {
	originalStdout := os.Stdout      // Save the original os.Stdout
	rStdOut, wStdOut, _ := os.Pipe() // Create a pipe and redirect os.Stdout to it
	os.Stdout = wStdOut
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	mangoLogger := NewMangoLogger(logConfig)
	msg := "warning message here"
	logRecord := slog.NewRecord(time.Time{}, slog.LevelWarn, msg, 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = wStdOut.Close()
	gotStdOut := readPipe(rStdOut)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	expectedStdErr := "\"[WARN] - 0001-01-01T00:00:00Z - operation - " + msg + " - {}\"\n"
	assert.Equal(t, expectedStdErr, gotStdErr, "Log output to stdout should match")
	assert.Equal(t, "", gotStdOut, "Log output to stdout should match")

	validateLogFile(t, slog.LevelWarn, msg)

	logConfig.Out.Cli.Verbose = false
	logConfig.Out.Cli.VerboseFormat = DefaultVerboseFormat

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_ValidError(t *testing.T) {
	originalStdout := os.Stdout      // Save the original os.Stdout
	rStdOut, wStdOut, _ := os.Pipe() // Create a pipe and redirect os.Stdout to it
	os.Stdout = wStdOut
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	mangoLogger := NewMangoLogger(logConfig)
	msg := "Something went wrong so we can now log it"
	logRecord := slog.NewRecord(time.Time{}, slog.LevelError, msg, 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = wStdOut.Close()
	gotStdOut := readPipe(rStdOut)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	expectedStdErr := "\"[ERROR] - 0001-01-01T00:00:00Z - operation - " + msg + " - {}\"\n"
	assert.Equal(t, expectedStdErr, gotStdErr, "Log output to stdout should match")
	assert.Equal(t, "", gotStdOut, "Log output to stdout should match")

	validateLogFile(t, slog.LevelError, msg)

	logConfig.Out.Cli.Verbose = false
	logConfig.Out.Cli.VerboseFormat = DefaultVerboseFormat

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_NoCLIFriendly_ValidError(t *testing.T) {
	originalStdout := os.Stdout      // Save the original os.Stdout
	rStdOut, wStdOut, _ := os.Pipe() // Create a pipe and redirect os.Stdout to it
	os.Stdout = wStdOut
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Cli.Friendly = false
	mangoLogger := NewMangoLogger(logConfig)
	msg := "Something went wrong so we can now log it"
	logRecord := slog.NewRecord(time.Time{}, slog.LevelError, msg, 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = wStdOut.Close()
	gotStdOut := readPipe(rStdOut)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	var logEntry StructuredLog
	err = json.Unmarshal([]byte(gotStdErr), &logEntry)
	if err != nil {
		t.Fatalf("Error unmarshalling log file: %v", err)
	}

	assert.Equal(t, logEntry.Level, slog.LevelError, "Level in file log entry should match")
	assert.Equal(t, logEntry.Timestamp, "0001-01-01T00:00:00Z", "Timestamp in file log entry should match")
	assert.Equal(t, logEntry.Type, "Business", "Type in file log entry should match")
	assert.Equal(t, logEntry.Application, "testApp", "Application in file log entry should match")
	assert.Equal(t, logEntry.Operation, "operation", "Operation in file log entry should match")
	testutils.AssertValidUUID(t, logEntry.Correlationid, "correlationId")
	testutils.AssertValidUUID(t, logEntry.LogId, "logId")
	assert.Equal(t, logEntry.Message, msg, "Message in file log entry should match")
	assert.Empty(t, logEntry.Attributes, "No Attributes should be set")

	assert.Equal(t, "", gotStdOut, "Log output to stdout should match")

	validateLogFile(t, slog.LevelError, msg)

	logConfig.Out.Cli.Verbose = false
	logConfig.Out.Cli.Friendly = true
	logConfig.Out.Cli.VerboseFormat = DefaultVerboseFormat

	t.Cleanup(removeTemp)
}

func TestMangoLogger_EffectivelyNoLogging(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	logConfig.Out.File.Enabled = false
	logConfig.Out.Cli.Enabled = false
	logConfig.Out.Syslog.Facility = ""

	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)

	expectedStdOut := "Effectively no logging enabled! The config.out.file.enabled, config.out.cli.enabled and config.out.syslog.facility flags are all false.\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.NoFileExists(t, TEST_OUT_FILE_PATH, "Log output file does not exist")

	// restore config
	logConfig.Out.File.Enabled = true
	logConfig.Out.Cli.Enabled = true

	t.Cleanup(removeTemp)
}

func TestMangoLogger_OutDisabled(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	logConfig.Out.Enabled = false

	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)

	expectedStdOut := "No logging enabled! Check config.out.enabled.\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.NoFileExists(t, TEST_OUT_FILE_PATH, "Log output file does not exist")

	// restore config
	logConfig.Out.Enabled = true

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Kern(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityKern
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_KERN|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_User(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityUser
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_USER|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Mail(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityMail
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_MAIL|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Daemon(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityDaemon
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_DAEMON|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Auth(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityAuth
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_AUTH|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Syslog(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilitySyslog
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_SYSLOG|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_News(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityNews
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_NEWS|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Uucp(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityUucp
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_UUCP|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Cron(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityCron
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_CRON|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Authpriv(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityAuthpriv
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_AUTHPRIV|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Ftp(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityFtp
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_FTP|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Local0(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityLocal0
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_LOCAL0|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Local1(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityLocal1
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_LOCAL1|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Local2(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityLocal2
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_LOCAL2|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Local3(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityLocal3
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_LOCAL3|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Local4(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityLocal4
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_LOCAL4|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Local5(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityLocal5
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_LOCAL5|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Local6(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityLocal6
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_LOCAL6|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Local7(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Syslog.Facility = SyslogFacilityLocal7
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		t.Errorf("Error not expected from handling this log record")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	assert.Equal(t, syslog.LOG_LOCAL7|syslog.LOG_INFO, logConfig.Out.Syslog.priority, "syslog.priority should match calculated priority")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func TestMangoLogger_Handle_StrictModeOn_Facility_Invalid(t *testing.T) {
	originalStdout := os.Stdout // Save the original os.Stdout
	r, w, _ := os.Pipe()        // Create a pipe and redirect os.Stdout to it
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	originalStderr := os.Stderr
	rStdErr, wStdErr, _ := os.Pipe() // Create a pipe and redirect os.Stderr to it
	os.Stderr = wStdErr
	defer func() { os.Stderr = originalStderr }()

	logConfig.MangoConfig.CorrelationId.AutoGenerate = true
	logConfig.Out.Cli.Friendly = true
	logConfig.Out.Syslog.Facility = "Test Failure"
	mangoLogger := NewMangoLogger(logConfig)
	logRecord := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello from sub1", 0)

	logContext := context.WithValue(context.Background(), OPERATION, "operation")
	logContext = context.WithValue(logContext, TYPE, BusinessType)
	logContext = context.WithValue(logContext, APPLICATION, "testApp")

	// Action to test
	err := mangoLogger.Handle(logContext, logRecord)
	if err != nil {
		_ = w.Close()
		_ = wStdErr.Close()
		gotStdErr := readPipe(rStdErr)

		assert.ErrorContains(t, err, "facility level not valid")
		assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr
	} else {
		t.Errorf("Error expected from handling invalid syslog facility")
	}

	_ = w.Close()
	gotStdOut := readPipe(r)
	_ = wStdErr.Close()
	gotStdErr := readPipe(rStdErr)

	assert.Equal(t, "", gotStdErr, "Stderr should be empty") // nothing to stderr

	expectedStdOut := "\"[INFO] - 0001-01-01T00:00:00Z - operation - hello from sub1 - {}\"\nFacility level not valid\n"
	assert.Equal(t, expectedStdOut, gotStdOut, "Log output to stdout should match")

	validateLogFile(t, slog.Level(0), "hello from sub1")

	t.Cleanup(removeTemp)
}

func readPipe(r *os.File) string {
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func validateLogFile(t *testing.T, level slog.Level, message string) {
	content, err := os.ReadFile(TEST_OUT_FILE_PATH)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	gotFile := string(content)
	var logEntry StructuredLog
	err = json.Unmarshal([]byte(gotFile), &logEntry)
	if err != nil {
		t.Fatalf("Error unmarshalling log file: %v", err)
	}

	assert.Equal(t, logEntry.Level, level, "Level in file log entry should match")
	assert.Equal(t, logEntry.Timestamp, "0001-01-01T00:00:00Z", "Timestamp in file log entry should match")
	assert.Equal(t, logEntry.Type, "Business", "Type in file log entry should match")
	assert.Equal(t, logEntry.Application, "testApp", "Application in file log entry should match")
	assert.Equal(t, logEntry.Operation, "operation", "Operation in file log entry should match")
	testutils.AssertValidUUID(t, logEntry.Correlationid, "correlationId")
	testutils.AssertValidUUID(t, logEntry.LogId, "logId")
	assert.Equal(t, logEntry.Message, message, "Message in file log entry should match")
	assert.Empty(t, logEntry.Attributes, "No Attributes should be set")

}
