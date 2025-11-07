//go:build !windows

package mango_logger

import (
	"fmt"
	"log/slog"
	"log/syslog"
)

func (sl MangoLogger) handleSyslogOutput(log *StructuredLog, jsonOut []byte) error {

	var severity = syslog.LOG_EMERG
	switch log.Level {
	case slog.LevelDebug:
		severity = syslog.LOG_DEBUG
	case slog.LevelInfo:
		severity = syslog.LOG_INFO
	case slog.LevelWarn:
		severity = syslog.LOG_WARNING
	case slog.LevelError:
		severity = syslog.LOG_ERR
	default:
		fmt.Println("Record level not one of: debug, info, warn or error")
	}
	if severity == syslog.LOG_EMERG {
		return fmt.Errorf("record level not one of: debug, info, warn or error")
	}

	switch sl.Config.Out.Syslog.Facility {
	case SyslogFacilityKern:
		sl.Config.Out.Syslog.priority = syslog.LOG_KERN | severity
	case SyslogFacilityUser:
		sl.Config.Out.Syslog.priority = syslog.LOG_USER | severity
	case SyslogFacilityMail:
		sl.Config.Out.Syslog.priority = syslog.LOG_MAIL | severity
	case SyslogFacilityDaemon:
		sl.Config.Out.Syslog.priority = syslog.LOG_DAEMON | severity
	case SyslogFacilityAuth:
		sl.Config.Out.Syslog.priority = syslog.LOG_AUTH | severity
	case SyslogFacilitySyslog:
		sl.Config.Out.Syslog.priority = syslog.LOG_SYSLOG | severity
	case SyslogFacilityNews:
		sl.Config.Out.Syslog.priority = syslog.LOG_NEWS | severity
	case SyslogFacilityUucp:
		sl.Config.Out.Syslog.priority = syslog.LOG_UUCP | severity
	case SyslogFacilityCron:
		sl.Config.Out.Syslog.priority = syslog.LOG_CRON | severity
	case SyslogFacilityAuthpriv:
		sl.Config.Out.Syslog.priority = syslog.LOG_AUTHPRIV | severity
	case SyslogFacilityFtp:
		sl.Config.Out.Syslog.priority = syslog.LOG_FTP | severity
	case SyslogFacilityLocal0:
		sl.Config.Out.Syslog.priority = syslog.LOG_LOCAL0 | severity
	case SyslogFacilityLocal1:
		sl.Config.Out.Syslog.priority = syslog.LOG_LOCAL1 | severity
	case SyslogFacilityLocal2:
		sl.Config.Out.Syslog.priority = syslog.LOG_LOCAL2 | severity
	case SyslogFacilityLocal3:
		sl.Config.Out.Syslog.priority = syslog.LOG_LOCAL3 | severity
	case SyslogFacilityLocal4:
		sl.Config.Out.Syslog.priority = syslog.LOG_LOCAL4 | severity
	case SyslogFacilityLocal5:
		sl.Config.Out.Syslog.priority = syslog.LOG_LOCAL5 | severity
	case SyslogFacilityLocal6:
		sl.Config.Out.Syslog.priority = syslog.LOG_LOCAL6 | severity
	case SyslogFacilityLocal7:
		sl.Config.Out.Syslog.priority = syslog.LOG_LOCAL7 | severity
	default:
		fmt.Println("Facility level not valid")
		return fmt.Errorf("facility level not valid")
	}

	syslogWriter, err := syslog.New(syslog.Priority(sl.Config.Out.Syslog.priority), log.Application)
	if err != nil {
		fmt.Println("Error writing to syslog")
		return fmt.Errorf("error writing to syslog: %w", err)
	}

	defer func(syslogWriter *syslog.Writer) {
		err := syslogWriter.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close syslog writer %w", err)
		}
	}(syslogWriter)

	_, err = syslogWriter.Write(jsonOut)
	return err
}
