package internal

import (
	"net"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/yl2chen/cidranger"

	"git.oxl.at/open-bot-list/log_flagger/internal/config"
)

var REGEX_COLLAPSE_HYPHENS_TEST = regexp.MustCompile(`-+`)

// Helper to reset and populate the global trie for tests
func setupTestIPTrie() {
	LoadedIPTrie = cidranger.NewPCTrieRanger()

	testData := map[string][]string{
		"src_net_vpn_tor":               {"10.0.0.0/8"},
		"src_net_crawler_google":        {"192.168.1.0/24"},
		"src_net_crawler_special":       {"192.168.1.10/32", "2001:db8:a::/64"},
		"src_net_crawler_google_common": {"66.249.76.0/23"},
	}

	for key, cidrs := range testData {
		for _, cidrStr := range cidrs {
			ipAddr, network, err := net.ParseCIDR(cidrStr)
			if err != nil {
				panic(err)
			}
			network.IP = ipAddr.Mask(network.Mask)

			rangerEntry := &FlagRangerEntry{
				IPNet: *network,
				Key:   key,
			}

			if err := LoadedIPTrie.Insert(rangerEntry); err != nil {
				panic(err)
			}
		}
	}
}

func TestLookupIP(t *testing.T) {
	setupTestIPTrie()

	tests := []struct {
		name     string
		clientIP string
		wantKeys []string
	}{
		{
			name:     "IPv4 match on /8",
			clientIP: "10.1.2.3",
			wantKeys: []string{"src_net_vpn_tor"},
		},
		{
			name:     "IPv4 match on /24",
			clientIP: "192.168.1.50",
			wantKeys: []string{"src_net_crawler_google"},
		},
		{
			name:     "IPv4 exact match on /32 and partial on /24 (Multiple matches)",
			clientIP: "192.168.1.10",
			// Should match "192.168.1.10/32" (special) and "192.168.1.0/24" (google)
			wantKeys: []string{"src_net_crawler_google", "src_net_crawler_special"},
		},
		{
			name:     "IPv6 match",
			clientIP: "2001:db8:a::5",
			wantKeys: []string{"src_net_crawler_special"},
		},
		{
			name:     "Googlebot net",
			clientIP: "66.249.77.224",
			wantKeys: []string{"src_net_crawler_google_common"},
		},
		{
			name:     "No match",
			clientIP: "8.8.8.8",
			wantKeys: nil,
		},
		{
			name:     "Invalid IP",
			clientIP: "not-an-ip",
			wantKeys: nil,
		},
		{
			name:     "Empty IP",
			clientIP: "",
			wantKeys: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys := LookupIP(tt.clientIP)

			if gotKeys != nil {
				sort.Strings(gotKeys)
			}
			if tt.wantKeys != nil {
				sort.Strings(tt.wantKeys)
			}

			if !reflect.DeepEqual(gotKeys, tt.wantKeys) {
				t.Errorf("LookupIP() got = %v, want %v", gotKeys, tt.wantKeys)
			}
		})
	}
}

