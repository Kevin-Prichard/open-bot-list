package internal

import (
	"fmt"
	util2 "git.oxl.at/open-bot-list/pkg/flagger/util"
	"net"
	"regexp"
	"slices"
	"strings"
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
	if ip == nil || LoadedIPTrie == nil {
		return nil
	}

	containingNetworks, err := LoadedIPTrie.ContainingNetworks(ip)
	if err != nil {
		return nil
	}

	for _, entry := range containingNetworks {
		if flagEntry, ok := entry.(*FlagRangerEntry); ok {
			if !slices.Contains(matchingKeys, flagEntry.Key) {
				matchingKeys = append(matchingKeys, flagEntry.Key)
			}
		}
	}

	return matchingKeys
}

// LookupUACategories performs a fast, case-insensitive substring search of the User-Agent
// against general category lists. This always returns all matches found.
func LookupUACategories(userAgent string) []string {
	if LoadedUACategoryMatcher == nil {
		return nil
	}

	var flags []string
	lowerUA := strings.ToLower(userAgent)

	matches := LoadedUACategoryMatcher.MatchThreadSafe([]byte(lowerUA))

	for _, index := range matches {
		if flagSet, exists := UACategoryIndexToFlags[index]; exists {
			for _, flag := range flagSet {
				if !slices.Contains(flags, flag) {
					flags = append(flags, flag)
				}
			}
		}
	}

	return flags
}

// LookupUASpecifics implements the priority logic for user-agent map entries:
// Only the match corresponding to the lowest index in LoadedUserAgentMapArray is returned ("first match wins").
func LookupUASpecifics(userAgent string) []string {
	if LoadedUASpecificMatcher == nil || LoadedUserAgentMapArray == nil {
		return nil
	}

	lowerUA := strings.ToLower(userAgent)
	matches := LoadedUASpecificMatcher.MatchThreadSafe([]byte(lowerUA))

	if len(matches) == 0 {
		return nil
	}

	minArrayIndex := -1
	winningACIndex := -1

	for _, acIndex := range matches {
		arrayIndex, exists := UASpecificACIndexToArrayIndex[acIndex]
		if !exists {
			continue
		}

		if minArrayIndex == -1 || arrayIndex < minArrayIndex {
			minArrayIndex = arrayIndex
			winningACIndex = acIndex
		}
	}

	if winningACIndex == -1 {
		return nil
	}

	if flag, exists := UASpecificIndexToFlags[winningACIndex]; exists {
		return []string{flag}
	}

	return nil
}

// LookupUserAgent is now the composite function that calls the two specialized lookups.
// It implements the priority: Specific Map match > General Category match.
func LookupUserAgent(userAgent string) []string {
	specificFlags := LookupUASpecifics(userAgent)
	return append(specificFlags, LookupUACategories(userAgent)...)
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
	if ptr == "" {
		return []string{}
	}
	flags := []string{}
	for baseFlag, endStrings := range LoadedPTRLists {
		for _, endString := range endStrings {
			if strings.HasSuffix(ptr, endString) {
				flags = append(flags, baseFlag)
			}
		}
	}

	for _, valueMap := range LoadedPTRMapArray {
		endString := valueMap[0]
		baseFlag := valueMap[1]
		if strings.HasSuffix(ptr, strings.ToLower(endString)) {
			flags = append(flags, baseFlag)
		}
	}

	return flags
}

// LookupPTRFlags queries the DNS-PTR-Record for the provided clientIP and returns the flags if any are configured for that PTR
func LookupPTRFlags(clientIP string) []string {
	ptr := util2.LookupPTR(clientIP)
	fmt.Printf("  > PTR: %v => %v\n", clientIP, ptr)
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

	// todo: asn map's (flags per ASN as seen in CSV-files)

	return flags
}

// LookupGeoIPASNFlags queries the ASN for the provided clientIP and returns the flags if any are configured for that ASN
func LookupGeoIPASNFlags(clientIP string) []string {
	asn, as_name := util2.LookupGeoIPASN(clientIP)
	return generateGeoIPASNFlags(asn, as_name)
}
