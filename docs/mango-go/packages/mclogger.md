#  mango4go - mclogger

Mango Logger is an enhancement package aimed to be flexible and lightweight for Go that uses slog and allows logging messages to both standard output (stdout) and file simultaneously. It supports various log levels (e.g., DEBUG, INFO, WARN, ERROR) to control the verbosity of log output. The package can be easily integrated into Go applications, providing timestamped log entries, log rotation, etc...

## Limitations

- `withGroup` a function of `slog` is not currently supported by mclogger.


## Usage

Create a new McLogger by using `mclogger.NewMcLogger` _constructor_ and pass in your desired configuration.

```go language=go
mcLogger := mclogger.NewMcLogger(&mclogger.LogConfig{<...>})
slog.SetDefault(slog.New(mcLogger))
```

Then throughout the application code feel free to make use of the standard `slog` methods as:

```go language=go
slog.Debug("message", "key", "val")
slog.Info("message", "key", "val")
slog.Warn("message", "key", "val")
slog.Error("message", "key", "val")
```

Or the context enabled equivalents:

```go language=go
slog.Debug(ctx, "message", "key", "val")
slog.Info(ctx, "message", "key", "val")
slog.Warn(ctx, "message", "key", "val")
slog.Error(ctx, "message", "key", "val")
```


## Configuration

This package is aimed to be highly configurable, and therefore defined a grouping of structs with `LogConfig` being the main Logging Configuration struct.

### Example

`yaml` example of the full configuration:

```yaml language=yaml
---
mango: # The mango specific configuration node
  strict: false # enforces the REQUIRED_FIELDS to be present in each log context
  correlationId:
    strict: false # enforces CorrelationId as part of the REQUIRED_FIELDS to be present in each log context
    autoGenerate: false # will generate a correlationId if missing from context.
out: # output
  enabled: false # overall output enabling flag
  file: # file output configurations
    enabled: false # flag on printing out to file
    debug: false # allows debug printout to file
    path: path/to/file.log # is the log file name - It uses <processname>-lumberjack.log in os.TempDir() if empty.
    maxSize: 100 # in MB before rotating - It defaults to 100 megabytes
    maxBackups: 5 # is the number of old log files to keep - The default is to retain all old log files
    maxAge: 0 # is the number of days to keep old log files - The default is not to remove old log files based on age
    compress: false # flag to compress old log files - The default is not to perform compression
  cli:
    enabled: false # allows stdout/stderr printouts
    friendly: false # enables a human friendly output to stdout/stderr
    #  When false it outputs json format as in file output
    verbose: false # also prints debug to stdout and includes correlation id
  syslog:
    facility: "local0" # allows for syslog output of the same data going to file
```

### Gotchas

- If `out.file.enabled` and `out.cli.enabled` are both false (defaults) then effectively no logging happening. Expect to see error saying: `Effectively no logging enabled! The config.out.file.enabled and config.out.cli.enabled flags are both false.`
- `out.file.debug` as well as `out.cli.verbose` have no effect if the logging is disabled using either of their respective flags `out.file.enabled`, `out.cli.enabled` respectively; or even when the overall `out.enabled` is set to `false`.
- `mc.correlation-id.auto-generate` does NOT take precedence over the `mc.correlation-id.strict` checks! If `correlation-id` is marked as strict (`mc.correlation-id.strict: true`) then `correlation-id` MUST be present in the logging context
-

## Context

In order for each log entry to satisfy your  These data-points are expected to be in the context.
Each context is expected to contain the following `REQUIRED_FIELDS`:
- `type` (`mclogger.TYPE`)
- `application` (`mclogger.APPLICATION`)
- `operation` (`mclogger.OPERATION`)

Strict checking of these minimal values can be enforced by setting `mc.strict` in the `LogConfig`. 
Additionally, as mentioned in the configuration examples above, the `correlationId` can become part of `REQUIRED_FIELDS` by setting `mc.correlationId.strict` flag.


## Structured Log Output

Every log entry (output) is structured following the struct `StructuredLog`.

### Out File

The file will contain marshaled json of the `StructuredLog` struct.

Sample output:

```json language=json
{
  "ts":"2024-09-30T19:53:07.842+0100",
  "type":"Business",
  "application":"appName",
  "operation":"certs-create",
  "correlationid":"494c656b-19f0-41a7-9b42-39bc77372cee",
  "logId":"00eb5878-14ca-4cca-beb5-2e3c6f3ae2dc",
  "level":"ERROR",
  "message":"[3] Failure in generating certificate for mada.com [failure requesting certificate to the RA [Post \"http://localhost:13560/registeredAuthority\": dial tcp 127.0.0.1:13560: connect: connection refused]]",
  "attributes":{}
}
```

### Out Cli

The prompt (stdout/stderr) output by default the same output as the out file, in the same format without any debug information.

The `out.cli.friendly` flag set, will make the output opinionated towards cli user friendly output, and at the moment this will result in the below format:

```text language=text
[%sl.level%]    %sl.Timestamp%  %sl.Message%    [%sl.Correlationid%]
```

Where the correlationId only showing up if `out.cli.verbose` flag is set or the log level for the message is WARN or higher.

