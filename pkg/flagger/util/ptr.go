package util

import (
	"context"
	"git.oxl.at/open-bot-list/pkg/flagger/config"
	"net"
	"strings"
)

var PTRCache = make(map[string]string)

// LookupPTR performs a Reverse DNS (PTR) lookup for the given client IP.
// It returns the first resolved hostname or an empty string if the lookup fails or times out after 1500ms.
func LookupPTR(clientIP string) string {
	// use cached value if exists
	if value, exists := PTRCache[clientIP]; exists {
		return value
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.DNS_RESOLVE_TIMEOUT)
	defer cancel()

	r := net.Resolver{}

	names, err := r.LookupAddr(ctx, clientIP)

	if err != nil {
		return ""
	}

	if len(names) > 0 {
		name := names[0]

		if len(name) > 0 && name[len(name)-1] == '.' {
			name = name[:len(name)-1]
		}
		ptr := strings.ToLower(name)
		PTRCache[clientIP] = ptr
		return ptr
	}

	return ""
}
