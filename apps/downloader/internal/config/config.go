package config

import (
	"os"
	"path/filepath"
	"time"
)

const (
	URL_BASE_MATCHES      = "https://raw.githubusercontent.com/O-X-L/open-bot-list/refs/heads/latest/matches"
	VALUE_MULTI_DELIMITER = "|"
	USER_AGENT_STRING     = "OXL Open-Bot-List (http://git.oxl.at/open-bot-list)"
	FILE_PREFIX_MATCH     = "match_"
	FILE_PREFIX_IPLIST    = "iplist_"
	FILE_PREFIX_ASNLIST   = "asnlist_"
	VERSION               = 1.4
)

var (
	PATH_RUNTIME           = filepath.Join(os.TempDir(), "oxl-open-bot-list")
	PATH_OUTPUT            = ""
	MIN_LIST_LIFETIME_HOUR = 4
	MIN_LIST_LIFETIME      = 4 * time.Hour
	WARN_IPLIST_COUNT      = 10000
	WARN_ASNLIST_COUNT     = 100
)

var IPListCategories = map[string]string{
	"_overall":   "ip_net/_overall.csv",
	"ai":         "ip_net/ai.csv",
	"cdn":        "ip_net/cdn.csv",
	"crawler":    "ip_net/crawler.csv",
	"ecommerce":  "ip_net/ecommerce.csv",
	"malicious":  "ip_net/malicious.csv",
	"monitoring": "ip_net/monitoring.csv",
	"proxy":      "ip_net/proxy.csv",
	"vpn":        "ip_net/vpn.csv",
}

var FingerprintCategories = map[string]string{
	"_overall": "fingerprint/_overall.csv",
	"crawler":  "fingerprint/crawler.csv",
	"scanner":  "fingerprint/scanner.csv",
	"script":   "fingerprint/script.csv",
}

var PtrCategories = map[string]string{
	"_overall":      "ptr/_overall.csv",
	"crawler":       "ptr/crawler.csv",
	"hosting":       "ptr/hosting.csv",
	"isp_base":      "ptr/isp_base.csv",
	"isp_cgnat":     "ptr/isp_cgnat.csv",
	"isp_copper":    "ptr/isp_copper.csv",
	"isp_dynamic":   "ptr/isp_dynamic.csv",
	"isp_satellite": "ptr/isp_satellite.csv",
	"isp_static":    "ptr/isp_static.csv",
	"isp_wireless":  "ptr/isp_wireless.csv",
}

var UserAgentCategories = map[string]string{
	"_overall":   "user_agent/_overall.csv",
	"ai":         "user_agent/ai.csv",
	"crawler":    "user_agent/crawler.csv",
	"ecommerce":  "user_agent/ecommerce.csv",
	"monitoring": "user_agent/monitoring.csv",
	"unknown":    "user_agent/unknown.csv",
	"scanner":    "user_agent/scanner.csv",
	"script":     "user_agent/script.csv",
	"software":   "user_agent/software.csv",
}

var ASNCategories = map[string]string{
	"_overall":  "asn/_overall.csv",
	"cdn":       "asn/cdn.csv",
	"cloud":     "asn/cloud.csv",
	"crawler":   "asn/crawler.csv",
	"education": "asn/education.csv",
	"hosting":   "asn/hosting.csv",
	"isp":       "asn/isp.csv",
	"proxy":     "asn/proxy.csv",
	"scanner":   "asn/scanner.csv",
	"vpn":       "asn/vpn.csv",
}

var ASNListCategories = map[string]string{
	"_overall":  "asn_list/_overall.csv",
	"malicious": "asn_list/malicious.csv",
}
