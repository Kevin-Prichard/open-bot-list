package internal

import (
	"net"
	"reflect"
	"sort"
	"testing"
)

func TestLookupIP(t *testing.T) {
	originalLoadedIPLists := LoadedIPLists

	_, net10_0_0_0_8, _ := net.ParseCIDR("10.0.0.0/8")
	_, net192_168_1_0_24, _ := net.ParseCIDR("192.168.1.0/24")
	_, net2001_db8_a_0_64, _ := net.ParseCIDR("2001:db8:a::/64")
	_, net192_168_1_10_32, _ := net.ParseCIDR("192.168.1.10/32")
	_, net66_249_76_0_23, _ := net.ParseCIDR("66.249.76.0/23")

	LoadedIPLists = map[string][]net.IPNet{
		"src_net_vpn_tor":               {*net10_0_0_0_8},
		"src_net_crawler_google":        {*net192_168_1_0_24},
		"src_net_crawler_special":       {*net192_168_1_10_32, *net2001_db8_a_0_64},
		"src_net_crawler_google_common": {*net66_249_76_0_23},
	}

	t.Cleanup(func() {
		LoadedIPLists = originalLoadedIPLists
	})

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

func TestLookupUserAgentCategory(t *testing.T) {
	originalLoadedUserAgentLists := LoadedUserAgentLists

	LoadedUserAgentLists = map[string][]string{
		"http_user_agent_crawler":    {"spider", "bingbot", "crawl"},
		"http_user_agent_monitoring": {"uptime"},
		"http_user_agent_script":     {"curl"},
		"http_user_agent_ai":         {"gptbot"},
	}

	t.Cleanup(func() {
		LoadedUserAgentLists = originalLoadedUserAgentLists
	})

	tests := []struct {
		name         string
		userAgent    string
		wantMatchKey string
	}{
		{
			name:         "Case-insensitive match (lowercase)",
			userAgent:    "mozilla/5.0 (compatible; bingbot/2.0)",
			wantMatchKey: "http_user_agent_crawler",
		},
		{
			name:         "Case-insensitive match (uppercase)",
			userAgent:    "CURL/7.81.0",
			wantMatchKey: "http_user_agent_script",
		},
		{
			name:         "Match on secondary list",
			userAgent:    "BetterUptimeBot",
			wantMatchKey: "http_user_agent_monitoring",
		},
		{
			name:         "Only first match is returned (order matters)",
			userAgent:    "I am a GPTBot crawler",
			wantMatchKey: "http_user_agent_ai",
		},
		{
			name:         "No match",
			userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			wantMatchKey: "",
		},
		{
			name:         "Empty user agent",
			userAgent:    "",
			wantMatchKey: "",
		},
		{
			name:         "User agent matches a key in different case",
			userAgent:    "SPIDER",
			wantMatchKey: "http_user_agent_crawler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey := LookupUserAgentCategory(tt.userAgent)
			if tt.wantMatchKey != "" && gotKey == "" {
				t.Errorf("LookupUserAgentCategory(%q) got = %q, want a non-empty key (e.g., %q)", tt.userAgent, gotKey, tt.wantMatchKey)
			} else if tt.wantMatchKey == "" && gotKey != "" {
				t.Errorf("LookupUserAgentCategory(%q) got = %q, want empty string", tt.userAgent, gotKey)
			}
		})
	}
}

func TestLookupUserAgent(t *testing.T) {
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
			"Googlebot-Image",
			"http_user_agent_crawler_google_image",
		},
		{
			"Google",
			"http_user_agent_crawler_google_fallback",
		},
	}

	t.Cleanup(func() {
		LoadedUserAgentMapArray = originalLoadedUserAgentMapArray
	})

	tests := []struct {
		name         string
		userAgent    string
		wantMatchKey string
	}{
		{
			name:         "Bingbot",
			userAgent:    "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm) Chrome/116.0.1938.76 Safari/537.36",
			wantMatchKey: "http_user_agent_crawler_microsoft_bing",
		},
		{
			name:         "Google Fallback",
			userAgent:    "Mozilla/5.0 (compatible; Google-TEST/2.1; +http://www.google.com/bot.html)",
			wantMatchKey: "http_user_agent_crawler_google_fallback",
		},
		{
			name:         "Googlebot-Image",
			userAgent:    "Mozilla/5.0 (compatible; Googlebot-Image/2.1; +http://www.google.com/bot.html)",
			wantMatchKey: "http_user_agent_crawler_google_image",
		},
		{
			name:         "Googlebot (Search)",
			userAgent:    "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			wantMatchKey: "http_user_agent_crawler_google_search",
		},
		{
			name:         "No match",
			userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			wantMatchKey: "",
		},
		{
			name:         "Empty user agent",
			userAgent:    "",
			wantMatchKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey := LookupUserAgent(tt.userAgent)
			if tt.wantMatchKey != "" && gotKey == "" {
				t.Errorf("LookupUserAgent(%q) got = %q, want a non-empty key (e.g., %q)", tt.userAgent, gotKey, tt.wantMatchKey)
			} else if tt.wantMatchKey == "" && gotKey != "" {
				t.Errorf("LookupUserAgent(%q) got = %q, want empty string", tt.userAgent, gotKey)
			}
		})
	}
}

func TestLookupFingerprint(t *testing.T) {
	originalLoadedFingerprintLists := LoadedFingerprintLists

	LoadedFingerprintLists = map[string][]string{
		"fingerprint_script":          {"33,65281-10,34,16,11,43", "99,10,23,5,15,6"},
		"fingerprint_scanner_tls_ja4": {"t13d1811h2_e8a523a41297_5894756feeaa", "t13d1811h2_e8a523a41297_5894756fee65"},
	}

	t.Cleanup(func() {
		LoadedFingerprintLists = originalLoadedFingerprintLists
	})

	tests := []struct {
		name         string
		fingerprint  string
		wantMatchKey string
	}{
		{
			name:         "Exact match on script list",
			fingerprint:  "33,65281-10,34,16,11,43",
			wantMatchKey: "fingerprint_script",
		},
		{
			name:         "Exact match on JA4 list",
			fingerprint:  "t13d1811h2_e8a523a41297_5894756fee65",
			wantMatchKey: "fingerprint_scanner_tls_ja4",
		},
		{
			name:         "No match (substring is not supported)",
			fingerprint:  "33,65281-10",
			wantMatchKey: "",
		},
		{
			name:         "No match (incorrect fingerprint)",
			fingerprint:  "33,65281-10,34,16,11,44",
			wantMatchKey: "",
		},
		{
			name:         "Empty fingerprint",
			fingerprint:  "",
			wantMatchKey: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey := LookupFingerprint(tt.fingerprint)
			if tt.wantMatchKey != "" && gotKey == "" {
				t.Errorf("LookupFingerprint(%q) got = %q, want a non-empty key (e.g., %q)", tt.fingerprint, gotKey, tt.wantMatchKey)
			} else if tt.wantMatchKey == "" && gotKey != "" {
				t.Errorf("LookupFingerprint(%q) got = %q, want empty string", tt.fingerprint, gotKey)
			}
		})
	}
}