Sample output:

```text language=text
[ERROR] 2024-09-30T20:04:05.318+0100    [3] Failure in generating certificate for mada.com [failure requesting certificate to the RA [Post "http://localhost:13560/registeredAuthority": dial tcp 127.0.0.1:13560: connect: connection refused]]           [620468e0-6523-4157-9c5a-b77d1c5369a8]
```


### Out Syslog

If enabled using the `out.syslog.facility` flag, it will output to syslog the logging information in the same format as the out.file configuration (json format).

`out.syslog.facility` expects a valid syslog facility value. Constants exist for each facility which use the naming convention `SyslogFacility*`

Syslog output is not supported on Windows machines.

Sample output:

```text language=text
Jan 15 09:53:34 ech-10-170-135-99 snackbox[4616]: {"ts":"2025-01-15T09:53:34.717-0500","type":"Business","application":"snackbox","operation":"certs-create","correlationid":"a52b0129-9d49-4f29-acbb-3575aa4442f4","logId":"67e36893-0a7e-476c-b799-4a2772e9bd17","level":"ERROR","message":"[56] SNACKBOX_HOME has not been set, please set this environment variable on your machine to use Snackbox. ","attributes":{"err":"SNACKBOX_HOME is not defined","version":"undefined"}}
```


## Working Example

Define the configuration (recommended to read this from configuration file of your application, or even derived from user input as applicable)

```go language=go
config := mclogger.LogConfig{
    Mc: &mclogger.McConfig{
        Strict: true,
        CorrelationId: &mclogger.CorrelationIdConfig{
            Strict:       true,
            AutoGenerate: false,
        },
    },
    Out: &mclogger.OutConfig{
        Enabled: true,
        File: &mclogger.FileOutputConfig{
            Enabled: true,
            Debug:   true,
            Path:    filepath.Join("some", "location", "output.log"),
        },
        Cli: &mclogger.CliConfig{
            Enabled:  true,
            Friendly: true,
            Verbose:  false, // can be updated based on -v argument for example
        },
    },
}
```

Create the logger and set your custom configurations:

```go language=go
mcLogger := mclogger.NewMcLogger(config)
slog.SetDefault(slog.New(mcLogger))
```

Making a context:

```go language=go
ctx := context.WithValue(ctx, mclogger.CORRELATION_ID, uuid.New().String())
ctx = context.WithValue(logContext, mclogger.TYPE, mclogger.BusinessType)
ctx = context.WithValue(logContext, mclogger.APPLICATION, "AppName")
```

Calling the slog Info:

```go language=go
slog.DebugContext(ctx, "A debugging message to help you later", "key", "VALUE")
slog.InfoContext(ctx, "Information message here")
slog.WarnContext(ctx, "One Warning!", "someKey", "AnotherValue")
slog.ErrorContext(ctx, "Example Error", "reason", "justFailed")
```

Will result in:
- Output:

```text language=text
[INFO]  2024-09-30T20:13:02.687+0100    Information message here
[WARN]  2024-09-30T20:13:02.687+0100    One Warning!     | someKey=AnotherValue [05e2e7d1-76c8-41ae-afcd-ac496806cba5]
[ERROR] 2024-09-30T20:13:02.687+0100    Example Error    | reason=justFailed    [05e2e7d1-76c8-41ae-afcd-ac496806cba5]
```

- File:

```json language=json
{
    "ts":"2024-09-30T20:14:28.07+0100",
    "type":"Business",
    "application":"appName",
    "operation":"certs-create",
    "correlationid":"159928d7-9822-452d-8211-ab928487bdbb",
    "logId":"d2fa6c5b-7233-4dee-b265-a76f4f63c994",
    "level":"DEBUG",
    "message":"A debugging message to help you later",
    "attributes": {
        "Key":"key",
        "Value":{}
    }
}
```


```json language=json
{
    "ts":"2024-09-30T20:14:28.071+0100",
    "type":"Business",
    "application":"appName",
    "operation":"certs-create",
    "correlationid":"159928d7-9822-452d-8211-ab928487bdbb",
    "logId":"f54569b7-ddbf-430b-9b20-e809345e0760",
    "level":"INFO",
    "message":"Information message here",
    "attributes":{}
}
```

```json language=json
{
    "ts":"2024-09-30T20:14:28.071+0100",
    "type":"Business",
    "application":"appName",
    "operation":"certs-create",
    "correlationid":"159928d7-9822-452d-8211-ab928487bdbb",
    "logId":"e7d88048-ef4f-4cdf-b2b0-bde1786725a7",
    "level":"WARN",
    "message":"One Warning!",
    "attributes": {
        "Key":"someKey",
        "Value":{}
    }
}
```

```json language=json
{
  "ts":"2024-09-30T20:14:28.071+0100",
  "type":"Business",
  "application":"appName",
  "operation":"certs-create",
  "correlationid":"159928d7-9822-452d-8211-ab928487bdbb",
  "logId":"bfcd0097-5572-4c68-8198-9a21f84539ba",
  "level":"ERROR",
  "message":"Example Error",
  "attributes":{
    "Key":"reason",
    "Value":{}
  }
}
```

