//go:build windows

package mango_logger

type SyslogConfig struct {
	// Facility refers to the syslog facility of a given log
	Facility SyslogFacility `yaml:"facility" json:"facility"`
}
