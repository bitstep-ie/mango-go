//go:build windows

package mango_logger

func (sl MangoLogger) handleSyslogOutput(log *StructuredLog, jsonOut []byte) error {
	return nil
}
