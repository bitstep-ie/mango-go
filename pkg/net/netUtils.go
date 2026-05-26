// Package net contains helpers for validating network address values.
package net

import "net/netip"

// IsValidIPv4 reports whether value is a syntactically valid IPv4 address.
func IsValidIPv4(value string) bool {
	addr, err := netip.ParseAddr(value)
	return err == nil && addr.Is4()
}

// IsValidIPv6 reports whether value is a syntactically valid IPv6 address.
func IsValidIPv6(value string) bool {
	addr, err := netip.ParseAddr(value)
	return err == nil && addr.Is6()
}
