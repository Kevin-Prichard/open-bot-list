package config

// NOTE: ls shieldwall-waf/testdata/frontend/log_flagger/open_bot_list/ | cut -d '.' -f1 | rev | cut -d '_' -f2- | rev | uniq

// IP/Network-List Files to load - mapping to file-specific flags
var IPLIST_FLAGS = map[string][]string{
	"src_net_cdn_bunnyway":              {"src_cdn", "org_cdn_bunnyway"},
	"src_net_cdn_cloudflare":            {"src_cdn", "org_cdn_cloudflare"},
	"src_net_cdn_fastly":                {"src_cdn", "org_cdn_fastly"},
	"src_net_crawler_ahrefs":            {"bot", "bot_crawler", "org_ahrefs"},
	"src_net_crawler_apple":             {"bot", "org_apple"},
	"src_net_crawler_commoncrawl":       {"bot", "bot_crawler", "crawler_ai_data", "org_commoncrawl"},
	"src_net_crawler_duckduckgo_aiuser": {"bot", "bot_crawler", "crawler_ai_user", "crawler_user", "org_duckduckgo"},
	"src_net_crawler_duckduckgo_search": {"bot", "bot_crawler", "crawler_search", "org_duckduckgo"},
	"src_net_crawler_ecom_stripe":       {"bot", "bot_ecommerce", "org_stripe"},
	"src_net_crawler_google_common":     {"bot", "org_google"},
	"src_net_crawler_google_special":    {"bot", "org_google"},
	"src_net_crawler_google_user1":      {"bot", "bot_crawler", "crawler_user", "org_google"},
	"src_net_crawler_google_user2":      {"bot", "bot_crawler", "crawler_user", "org_google"},
	"src_net_crawler_microsoft_bing":    {"bot", "org_microsoft"},
	"src_net_crawler_mistralai":         {"bot", "bot_crawler", "org_mistralai"},
	"src_net_crawler_mon_betterstack":   {"bot", "bot_monitoring", "org_betterstack"},
	"src_net_crawler_mon_catchpoint":    {"bot", "bot_monitoring", "org_catchpoint"},
	"src_net_crawler_mon_qualys":        {"bot", "bot_monitoring", "org_qualys"},
	"src_net_crawler_mon_solarwinds":    {"bot", "bot_monitoring", "org_solarwinds"},
	"src_net_crawler_mon_uptimerobot":   {"bot", "bot_monitoring", "org_uptimerobot"},
	"src_net_crawler_openai_aidata":     {"bot", "bot_crawler", "crawler_ai_data", "org_openai"},
	"src_net_crawler_openai_search":     {"bot", "bot_crawler", "org_openai"},
	"src_net_crawler_openai_user":       {"bot", "bot_crawler", "org_openai"},
	"src_net_crawler_perplexity_aidata": {"bot", "bot_crawler", "crawler_ai_data", "org_perplexity"},
	"src_net_crawler_perplexity_user":   {"bot", "bot_crawler", "org_perplexity"},
	"src_net_crawler_qwant":             {"bot", "bot_crawler", "org_qwant"},
	"src_net_crawler_seekport":          {"bot", "bot_crawler", "org_seekport"},
	"src_net_crawler_telegram":          {"bot", "bot_software", "org_telegram"},
	"src_net_proxy_tor":                 {"src_ip_proxy", "org_proxy_tor"},
	"src_net_vpn_apple_privacyrelay":    {"src_ip_vpn", "org_vpn_apple"},
}

// User-Agent Files to load - mapping to file-specific flags
var USER_AGENT_LIST_FLAGS = map[string][]string{
	"http_user_agent_ai":         {"bot"},
	"http_user_agent_crawler":    {"bot"},
	"http_user_agent_ecommerce":  {"bot"},
	"http_user_agent_monitoring": {"bot"},
	"http_user_agent_random":     {"bot"},
	"http_user_agent_scanner":    {"bot"},
	"http_user_agent_script":     {"bot"},
	"http_user_agent_software":   {"bot"},
}

// Fingerprint Files to load & kinds to process - mapping to file-specific flags
// do not directly flag as bots as these might still be false-positives..
var FINGERPRINT_LIST_FLAGS = map[string][]string{
	"fingerprint_crawler":         {},
	"fingerprint_crawler_tls_ja4": {"fingerprint_crawler_tls_ja4"},
	"fingerprint_scanner":         {},
	"fingerprint_scanner_tls_ja4": {"fingerprint_scanner_tls_ja4"},
	"fingerprint_script":          {},
	"fingerprint_script_tls_ja4":  {"fingerprint_script_tls_ja4"},
}

// PTR Files to load - mapping to file-specific flags
var PTR_LIST_FLAGS = map[string][]string{
	"ptr_crawler": {"bot", "bot_crawler"},
}

// ASN Files to load - mapping to file-specific flags
var ASN_LIST_FLAGS = map[string][]string{
	"src_asn_cdn":       {},
	"src_asn_cloud":     {},
	"src_asn_crawler":   {},
	"src_asn_education": {},
	"src_asn_hosting":   {},
	"src_asn_isp":       {},
	"src_asn_proxy":     {},
	"src_asn_scanner":   {},
	"src_asn_vpn":       {},
}
