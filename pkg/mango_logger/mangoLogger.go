package mango_logger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/itchyny/gojq"
	"github.com/natefinch/lumberjack"
	"log/slog"
	"os"
	"slices"
	"unsafe"
)

type MangoLogger struct {
	attrs     []slog.Attr
	Config    *LogConfig
	LogWriter *lumberjack.Logger
}

var StrictModeOn = errors.New(fmt.Sprintf("[STRICT_MODE ON] without required context fields %v", REQUIRED_FIELDS))

func NewMangoLogger(config *LogConfig) *MangoLogger {
	// Future idea to have multiple "appenders" in the mangoLogger that one can add, each with it's own logging configuration that it looks at
	logger := &MangoLogger{
		Config: applyDefaultFormats(*config),
		LogWriter: &lumberjack.Logger{
			Filename:   config.Out.File.Path,
			MaxSize:    config.Out.File.MaxSize,
			MaxBackups: config.Out.File.MaxBackups,
			MaxAge:     config.Out.File.MaxAge,
			Compress:   config.Out.File.Compress,
		},
	}
	return logger
}

// applyDefaultFormats to the configuration to ensure verbose and cli-friendly default formats are applied
func applyDefaultFormats(config LogConfig) *LogConfig {
	merged := config
	if config.Out.Cli.VerboseFormat == "" {
		merged.Out.Cli.VerboseFormat = DefaultVerboseFormat
	}
	if config.Out.Cli.FriendlyFormat == "" {
		merged.Out.Cli.FriendlyFormat = DefaultFriendlyFormat
	}
	return &merged
}

func (sl MangoLogger) Enabled(context context.Context, level slog.Level) bool {
	switch level {
	case slog.LevelDebug:
		fallthrough
	case slog.LevelInfo:
		fallthrough
	case slog.LevelWarn:
		fallthrough
	case slog.LevelError:
		return true
	default:
		return false
	}
}

func (sl MangoLogger) Handle(context context.Context, record slog.Record) error {
	if !sl.Config.Out.Enabled { // no logging enabled
		fmt.Println("No logging enabled! Check config.out.enabled.")
		return nil
	}

	if !sl.Config.Out.File.Enabled && !sl.Config.Out.Cli.Enabled && sl.Config.Out.Syslog.Facility == "" {
		fmt.Println("Effectively no logging enabled! The config.out.file.enabled, config.out.cli.enabled and config.out.syslog.facility flags are all false.")
		return nil
	}

	log, err := sl.buildLog(context, record)
	if err != nil {
		return err
	}

	jsonOut, err := json.Marshal(log)
	if err != nil {
		fmt.Println("Failed to marshal the StructuredLog. Internal error, should never happen")
		return err
	}

	if sl.Config.Out.Cli.Enabled {
		err := sl.handlePromptOutput(log, string(jsonOut))
		if err != nil {
			return err
		}
	}

	if sl.Config.Out.File.Enabled {
		err := sl.handleFileOutput(log, string(jsonOut))
		if err != nil {
			return err
		}
	}

	if sl.Config.Out.Syslog.Facility != "" {
		err := sl.handleSyslogOutput(log, jsonOut)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sl MangoLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	sl.attrs = append(sl.attrs, attrs...)
	return sl
}

func (sl MangoLogger) WithGroup(name string) slog.Handler {
	panic("unimplemented") // add a way of storing a generic map of maps of attributes
}

func (sl MangoLogger) writeStringToLogFile(s string) error {
	if sl.Config.Out.Enabled {
		s += "\n"
		b := unsafe.Slice(unsafe.StringData(s), len(s))
		_, err := sl.LogWriter.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func formatWithGoJQ(obj string, query string) (string, error) {
	// Unmarshal the JSON into an interface{} so gojq can process it
	var objInterface map[string]interface{}
	if err := json.Unmarshal([]byte(obj), &objInterface); err != nil {
		return "", err
	}

	// Parse the jq query
	jqQuery, err := gojq.Parse(query)
	if err != nil {
		return "", err
	}

	// Create a jq code executor
	iter := jqQuery.Run(objInterface)

	// Retrieve the result
	var result interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", err
		}
		result = v
	}

	// Convert result to string
	resultStr, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(resultStr), nil
}

func (sl MangoLogger) handleFileOutput(log *StructuredLog, jsonOut string) error {
	switch log.Level {
	case slog.LevelDebug:
		if sl.Config.Out.File.Debug {
			return sl.writeStringToLogFile(jsonOut)
		}
	case slog.LevelInfo:
		return sl.writeStringToLogFile(jsonOut)
	case slog.LevelWarn:
		fallthrough
	case slog.LevelError:
		return sl.writeStringToLogFile(jsonOut)
	default:
		fmt.Println("Record level not one of: debug, info, warn or error")
		return fmt.Errorf("record level not one of: debug, info, warn or error")
	}
	return nil
}

func (sl MangoLogger) handlePromptOutput(log *StructuredLog, jsonOut string) error {
	switch log.Level {
	case slog.LevelDebug:
		if sl.Config.Out.Cli.Verbose {
			result, _ := formatWithGoJQ(jsonOut, sl.Config.Out.Cli.VerboseFormat)
			_, _ = fmt.Fprintln(os.Stdout, result)
		}
	case slog.LevelInfo:
		if sl.Config.Out.Cli.Friendly {
			result, _ := formatWithGoJQ(jsonOut, sl.Config.Out.Cli.FriendlyFormat)
			_, _ = fmt.Fprintln(os.Stdout, result)
		} else {
			_, _ = fmt.Fprintln(os.Stdout, jsonOut)
		}
	case slog.LevelWarn:
		fallthrough
	case slog.LevelError:
		if sl.Config.Out.Cli.Friendly {
			result, _ := formatWithGoJQ(jsonOut, sl.Config.Out.Cli.FriendlyFormat)
			_, _ = fmt.Fprintln(os.Stderr, result)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, jsonOut)
		}
	default:
		fmt.Println("Record level not one of: debug, info, warn or error")
		return fmt.Errorf("record level not one of: debug, info, warn or error")
	}
	return nil
}

