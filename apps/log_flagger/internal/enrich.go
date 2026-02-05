package internal

import (
	"fmt"
	config2 "git.oxl.at/open-bot-list/pkg/flagger/config"
	"slices"
	"strings"
)

func EnrichLog(logEntry config2.LogEntry) config2.EnrichedLog {
	var flags []string
	debug := config2.DEBUG || strings.Contains(logEntry.UserAgent, config2.DEBUG_UA)

	// check if the request has any matches in the open-bot-list config
	flags = append(flags, LookupIP(logEntry.ClientIP)...)

	flags = append(flags, LookupUserAgent(logEntry.UserAgent)...)

	flags = append(
		flags,
		LookupMapsGenericSubstring(LoadedUserAgentMapArray, logEntry.UserAgent)...,
	)

	flags = append(
		flags,
		LookupListsGenericExact(LoadedFingerprintLists, logEntry.FingerprintJA4)...,
	)

	if config2.LOOKUP_PTR {
		flags = append(
			flags,
			LookupPTRFlags(logEntry.ClientIP)...,
		)
	}

	if config2.PATH_GEOIP_ASN_DB != "" {
		flags = append(
			flags,
			LookupGeoIPASNFlags(logEntry.ClientIP)...,
		)
	}

	// apply flags specific to match-files
	for _, flag := range flags {
		if value, exists := config2.IPLIST_FLAGS[flag]; exists {
			flags = append(flags, value...)
		}
		if value, exists := config2.USER_AGENT_LIST_FLAGS[flag]; exists {
			flags = append(flags, value...)
		}
		if value, exists := config2.FINGERPRINT_LIST_FLAGS[flag]; exists {
			flags = append(flags, value...)
		}
		if value, exists := config2.PTR_LIST_FLAGS[flag]; exists {
			flags = append(flags, value...)
		}
		if value, exists := config2.ASN_LIST_FLAGS[flag]; exists {
			flags = append(flags, value...)
		}
	}

	// process ruleset and apply flags if rule matches (only one match)
	for ruleId, rule := range config2.FLAGGING_RULESET {
		matching := true
		for _, requireFlag := range rule[0] {
			if !slices.Contains(flags, requireFlag) {
				matching = false
				break
			}
		}
		if !matching {
			continue
		}
		if debug {
			fmt.Printf("MATCHING RULE: #%d %v\n", ruleId, rule)
		}
		flags = append(flags, rule[1]...)
		break
	}

	// custom fallback flags (might be migrated to open-bot-list config later on)
	if logEntry.UserAgent == "" {
		flags = append(flags, "bot", "bot_unknown")

	} else if !slices.Contains(flags, "crawler_verified") {
		for _, uaSub := range config2.FALLBACK_SPOOFED_UA_SUB {
			if strings.Contains(logEntry.UserAgent, uaSub) {
				flags = append(flags, "crawler_spoofed")
				break
			}
		}
	}

	flaggedAsBot := slices.Contains(flags, "bot")
	if !slices.Contains(flags, "bot") {
		for _, flag := range flags {
			if strings.HasPrefix(flag, "bot_") || strings.HasPrefix(flag, "crawler_") {
				flags = append(flags, "bot")
				flaggedAsBot = true
				break
			}
		}
	}
	if !flaggedAsBot {
		uaLower := strings.ToLower(logEntry.UserAgent)
		if strings.Contains(uaLower, "bot") || strings.Contains(uaLower, "crawler") || strings.Contains(uaLower, "spider") {
			flags = append(flags, "bot", "bot_unknown")
		}
	}

	if debug {
		fmt.Printf("  > Flags: %v\n  > UA: %v\n", flags, logEntry.UserAgent)
	}
	/*
		if len(flags) > 0 {

				fmt.Printf(
					"Flags: \"%v\"\n  -> (IP='%s', UA='%s', FP='%s')\n\n",
					flagsStr, logEntry.ClientIP, logEntry.UserAgent, logEntry.FingerprintJA4,
				)
		}
	*/

	return PrepareOutput(logEntry, SerializeFlags(flags))
}
