// Package net contains helpers for validating network address values.
package net

import (
	"net/netip"
	"net/url"
)

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

// IsValidIP reports whether value is a syntactically valid IPv4 or IPv6 address.
func IsValidIP(value string) bool {
	_, err := netip.ParseAddr(value)
	return err == nil
}

// IsValidURL reports whether value is a syntactically valid URL with a required scheme and host.
// The URL must have a scheme (http, https, etc.) and a host to be considered valid.
// This is stricter than net.ParseURL and rejects permissive edge cases like scheme-only or host-only URLs.
func IsValidURL(value string) bool {
	if value == "" {
		return false
	}

	u, err := url.Parse(value)
	if err != nil {
		return false
	}

	// Require a scheme
	if u.Scheme == "" {
		return false
	}

	// Require a host
	if u.Host == "" {
		return false
	}

	return true
}