// mergeAttrs with list2 taking precedence
func mergeAttrs(list1, list2 []slog.Attr) []slog.Attr {
	attrMap := make(map[string]slog.Attr)
	for _, attr := range list1 {
		attrMap[attr.Key] = attr
	}
	for _, attr := range list2 {
		attrMap[attr.Key] = attr
	}
	mergedAttrs := make([]slog.Attr, 0, len(attrMap))
	for _, attr := range attrMap {
		mergedAttrs = append(mergedAttrs, attr)
	}

	return mergedAttrs
}

func getAllAttrs(record slog.Record) []slog.Attr {
	var attrs []slog.Attr

	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true // Continue the iteration
	})

	return attrs
}

func (sl MangoLogger) handleRequiredFields(context context.Context, logOutput *StructuredLog) error {
	if sl.Config.MangoConfig.CorrelationId.Strict {
		REQUIRED_FIELDS = append(REQUIRED_FIELDS, CORRELATION_ID)
	}
	for _, label := range REQUIRED_FIELDS {
		err := handleEachField(context, logOutput, label, sl)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleEachField(context context.Context, logOutput *StructuredLog, label ctxKey, sl MangoLogger) error {
	if value, ok := context.Value(label).(string); !ok {
		err := handleValueMissing(label, sl, logOutput)
		if err != nil {
			return err
		}
	} else {
		err := handleExistentValues(label, logOutput, value, sl)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleValueMissing(label ctxKey, sl MangoLogger, logOutput *StructuredLog) error {
	if CORRELATION_ID == label {
		if sl.Config.MangoConfig.CorrelationId.AutoGenerate {
			logOutput.Correlationid = uuid.New().String() // generate new UUID for correlation if missing from context
		} else {
			return fmt.Errorf("%w - required in context and not present (or wrong type - expected string). This can be added by doing: context.WithValue(newCtx, mangologger.%s, \"desiredValue\")", StrictModeOn, label)
		}
	} else {
		if sl.Config.MangoConfig.Strict {
			return fmt.Errorf("%w - required in context and not present (or wrong type - expected string). This can be added by doing: context.WithValue(newCtx, mangologger.%s, \"desiredValue\")", StrictModeOn, label)
		}
	}
	return nil
}

func handleExistentValues(label ctxKey, logOutput *StructuredLog, value string, sl MangoLogger) error {
	switch label { // set actual
	case OPERATION:
		logOutput.Operation = value
	case APPLICATION:
		logOutput.Application = value
	case TYPE:
		if sl.Config.MangoConfig.Strict {
			if !slices.Contains(ALLOWED_TYPES, value) {
				return fmt.Errorf("%w - [%s] required in context and not present (or wrong type - expected string). Current value [%s] is not in the allowed list: %+q", StrictModeOn, label, value, ALLOWED_TYPES)
			}
		}
		logOutput.Type = value
	}
	return nil
}

func (sl MangoLogger) buildLog(context context.Context, record slog.Record) (*StructuredLog, error) {
	logOutput := sl.makeBaseLog(record)

	err := sl.handleRequiredFields(context, logOutput)
	if err != nil {
		fmt.Printf("Required fields are not present. %s\n", err.Error())
		return logOutput, err
	}

	if value, ok := context.Value(CORRELATION_ID).(string); ok {
		logOutput.Correlationid = value
	}

	return logOutput, nil
}

func (sl MangoLogger) makeBaseLog(record slog.Record) *StructuredLog {
	logOutput := &StructuredLog{}
	logOutput.Timestamp = record.Time.Format(RFC3339NanoMC)
	logOutput.LogId = uuid.New().String() // generate a new UUID for each log entry
	logOutput.Level = record.Level
	logOutput.Operation = "unknownOperation"
	logOutput.Application = "unknownApplication"
	logOutput.Type = "unknownType"
	logOutput.Correlationid = ""
	logOutput.Message = record.Message
	logOutput.Attributes = ToMap(mergeAttrs(sl.attrs, getAllAttrs(record)))
	return logOutput
}
