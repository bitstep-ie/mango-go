package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidIPv4(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "valid localhost", value: "127.0.0.1", want: true},
		{name: "valid private ip", value: "192.168.1.10", want: true},
		{name: "valid zero", value: "0.0.0.0", want: true},
		{name: "valid broadcast", value: "255.255.255.255", want: true},
		{name: "empty", value: "", want: false},
		{name: "ipv6", value: "2001:db8::1", want: false},
		{name: "out of range octet", value: "256.1.2.3", want: false},
		{name: "too few octets", value: "10.0.1", want: false},
		{name: "too many octets", value: "10.0.1.2.3", want: false},
		{name: "contains spaces", value: " 192.168.1.1 ", want: false},
		{name: "contains letters", value: "a.b.c.d", want: false},
		{name: "ipv4 mapped ipv6", value: "::ffff:192.168.1.1", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidIPv4(tt.value))
		})
	}
}

func TestIsValidIPv6(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "valid loopback", value: "::1", want: true},
		{name: "valid compressed", value: "2001:db8::1", want: true},
		{name: "valid full", value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", want: true},
		{name: "valid ipv4 mapped", value: "::ffff:192.168.1.1", want: true},
		{name: "empty", value: "", want: false},
		{name: "ipv4", value: "192.168.1.10", want: false},
		{name: "too many groups", value: "2001:db8:1:2:3:4:5:6:7", want: false},
		{name: "invalid hex", value: "2001:db8::zzzz", want: false},
		{name: "contains spaces", value: " 2001:db8::1 ", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidIPv6(tt.value))
		})
	}
}

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		// IPv4 cases
		{name: "valid ipv4 localhost", value: "127.0.0.1", want: true},
		{name: "valid ipv4 private", value: "192.168.1.10", want: true},
		{name: "valid ipv4 zero", value: "0.0.0.0", want: true},
		{name: "valid ipv4 broadcast", value: "255.255.255.255", want: true},
		// IPv6 cases
		{name: "valid ipv6 loopback", value: "::1", want: true},
		{name: "valid ipv6 compressed", value: "2001:db8::1", want: true},
		{name: "valid ipv6 full", value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", want: true},
		{name: "valid ipv6 mapped ipv4", value: "::ffff:192.168.1.1", want: true},
		// Invalid cases
		{name: "empty", value: "", want: false},
		{name: "invalid ipv4 octet", value: "256.1.2.3", want: false},
		{name: "invalid ipv4 count", value: "10.0.1", want: false},
		{name: "invalid ipv6 hex", value: "2001:db8::zzzz", want: false},
		{name: "contains whitespace", value: " 192.168.1.1 ", want: false},
		{name: "invalid format", value: "not-an-ip", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidIP(tt.value))
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		// Valid URLs
		{name: "valid http url", value: "http://example.com", want: true},
		{name: "valid https url", value: "https://example.com", want: true},
		{name: "valid url with path", value: "https://example.com/path/to/resource", want: true},
		{name: "valid url with query", value: "https://example.com?key=value", want: true},
		{name: "valid url with fragment", value: "https://example.com#section", want: true},
		{name: "valid url with port", value: "https://example.com:8080", want: true},
		{name: "valid url with user info", value: "https://user:pass@example.com", want: true},
		{name: "valid ftp url", value: "ftp://files.example.com", want: true},
		{name: "valid url with subdomain", value: "https://api.v2.example.com/v1/users", want: true},
		{name: "valid localhost", value: "http://localhost:3000", want: true},
		{name: "valid ip url", value: "http://192.168.1.1:8080", want: true},
		{name: "valid trailing slash", value: "https://example.com/", want: true},
		// Invalid URLs
		{name: "empty string", value: "", want: false},
		{name: "no scheme", value: "example.com", want: false},
		{name: "scheme only", value: "https://", want: false},
		{name: "host only", value: "example.com", want: false},
		{name: "whitespace in url", value: "https://exam ple.com", want: false},
		{name: "invalid scheme format", value: "ht!tp://example.com", want: false},
		{name: "malformed url", value: "://example.com", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidURL(tt.value))
		})
	}
}
