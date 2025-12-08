package internal

import (
	"net"
	"strings"
)

// LookupIP performs a lookup of the clientIP against all loaded IP Network lists.
// It returns a list of keys (base flags) for the matching lists.
func LookupIP(clientIP string) []string {
	var matchingKeys []string

	ip := net.ParseIP(clientIP)
	if ip == nil {
		return nil
	}

	for key, networks := range LoadedIPLists {
		for _, network := range networks {
			if network.Contains(ip) {
				matchingKeys = append(matchingKeys, key)
				break
			}
		}
	}

	return matchingKeys
}

// LookupUserAgent performs a case-insensitive substring match of the userAgent
// against all loaded User-Agent lists.
// It returns a list of keys (base flags) for the matching lists.
func LookupUserAgentCategory(userAgent string) string {
	lowerUserAgent := strings.ToLower(userAgent)

	for baseFlag, subStrings := range LoadedUserAgentLists {
		for _, substring := range subStrings {
			if strings.Contains(lowerUserAgent, strings.ToLower(substring)) {
				return baseFlag
			}
		}
	}

	return ""
}

func LookupUserAgent(userAgent string) string {
	lowerUserAgent := strings.ToLower(userAgent)

	for _, userAgentMap := range LoadedUserAgentMapArray {
		substring := userAgentMap[0]
		baseFlag := userAgentMap[1]
		if strings.Contains(lowerUserAgent, strings.ToLower(substring)) {
			return baseFlag
		}
	}

	return ""
}

// LookupFingerprint performs an exact string match of the fingerprint
// against all loaded Fingerprint lists.
// It returns a list of keys (base flags) for the matching lists.
func LookupFingerprint(fingerprint string) string {
	for key, fingerprints := range LoadedFingerprintLists {
		for _, fp := range fingerprints {
			if fp == fingerprint {
				return key
			}
		}
	}

	return ""
}
