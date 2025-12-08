package internal

import (
	"fmt"
	"net"
	"strings"

	"git.oxl.at/open-bot-list/log_flagger/internal/config"
)

// LoadedIPLists holds lists of pre-parsed IP networks (CIDRs).
// The keys match those in IPLIST_FLAGS (e.g., "src_net_crawler_google_common").
var LoadedIPLists = make(map[string][]net.IPNet)

// LoadedUserAgentLists holds lists of User-Agent substring strings from *.lst files.
// The keys match those in USER_AGENT_LIST_FLAGS.
var LoadedUserAgentLists = make(map[string][]string)

// LoadedUserAgentMap holds a map of User-Agent substrings to their corresponding base flag from map files.
// The key is the substring, and the value is the base flag string.
var LoadedUserAgentMapArray [][]string

// LoadedFingerprintLists holds lists of exact fingerprint strings (JA4).
// The keys match those in FINGERPRINT_LIST_FLAGS.
var LoadedFingerprintLists = make(map[string][]string)

// LoadListContents processes the content of a single list file and loads it into the appropriate global map.
func LoadListContents(fileName string, content string) error {
	lines := strings.Split(content, "\n")
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

	// --- IP List Logic (src_net_*.lst) ---
	if strings.HasPrefix(fileName, "src_net_") {
		key := strings.ReplaceAll(fileName, "_net4.lst", "")
		key = strings.ReplaceAll(key, "_net6.lst", "")
		key = strings.TrimSuffix(key, ".lst")
		// Final key (e.g., "src_net_crawler_google_common")

		if _, exists := config.IPLIST_FLAGS[key]; exists {
			for _, entry := range entries {
				_, network, err := net.ParseCIDR(entry)
				if err != nil {
					fmt.Printf("Warning: Failed to parse CIDR %s in file %s: %v\n", entry, fileName, err)
					continue
				}
				LoadedIPLists[key] = append(LoadedIPLists[key], *network)
			}
			return nil
		}
	}

	// --- User-Agent List Logic (http_user_agent_*_sub.lst) ---
	if strings.HasPrefix(fileName, "http_user_agent_") && strings.HasSuffix(fileName, "_sub.lst") {
		key := strings.ReplaceAll(fileName, "_sub.lst", "")
		key = strings.TrimSuffix(key, ".lst")

		if _, exists := config.USER_AGENT_LIST_FLAGS[key]; exists {
			LoadedUserAgentLists[key] = append(LoadedUserAgentLists[key], entries...)
			return nil
		}
	}

	// --- User-Agent Map Logic (http_user_agent_*.map) ---
	if strings.HasPrefix(fileName, "http_user_agent_") && strings.HasSuffix(fileName, ".map") {
		// Map files contain <base-flag><space><user-agent-substring> on each line
		for _, entry := range entries {
			parts := strings.SplitN(entry, " ", 2)
			if len(parts) != 2 {
				fmt.Printf("Warning: Skipping malformed map entry in file %s: %s\n", fileName, entry)
				continue
			}
			baseFlag := parts[0]
			uaSubstring := parts[1]

			// Store the base flag against the substring.
			// The key is the substring for lookup, value is the base flag to return.
			LoadedUserAgentMapArray = append(LoadedUserAgentMapArray, []string{uaSubstring, baseFlag})
			// LoadedUserAgentMap[uaSubstring] = baseFlag
		}
		return nil
	}

	// --- Fingerprint List Logic (fingerprint_*.lst) ---
	if strings.HasPrefix(fileName, "fingerprint_") && strings.HasSuffix(fileName, ".lst") {
		key := strings.TrimSuffix(fileName, ".lst")

		if _, exists := config.FINGERPRINT_LIST_FLAGS[key]; exists {
			LoadedFingerprintLists[key] = append(LoadedFingerprintLists[key], entries...)
			return nil
		}
	}

	return fmt.Errorf("file %s does not match any known list category for loading", fileName)
}
