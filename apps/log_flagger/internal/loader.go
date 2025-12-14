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

			// Store the base flag against the substring.
			// The key is the substring for lookup, value is the base flag to return.
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
	if strings.HasPrefix(fileName, "ptr_") && isList {
		if _, exists := config.PTR_LIST_FLAGS[fileName]; exists {
			LoadedPTRLists[fileName] = append(LoadedPTRLists[fileName], entries...)
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
