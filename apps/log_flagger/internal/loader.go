package internal

import (
	"fmt"
	"net"
	"slices"
	"strings"

	"github.com/cloudflare/ahocorasick"
	"github.com/yl2chen/cidranger"

	"git.oxl.at/open-bot-list/log_flagger/internal/config"
)

// Custom Ranger Entry to implement cidranger.RangerEntry.
// It holds the net.IPNet and the associated flag key as a string.
type FlagRangerEntry struct {
	net.IPNet
	Key string
}

func (e *FlagRangerEntry) Network() net.IPNet {
	return e.IPNet
}

// LoadedIPTrie replaces LoadedIPLists and the two kentik/patricia tries with a single Ranger.
// It is initialized immediately, making it safe for concurrent read operations after initial loading.
var LoadedIPTrie cidranger.Ranger = cidranger.NewPCTrieRanger()

// LoadedUACategoryMatcher: For patterns from *.lst files (general categories/always match all).
var LoadedUACategoryMatcher *ahocorasick.Matcher
var UACategoryIndexToFlags = make(map[int][]string)

// LoadedUASpecificMatcher: For patterns from *.map files (specific substrings/priority logic applies).
var LoadedUASpecificMatcher *ahocorasick.Matcher

// UASpecificIndexToFlags maps the Aho-Corasick index (int) to the original flag string (string).
var UASpecificIndexToFlags = make(map[int]string)

// UASpecificACIndexToArrayIndex maps the Aho-Corasick match index (int) back to the original
// index in LoadedUserAgentMapArray (int). This is used to enforce the 'first match wins' priority rule.
var UASpecificACIndexToArrayIndex = make(map[int]int)

// Loaded*Lists hold lists of values from *.lst files.
var LoadedUserAgentLists = make(map[string][]string)
var LoadedFingerprintLists = make(map[string][]string)
var LoadedPTRLists = make(map[string][]string)
var LoadedASNLists = make(map[string][]string)

// Loaded*Map hold a map of values to their corresponding base flag from map files.
var LoadedUserAgentMapArray [][]string
var LoadedFingerprintMapArray [][]string
var LoadedPTRMapArray [][]string
var LoadedASNMapArray [][]string

// CompileUAMatcher builds the two separate Aho-Corasick matchers.
func CompileUAMatcher() {
	// --- 1. Category Matcher (General Lists, allows multiple matches) ---
	var categoryPatterns []string
	categoryPatternToFlags := make(map[string][]string)

	for listKey, subStrings := range LoadedUserAgentLists {
		baseFlags := config.USER_AGENT_LIST_FLAGS[listKey]
		for _, subString := range subStrings {
			subString = strings.ToLower(subString)
			existingFlags := categoryPatternToFlags[subString]
			existingFlags = append(existingFlags, listKey)
			for _, flag := range baseFlags {
				if !slices.Contains(existingFlags, flag) {
					existingFlags = append(existingFlags, flag)
				}
			}
			categoryPatternToFlags[subString] = existingFlags
		}
	}

	for pattern, flags := range categoryPatternToFlags {
		categoryPatterns = append(categoryPatterns, pattern)
		UACategoryIndexToFlags[len(categoryPatterns)-1] = flags
	}

	LoadedUACategoryMatcher = ahocorasick.NewStringMatcher(categoryPatterns)

	// --- 2. Specific Matcher (Map Files, implements priority/fallback logic) ---
	var specificPatterns []string
	// Temporary map to ensure uniqueness and track the source index
	specificPatternMap := make(map[string]int) // pattern -> original array index in LoadedUserAgentMapArray

	// We must iterate over the map array directly to preserve the priority order (top-to-bottom)
	for arrayIndex, valueMap := range LoadedUserAgentMapArray {
		uaSubstring := strings.ToLower(valueMap[0])
		baseFlag := valueMap[1]

		// If the pattern is already added, we skip it to preserve the first-occurrence rule.
		if _, exists := specificPatternMap[uaSubstring]; exists {
			continue
		}

		specificPatterns = append(specificPatterns, uaSubstring)
		acIndex := len(specificPatterns) - 1 // Aho-Corasick index

		UASpecificIndexToFlags[acIndex] = baseFlag          // Store only the single resulting flag
		UASpecificACIndexToArrayIndex[acIndex] = arrayIndex // Store the original map array index
		specificPatternMap[uaSubstring] = arrayIndex        // Mark as added
	}

	LoadedUASpecificMatcher = ahocorasick.NewStringMatcher(specificPatterns)

	if config.DEBUG {
		fmt.Printf("Compiled Aho-Corasick Category Matcher with %d patterns.\n", len(categoryPatterns))
		fmt.Printf("Compiled Aho-Corasick Specific Matcher with %d patterns.\n", len(specificPatterns))
	}
}

