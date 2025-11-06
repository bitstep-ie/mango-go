//go:build !windows

package mango_logger

import "log/syslog"

type SyslogConfig struct {
	// Facility refers to the syslog facility of a given log
	Facility SyslogFacility `yaml:"facility" json:"facility"`

	// priority allows you to indicate the facility and severity of a given log
	priority syslog.Priority
}
