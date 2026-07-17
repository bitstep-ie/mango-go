// Package logger is a specific logging library on top of slog with additional goodness
package logger

// Default output formats
const (
	// DefaultVerboseFormat is the default format for verbose (DEBUG to stdout) output
	DefaultVerboseFormat = "."

	// DefaultFriendlyFormat is the default format for all CLI friendly output (INFO and above to stdout)
	DefaultFriendlyFormat = `"[\(.level)] - \(.ts) - \(.operation) - \(.message) - \(.attributes)"`
)

type SyslogFacility string

const (
	SyslogFacilityKern     = "kern"
	SyslogFacilityUser     = "user"
	SyslogFacilityMail     = "mail"
	SyslogFacilityDaemon   = "daemon"
	SyslogFacilityAuth     = "auth"
	SyslogFacilitySyslog   = "syslog"
	SyslogFacilityNews     = "news"
	SyslogFacilityUucp     = "uucp"
	SyslogFacilityCron     = "cron"
	SyslogFacilityAuthpriv = "authpriv"
	SyslogFacilityFtp      = "ftp"
	SyslogFacilityLocal0   = "local0"
	SyslogFacilityLocal1   = "local1"
	SyslogFacilityLocal2   = "local2"
	SyslogFacilityLocal3   = "local3"
	SyslogFacilityLocal4   = "local4"
	SyslogFacilityLocal5   = "local5"
	SyslogFacilityLocal6   = "local6"
	SyslogFacilityLocal7   = "local7"
)

// LogConfig is the main configuration struct for Mango logging
type LogConfig struct {
	// ContextConfig is the mango configuration node
	ContextConfig *ContextConfig `yaml:"context" json:"context"`

	// Out is the node holding configuration about the file output
	Out *OutConfig `yaml:"out" json:"out"`
}

type ContextConfig struct {
	// Strict Will enforce the fields defined to be required are set for each log entry
	Strict bool `yaml:"strict" json:"strict"`

	// Required are the list of context configuration fields to be enforced as part of the log context based on the strict flag
	Required *[]ContextConfigField `yaml:"required" json:"required"`
}

// OutConfig provides a structure for defining the configuration of all the logging output
type OutConfig struct {
	// Overall output enable - killer switch to all output of mangologger
	Enabled bool `yaml:"enabled" json:"enabled"`

	// File output configuration
	File *FileOutputConfig `yaml:"file" json:"file"`

	// Cli configuration node for CLI output options
	Cli *CliConfig `yaml:"cli" json:"cli"`

	// Syslog configuration node for Syslog output options
	Syslog *SyslogConfig `yaml:"syslog" json:"syslog"`
}

// ContextConfigField defines the configuration of correlationId across mangologger
type ContextConfigField struct {
	// Name is the name of the context field to be enforced as part of the log context based on the strict flag
	Name string `yaml:"name" json:"name"`

	Value string `yaml:"value" json:"value"`
	// AutoGenerate is a placeholder for now - TODO
	AutoGenerate bool `yaml:"auto-generate" json:"autoGenerate"`
}

type FileOutputConfig struct {
	// Enabled switches on printing out to file
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Debug allows debug printout to file
	Debug bool `yaml:"debug" json:"debug"`

	// Path is the log file name - It uses <processname>-lumberjack.log in os.TempDir() if empty.
	Path string `yaml:"path" json:"path"`

	// MaxSize in MB before rotating - It defaults to 100 megabytes
	MaxSize int `yaml:"max-size" json:"maxSize"`

	// MaxBackups is the number of old log files to keep - The default is to retain all old log files
	MaxBackups int `yaml:"max-backups" json:"maxBackups"`

	// MaxAge is the number of days to keep old log files - The default is not to remove old log files based on age
	MaxAge int `yaml:"max-age" json:"maxAge"`

	// Compress old log files - The default is not to perform compression
	Compress bool `yaml:"compress" json:"compress"`
}

type CliConfig struct {
	// Enabled allows stdout/stderr printouts
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Friendly enables a human friendly output to stdout/stderr
	// When false it outputs json format as in file output
	Friendly bool `yaml:"friendly" json:"friendly"`

	// FriendlyFormat of the output in normal scenarios and if Friendly enabled
	// Defaults to DefaultFriendlyFormat applied to all log statements info+
	FriendlyFormat string `yaml:"friendly-format" json:"friendlyFormat"`

	// Verbose Enable debug to come out to std out following the VerboseFormat
	Verbose bool `yaml:"verbose" json:"verbose"`

	// VerboseFormat of the DEBUG statements output in verbose mode
	// Defaults to print the whole json object of logger.StructuredLog (using DefaultVerboseFormat)
	VerboseFormat string `yaml:"verbose-format" json:"verboseFormat"`
}
