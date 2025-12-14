package internal

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"git.oxl.at/open-bot-list/log_flagger/internal/util"
)

var (
	REGEX_AS_NAME_SAFE     = regexp.MustCompile(`[^a-zA-Z0-9\._\-\(\)\[\]]`)
	REGEX_COLLAPSE_HYPHENS = regexp.MustCompile(`-+`)
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

// LookupListsGenericSubstring performs a case-insensitive substring match of the value against all supplied lists.
// It returns a list of keys (base flags) for the matching lists.
func LookupListsGenericSubstring(loadedLists map[string][]string, toMatch string) []string {
	var flags []string
	lowerToMatch := strings.ToLower(toMatch)

	for baseFlag, subStrings := range loadedLists {
		for _, substring := range subStrings {
			if strings.Contains(lowerToMatch, strings.ToLower(substring)) {
				flags = append(flags, baseFlag)
			}
		}
	}

	return flags
}

func LookupMapsGenericSubstring(loadedMaps [][]string, toMatch string) []string {
	var flags []string
	lowerToMatch := strings.ToLower(toMatch)

	for _, valueMap := range loadedMaps {
		substring := valueMap[0]
		baseFlag := valueMap[1]
		if strings.Contains(lowerToMatch, strings.ToLower(substring)) {
			flags = append(flags, baseFlag)
		}
	}

	return flags
}

// LookupListsGenericExact performs an exact string match of the value against all supplied lists.
// It returns a list of keys (base flags) for the matching lists.
func LookupListsGenericExact(loadedLists map[string][]string, toMatch string) []string {
	var flags []string
	for key, values := range loadedLists {
		for _, value := range values {
			if value == toMatch {
				flags = append(flags, key)
			}
		}
	}

	return flags
}

func generatePTRFlags(ptr string) []string {
	var flags []string
	for baseFlag, endStrings := range LoadedPTRLists {
		for _, endString := range endStrings {
			if strings.HasSuffix(ptr, endString) {
				flags = append(flags, baseFlag)
			}
		}
	}

	return flags
}

func LookupPTRFlags(clientIP string) []string {
	ptr := util.LookupPTR(clientIP)
	return generatePTRFlags(ptr)
}

func generateGeoIPASNFlags(asn string, as_name string) []string {
	if asn == "" {
		return []string{}
	}
	flags := LookupListsGenericExact(
		LoadedASNLists,
		asn,
	)
	flags = append(flags, fmt.Sprintf("src_asn_%s", asn))
	if as_name != "" {
		as_name_safe := REGEX_AS_NAME_SAFE.ReplaceAllString(as_name, "-")
		as_name_safe = REGEX_COLLAPSE_HYPHENS.ReplaceAllString(as_name_safe, "-")
		flags = append(flags, fmt.Sprintf("src_as_name_%s", as_name_safe))
	}
	return flags
}

func LookupGeoIPASNFlags(clientIP string) []string {
	// todo: add as-name to logs
	asn, as_name := util.LookupGeoIPASN(clientIP)
	return generateGeoIPASNFlags(asn, as_name)
}
