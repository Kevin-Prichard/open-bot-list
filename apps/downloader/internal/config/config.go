package config

import (
	"os"
	"path/filepath"
)

const (
	URL_BASE_MATCHES      = "https://raw.githubusercontent.com/O-X-L/open-bot-list/refs/heads/latest/matches"
	VALUE_MULTI_DELIMITER = "|"
	USER_AGENT_STRING     = "OXL Open-Bot-List (http://git.oxl.at/open-bot-list)"
	FILE_PREFIX_MATCH     = "match_"
	FILE_PREFIX_IPLIST    = "iplist_"
	VERSION               = 1.1
)

var (
	PATH_RUNTIME = filepath.Join(os.TempDir(), "oxl-open-bot-list")
	PATH_OUTPUT  = ""
)

var IPListCategories = map[string]string{
	"_overall":   "ip_net/_overall.csv",
	"ai":         "ip_net/ai.csv",
	"cdn":        "ip_net/cdn.csv",
	"crawler":    "ip_net/crawler.csv",
	"ecommerce":  "ip_net/ecommerce.csv",
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
	"_overall": "ptr/_overall.csv",
	"crawler":  "ptr/crawler.csv",
}

var UserAgentCategories = map[string]string{
	"_overall":   "user_agent/_overall.csv",
	"ai":         "user_agent/ai.csv",
	"crawler":    "user_agent/crawler.csv",
	"ecommerce":  "user_agent/ecommerce.csv",
	"monitoring": "user_agent/monitoring.csv",
	"random":     "user_agent/random.csv",
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
