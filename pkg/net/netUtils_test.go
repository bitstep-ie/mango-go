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