func TestLookupUserAgent(t *testing.T) {
	originalLoadedUserAgentLists := LoadedUserAgentLists
	originalLoadedUserAgentMapArray := LoadedUserAgentMapArray
	originalLoadedUACategoryMatcher := LoadedUACategoryMatcher
	originalUACategoryIndexToFlags := UACategoryIndexToFlags
	originalLoadedUASpecificMatcher := LoadedUASpecificMatcher
	originalUASpecificIndexToFlags := UASpecificIndexToFlags
	originalUASpecificACIndexToArrayIndex := UASpecificACIndexToArrayIndex

	t.Cleanup(func() {
		LoadedUserAgentLists = originalLoadedUserAgentLists
		LoadedUserAgentMapArray = originalLoadedUserAgentMapArray
		LoadedUACategoryMatcher = originalLoadedUACategoryMatcher
		UACategoryIndexToFlags = originalUACategoryIndexToFlags
		LoadedUASpecificMatcher = originalLoadedUASpecificMatcher
		UASpecificIndexToFlags = originalUASpecificIndexToFlags
		UASpecificACIndexToArrayIndex = originalUASpecificACIndexToArrayIndex
	})

	// Clear and set up mock data for testing the compilation logic
	LoadedUserAgentLists = map[string][]string{
		"http_user_agent_crawler":    {"spider", "crawler"},
		"http_user_agent_monitoring": {"uptime"},
	}

	// The order here defines priority: Index 0 is highest priority.
	// NOTE: The pattern 'bot' is included as a low-priority, general match.
	LoadedUserAgentMapArray = [][]string{
		{"crawler", "http_user_agent_crawler_general"},       // Array Index 0 (Highest Priority)
		{"my-specific-bot", "http_user_agent_crawler_niche"}, // Array Index 1 (Medium Priority)
		{"BOT", "http_user_agent_fallback_bot"},              // Array Index 2 (Lowest Priority/Fallback)
	}

	// Set up mock config required by CompileUAMatcher to associate list names with flags
	originalUAFlags := config.USER_AGENT_LIST_FLAGS
	config.USER_AGENT_LIST_FLAGS = map[string][]string{
		"http_user_agent_crawler":    {"bot", "bot_crawler"},
		"http_user_agent_monitoring": {"bot", "bot_monitoring"},
	}
	t.Cleanup(func() {
		config.USER_AGENT_LIST_FLAGS = originalUAFlags
	})

	CompileUAMatcher()

	tests := []struct {
		name      string
		userAgent string
		wantFlags []string
	}{
		{
			name:      "Case 1: Only Fallback Matches (Lowest Priority Wins)",
			userAgent: "I am a simple Bot",
			// Matches only "bot" (Array Index 2).
			// Result: Flag from Index 2.
			wantFlags: []string{"http_user_agent_fallback_bot"},
		},
		{
			name:      "Case 2: Priority Clash - Index 0 (crawler) wins over Index 2 (bot)",
			userAgent: "Fast crawler v1.0",
			// Matches "crawler" (Index 0) and "bot" (Index 2).
			// Matches list "crawler"
			// Result: Flag from Index 0 (lowest array index wins).
			wantFlags: []string{"bot", "bot_crawler", "http_user_agent_crawler", "http_user_agent_crawler_general"},
		},
		{
			name:      "Case 3: Priority Clash - Index 1 (niche) wins over Index 2 (bot)",
			userAgent: "Mozilla/5.0 (compatible; my-specific-bot)",
			// Matches "my-specific-bot" (Index 1) and "bot" (Index 2).
			// Result: Flag from Index 1 (lowest array index wins).
			wantFlags: []string{"http_user_agent_crawler_niche"},
		},
		{
			name:      "Case 4: Category Match Only (No Specific Match)",
			userAgent: "myspider checks uptime",
			// Matches "spider" and "uptime" in the Category Matcher.
			// Result: Fall back to Category results (all found).
			wantFlags: []string{"bot", "bot_crawler", "bot_monitoring", "http_user_agent_crawler", "http_user_agent_monitoring"},
		},
		{
			name:      "Case 5: No Match",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			wantFlags: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags := LookupUserAgent(tt.userAgent)

			if len(gotFlags) > 0 && tt.wantFlags != nil {
				sort.Strings(gotFlags)
				sort.Strings(tt.wantFlags)
			}

			if !reflect.DeepEqual(gotFlags, tt.wantFlags) {
				t.Errorf("LookupUserAgent(%q) got = %v, want %v", tt.userAgent, gotFlags, tt.wantFlags)
			}
		})
	}
}

func TestLookupListsGenericSubstring(t *testing.T) {
	originalLoadedUserAgentLists := LoadedUserAgentLists

	LoadedUserAgentLists = map[string][]string{
		"http_user_agent_crawler":    {"spider", "bingbot", "crawl"},
		"http_user_agent_monitoring": {"uptime"},
		"http_user_agent_script":     {"curl"},
		"http_user_agent_ai":         {"gptbot"},
		"http_user_agent_multiple":   {"bot"}, // for multiple matches
	}

	t.Cleanup(func() {
		LoadedUserAgentLists = originalLoadedUserAgentLists
	})

	tests := []struct {
		name      string
		userAgent string
		wantFlags []string
	}{
		{
			name:      "Case-insensitive match (lowercase)",
			userAgent: "mozilla/5.0 (compatible; bingbot/2.0)",
			// "bingbot" matches "bingbot" (crawler) AND "bot" (multiple)
			wantFlags: []string{"http_user_agent_crawler", "http_user_agent_multiple"},
		},
		{
			name:      "Multiple matches",
			userAgent: "I am a GPTBot bot",
			wantFlags: []string{"http_user_agent_ai", "http_user_agent_multiple"},
		},
		{
			name:      "No match",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			wantFlags: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags := LookupListsGenericSubstring(LoadedUserAgentLists, tt.userAgent)

			if len(gotFlags) > 0 && tt.wantFlags != nil {
				sort.Strings(gotFlags)
				sort.Strings(tt.wantFlags)
			}

			if !reflect.DeepEqual(gotFlags, tt.wantFlags) {
				t.Errorf("LookupListsGenericSubstring(%q) got = %v, want %v", tt.userAgent, gotFlags, tt.wantFlags)
			}
		})
	}
}