// LoadListContents processes the content of a single list file and loads it into the appropriate global map.
func LoadListContents(fileName string, content string) error {
	var isList = strings.HasSuffix(fileName, ".lst")
	var isMap = strings.HasSuffix(fileName, ".map")

	lines := strings.Split(content, "\n")
	fileName = strings.TrimSuffix(fileName, ".lst")
	fileName = strings.TrimSuffix(fileName, ".map")

	var entries []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "[") || strings.HasPrefix(line, "#") {
			continue
		}
		entries = append(entries, line)
	}

	if len(entries) == 0 {
		return nil
	}

	// IP List Logic (src_net_*.lst)
	if strings.HasPrefix(fileName, "src_net_") {
		key := strings.TrimSuffix(fileName, "_net4")
		key = strings.TrimSuffix(key, "_net6")
		// Final key (e.g., "src_net_crawler_google_common")

		if _, exists := config.IPLIST_FLAGS[key]; exists {
			for _, entry := range entries {
				ipAddr, network, err := net.ParseCIDR(entry)
				if err != nil {
					fmt.Printf("Warning: Failed to parse CIDR %s in file %s: %v\n", entry, fileName, err)
					continue
				}
				// cidranger expects the network address from ParseCIDR
				network.IP = ipAddr.Mask(network.Mask)

				rangerEntry := &FlagRangerEntry{
					IPNet: *network,
					Key:   key,
				}

				if err := LoadedIPTrie.Insert(rangerEntry); err != nil {
					fmt.Printf("Warning: Ranger Insert error for %s in file %s: %v\n", entry, fileName, err)
				}
			}
			return nil
		}
	}

	// User-Agent List Logic (http_user_agent_*_sub.lst)
	if strings.HasPrefix(fileName, "http_user_agent_") && isList && strings.HasSuffix(fileName, "_sub") {
		key := strings.TrimSuffix(fileName, "_sub")

		if _, exists := config.USER_AGENT_LIST_FLAGS[key]; exists {
			LoadedUserAgentLists[key] = append(LoadedUserAgentLists[key], entries...)
			return nil
		}
	}

	// User-Agent Map Logic (http_user_agent_*.map)
	if strings.HasPrefix(fileName, "http_user_agent_") && isMap {
		// Map files contain <base-flag><space><user-agent-substring> on each line
		for _, entry := range entries {
			parts := strings.SplitN(entry, " ", 2)
			if len(parts) != 2 {
				fmt.Printf("Warning: Skipping malformed map entry in file %s: %s\n", fileName, entry)
				continue
			}
			baseFlag := parts[0]
			uaSubstring := parts[1]

			LoadedUserAgentMapArray = append(LoadedUserAgentMapArray, []string{uaSubstring, baseFlag})
		}
		return nil
	}

	// Fingerprint List
	if strings.HasPrefix(fileName, "fingerprint_") && isList {
		if _, exists := config.FINGERPRINT_LIST_FLAGS[fileName]; exists {
			LoadedFingerprintLists[fileName] = append(LoadedFingerprintLists[fileName], entries...)
			return nil
		}
	}

	if strings.HasPrefix(fileName, "fingerprint_") && isMap {
		for _, entry := range entries {
			parts := strings.SplitN(entry, " ", 2)
			if len(parts) != 2 {
				fmt.Printf("Warning: Skipping malformed map entry in file %s: %s\n", fileName, entry)
				continue
			}
			baseFlag := parts[0]
			value := parts[1]

			LoadedFingerprintMapArray = append(LoadedFingerprintMapArray, []string{value, baseFlag})
		}
		return nil
	}

	// ASN lists
	if strings.HasPrefix(fileName, "src_asn_") && isList {
		if _, exists := config.ASN_LIST_FLAGS[fileName]; exists {
			LoadedASNLists[fileName] = append(LoadedASNLists[fileName], entries...)
			return nil

		} else {
			fmt.Printf("ASN category not found: %s\n", fileName)
		}
	}

	if strings.HasPrefix(fileName, "src_asn_") && isMap {
		for _, entry := range entries {
			parts := strings.SplitN(entry, " ", 2)
			if len(parts) != 2 {
				fmt.Printf("Warning: Skipping malformed map entry in file %s: %s\n", fileName, entry)
				continue
			}
			baseFlag := parts[0]
			value := parts[1]

			LoadedASNMapArray = append(LoadedASNMapArray, []string{value, baseFlag})
		}
		return nil
	}

	// PTR lists
	if strings.HasPrefix(fileName, "ptr_") && isList && strings.HasSuffix(fileName, "_end") {
		key := strings.TrimSuffix(fileName, "_end")

		if _, exists := config.PTR_LIST_FLAGS[key]; exists {
			LoadedPTRLists[key] = append(LoadedPTRLists[key], entries...)
			return nil
		}
	}

	if strings.HasPrefix(fileName, "ptr_") && isMap {
		for _, entry := range entries {
			parts := strings.SplitN(entry, " ", 2)
			if len(parts) != 2 {
				fmt.Printf("Warning: Skipping malformed map entry in file %s: %s\n", fileName, entry)
				continue
			}
			baseFlag := parts[0]
			value := parts[1]

			LoadedPTRMapArray = append(LoadedPTRMapArray, []string{value, baseFlag})
		}
		return nil
	}

	return fmt.Errorf("file %s does not match any known list category for loading", fileName)
}
