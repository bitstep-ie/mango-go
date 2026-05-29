# `pkg/net`

Helpers for validating network address strings. The package currently exposes strict IPv4 and IPv6 checkers for input validation in config, API payloads, and CLI flags.

## Quick Start

```go
import mangonet "github.com/bitstep-ie/mango-go/pkg/net"

func allowHosts(ipv4 string, ipv6 string) bool {
	return mangonet.IsValidIPv4(ipv4) && mangonet.IsValidIPv6(ipv6)
}
```

## API Cheatsheet

| Function | Purpose |
| --- | --- |
| `IsValidIP(value string) bool` | returns `true` only for syntactically valid IPv4 or IPv6 addresses |
| `IsValidIPv4(value string) bool` | returns `true` only for syntactically valid IPv4 addresses |
| `IsValidIPv6(value string) bool` | returns `true` only for syntactically valid IPv6 addresses |
| `IsValidURL(value string) bool` | returns `true` only for syntactically valid URLs with scheme and host |

## Examples

### Accept any valid IP (v4 or v6)

```go
mangonet.IsValidIP("127.0.0.1")       // true
mangonet.IsValidIP("::1")             // true
mangonet.IsValidIP("invalid")         // false
```

### Accept valid IPv4

```go
mangonet.IsValidIPv4("127.0.0.1")     // true
mangonet.IsValidIPv4("192.168.1.10")  // true
mangonet.IsValidIPv4("255.255.255.255") // true
```

### Reject invalid input

```go
mangonet.IsValidIPv4("2001:db8::1")   // false (IPv6)
mangonet.IsValidIPv4("256.1.2.3")     // false (octet out of range)
mangonet.IsValidIPv4("10.0.1")        // false (not 4 octets)
mangonet.IsValidIPv4(" 192.168.1.1 ") // false (whitespace)
```

### Accept valid IPv6

```go
mangonet.IsValidIPv6("::1")                                   // true
mangonet.IsValidIPv6("2001:db8::1")                           // true
mangonet.IsValidIPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334") // true
mangonet.IsValidIPv6("::ffff:192.168.1.1")                    // true (IPv4-mapped IPv6)
```

### Reject invalid IPv6 input

```go
mangonet.IsValidIPv6("192.168.1.1")        // false (IPv4)
mangonet.IsValidIPv6("2001:db8::zzzz")     // false (invalid hex)
mangonet.IsValidIPv6("2001:db8:1:2:3:4:5:6:7") // false (too many groups)
mangonet.IsValidIPv6(" 2001:db8::1 ")      // false (whitespace)
```

### Accept valid URLs

```go
mangonet.IsValidURL("http://example.com")              // true
mangonet.IsValidURL("https://api.example.com/v1/users") // true
mangonet.IsValidURL("https://example.com:8080/path?key=value") // true
mangonet.IsValidURL("ftp://files.example.com")         // true
```

### Reject invalid URLs

```go
mangonet.IsValidURL("example.com")          // false (no scheme)
mangonet.IsValidURL("https://")             // false (no host)
mangonet.IsValidURL("")                     // false (empty)
mangonet.IsValidURL("not a url")            // false (invalid format)
```

## Behavior Notes

- Validation is syntactic only; it does not check host reachability.
- Reserved/private/public ranges are all treated as valid if the address format is correct.
- The helper does not trim input. Normalize user input before validation if needed.

### URL Validation Specifics

- `IsValidURL` requires both a **scheme** (http, https, ftp, etc.) and a **host** to be present.
- This is stricter than `net.ParseURL`, rejecting permissive edge cases like scheme-only or host-only strings.
- Supports any valid scheme (http, https, ftp, ws, custom protocols, etc.).
- Accepts URLs with ports, paths, query strings, fragments, and user info.
- Does not validate the scheme against a whitelist; any non-empty valid scheme is accepted (use additional logic if you need to restrict schemes).