func TestLookupMapsGenericSubstring(t *testing.T) {
	originalLoadedUserAgentMapArray := LoadedUserAgentMapArray

	LoadedUserAgentMapArray = [][]string{
		{
			"Bingbot",
			"http_user_agent_crawler_microsoft_bing",
		},
		{
			"Googlebot/",
			"http_user_agent_crawler_google_search",
		},
		{
			"Google",
			"http_user_agent_crawler_google_fallback", // Match for both, should return both
		},
	}

	t.Cleanup(func() {
		LoadedUserAgentMapArray = originalLoadedUserAgentMapArray
	})

	tests := []struct {
		name      string
		userAgent string
		wantFlags []string
	}{
		{
			name:      "Bingbot (single match)",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0)",
			wantFlags: []string{"http_user_agent_crawler_microsoft_bing"},
		},
		{
			name:      "Googlebot (multiple matches)",
			userAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			wantFlags: []string{"http_user_agent_crawler_google_search", "http_user_agent_crawler_google_fallback"},
		},
		{
			name:      "No match",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			wantFlags: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags := LookupMapsGenericSubstring(LoadedUserAgentMapArray, tt.userAgent)

			if len(gotFlags) > 0 && tt.wantFlags != nil {
				sort.Strings(gotFlags)
				sort.Strings(tt.wantFlags)
			}

			if !reflect.DeepEqual(gotFlags, tt.wantFlags) {
				t.Errorf("LookupMapsGenericSubstring(%q) got = %v, want %v", tt.userAgent, gotFlags, tt.wantFlags)
			}
		})
	}
}

func TestLookupFingerprint(t *testing.T) { // Tests LookupListsGenericExact
	originalLoadedFingerprintLists := LoadedFingerprintLists

	LoadedFingerprintLists = map[string][]string{
		"fingerprint_script":          {"33,65281-10,34,16,11,43", "99,10,23,5,15,6"},
		"fingerprint_scanner_tls_ja4": {"t13d1811h2_e8a523a41297_5894756feeaa"},
		"fingerprint_multiple":        {"t13d1811h2_e8a523a41297_5894756feeaa"},
	}

	t.Cleanup(func() {
		LoadedFingerprintLists = originalLoadedFingerprintLists
	})

	tests := []struct {
		name        string
		toMatch     string
		loadedLists map[string][]string
		wantFlags   []string
	}{
		{
			name:        "Exact match on script list",
			toMatch:     "99,10,23,5,15,6",
			loadedLists: LoadedFingerprintLists,
			wantFlags:   []string{"fingerprint_script"},
		},
		{
			name:        "Multiple exact matches",
			toMatch:     "t13d1811h2_e8a523a41297_5894756feeaa",
			loadedLists: LoadedFingerprintLists,
			wantFlags:   []string{"fingerprint_multiple", "fingerprint_scanner_tls_ja4"},
		},
		{
			name:        "No match",
			toMatch:     "t13d1811h2_e8a523a41297_5894756fee65",
			loadedLists: LoadedFingerprintLists,
			wantFlags:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags := LookupListsGenericExact(tt.loadedLists, tt.toMatch)

			if len(gotFlags) > 0 && tt.wantFlags != nil {
				sort.Strings(gotFlags)
				sort.Strings(tt.wantFlags)
			}

			if !reflect.DeepEqual(gotFlags, tt.wantFlags) {
				t.Errorf("LookupListsGenericExact(%q) got = %v, want %v", tt.toMatch, gotFlags, tt.wantFlags)
			}
		})
	}
}

func TestGeneratePTRFlags(t *testing.T) {
	// NOTE: This test verifies the flag generation logic within LookupPTRFlags,
	// assuming the dependency (util.LookupPTR) returns a known PTR record.
	originalLoadedPTRLists := LoadedPTRLists
	originalLookupPtr := config.LOOKUP_PTR
	config.LOOKUP_PTR = true

	LoadedPTRLists = map[string][]string{
		"ptr_crawler":  {"google.com", "yandex.net"},
		"ptr_hoster":   {"digitalocean.com"},
		"ptr_multiple": {"google.com"},
	}

	t.Cleanup(func() {
		LoadedPTRLists = originalLoadedPTRLists
		config.LOOKUP_PTR = originalLookupPtr
	})

	tests := []struct {
		name      string
		ptr       string // The simulated result of util.LookupPTR(clientIP)
		wantFlags []string
	}{
		{
			name:      "Google PTR match (multiple matches)",
			ptr:       "some-bot.google.com",
			wantFlags: []string{"ptr_crawler", "ptr_multiple"},
		},
		{
			name:      "Yandex PTR match (exact suffix)",
			ptr:       "yandex.net",
			wantFlags: []string{"ptr_crawler"},
		},
		{
			name:      "No match",
			ptr:       "example.org",
			wantFlags: []string{},
		},
		{
			name:      "Empty PTR",
			ptr:       "",
			wantFlags: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic of LookupPTRFlags using the predefined PTR result.
			flags := generatePTRFlags(tt.ptr)

			if len(flags) > 0 && tt.wantFlags != nil {
				sort.Strings(flags)
				sort.Strings(tt.wantFlags)
			}

			if !reflect.DeepEqual(flags, tt.wantFlags) {
				t.Errorf("generatePTRFlags() got = %v, want %v (Simulated PTR: %q)", flags, tt.wantFlags, tt.ptr)
			}
		})
	}
}

func TestGenerateGeoIPASNFlags(t *testing.T) {
	// NOTE: This test verifies the ASN flagging logic and AS-Name sanitization,
	// assuming the dependency (util.LookupGeoIPASN) returns known values.
	originalLoadedASNLists := LoadedASNLists

	LoadedASNLists = map[string][]string{
		"src_asn_cdn":     {"15169"}, // Google
		"src_asn_cloud":   {"16509"}, // Amazon
		"src_asn_crawler": {"16509"}, // Amazon (multiple match)
	}

	t.Cleanup(func() {
		LoadedASNLists = originalLoadedASNLists
	})

	tests := []struct {
		name      string
		asn       string
		asName    string
		wantFlags []string
	}{
		{
			name:      "Google ASN match",
			asn:       "15169",
			asName:    "Google LLC",
			wantFlags: []string{"src_asn_cdn", "src_asn_15169", "src_as_name_Google-LLC"},
		},
		{
			name:   "Amazon ASN match (multiple lists) and complex name sanitization (Fix for '---')",
			asn:    "16509",
			asName: "AMAZON-02 - Amazon.com, Inc. [Web]",
			// Step 1 (REGEX_AS_NAME_SAFE): "AMAZON-02--Amazon.com--Inc.-[Web]"
			// Step 2 (REGEX_COLLAPSE_HYPHENS): "AMAZON-02-Amazon.com-Inc.-[Web]"
			wantFlags: []string{"src_asn_cloud", "src_asn_crawler", "src_asn_16509", "src_as_name_AMAZON-02-Amazon.com-Inc.-[Web]"},
		},
		{
			name:      "No ASN match (generic ASN)",
			asn:       "64512",
			asName:    "TEST-AS-NAME",
			wantFlags: []string{"src_asn_64512", "src_as_name_TEST-AS-NAME"},
		},
		{
			name:      "Empty ASN",
			asn:       "",
			asName:    "",
			wantFlags: []string{},
		},
		{
			name:      "No AS Name",
			asn:       "15169",
			asName:    "",
			wantFlags: []string{"src_asn_cdn", "src_asn_15169"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic of LookupGeoIPASNFlags using the predefined ASN results.
			flags := generateGeoIPASNFlags(tt.asn, tt.asName)

			if len(flags) > 0 && tt.wantFlags != nil {
				sort.Strings(flags)
				sort.Strings(tt.wantFlags)
			}

			if !reflect.DeepEqual(flags, tt.wantFlags) {
				t.Errorf("generateGeoIPASNFlags() got = %v, want %v (Simulated ASN/Name: %q/%q)", flags, tt.wantFlags, tt.asn, tt.asName)
			}
		})
	}
}
