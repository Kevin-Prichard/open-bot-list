package config

// List of check-result flag-pairs
// If all check-flags exist - the resulting flags get applied
// Only one rule can be matched! Fallback/Catchall rules should be placed below specific ones!
// see also: https://github.com/O-X-L/open-bot-list/tree/latest/flagging
var FLAGGING_RULESET = [][][]string{
	{ // Scanners from Tor-Network
		[]string{"org_proxy_tor", "fingerprint_scanner_tls_ja4"}, // <== flags to check/match
		[]string{"bot_scanner"},                                  // <== resulting flags if all matched
	},
	{ // Script-Bots from Tor-Network
		[]string{"org_proxy_tor", "fingerprint_script_tls_ja4"},
		[]string{"bot_script"},
	},
	{
		[]string{"org_proxy_tor", "http_user_agent_crawler"},
		[]string{"bot_crawler", "crawler_spoofed"},
	},
	{
		[]string{"org_proxy_tor", "http_user_agent_ai"},
		[]string{"bot_crawler", "crawler_spoofed"},
	},
	// ### AI ###
	{
		[]string{"src_net_crawler_apple", "http_user_agent_crawler_applebot_extended"},
		[]string{"bot_crawler", "crawler_ai_data", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_openai_search", "http_user_agent_crawler_openai_search"},
		[]string{"bot_crawler", "crawler_search", "crawler_ai_user", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_openai_user", "http_user_agent_crawler_openai_user"},
		[]string{"bot_crawler", "crawler_user", "crawler_ai_user", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_openai_aidata", "src_net_crawler_openai_aidata"},
		[]string{"bot_crawler", "crawler_ai_data", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_perplexity_user", "http_user_agent_crawler_perplexity_user"},
		[]string{"bot_crawler", "crawler_user", "crawler_ai_user", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_perplexity_aidata", "http_user_agent_crawler_perplexity_aidata"},
		[]string{"bot_crawler", "crawler_ai_data", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_mistralai", "http_user_agent_crawler_mistral_user"},
		[]string{"bot_crawler", "crawler_user", "crawler_ai_user", "crawler_ai_data", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_duckduckgo_aiuser", "http_user_agent_crawler_duckduckgo_aiuser"},
		[]string{"bot_crawler", "crawler_user", "crawler_ai_user", "crawler_verified"},
	},
	{ // todo: amazon IP-verification
		[]string{"http_user_agent_crawler_amazon_main"},
		[]string{"bot_crawler", "crawler_ai_data", "org_amazon"},
	},
	{ // todo: meta/facebook IP-verification
		[]string{"http_user_agent_crawler_meta_aidata"},
		[]string{"bot_crawler", "crawler_ai_data", "org_meta"},
	},
	{ // todo: alibaba IP-verification
		[]string{"http_user_agent_crawler_alibaba_aidata"},
		[]string{"bot_crawler", "crawler_ai_data", "org_alibaba"},
	},
	{ // todo: exaAI IP-verification
		[]string{"http_user_agent_crawler_exa_aisearch"},
		[]string{"bot_crawler", "crawler_ai_user", "org_exa"},
	},
	{ // todo: cohereAI IP-verification
		[]string{"http_user_agent_crawler_cohere_ai"},
		[]string{"bot_crawler", "crawler_user", "crawler_ai_user", "crawler_ai_data", "org_cohere"},
	},
	{ // Google-Extended
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_aidata1"},
		[]string{"bot_crawler", "crawler_ai_data", "crawler_verified"},
	},
	{ // Google-LLM-Research
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_aidata2"},
		[]string{"bot_crawler", "crawler_ai_data", "crawler_verified"},
	},
	{ // Google Gemini Deep-Research
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_aidata3"},
		[]string{"bot_crawler", "crawler_ai_data", "crawler_verified"},
	},
	{ // Google-CloudVertexBot
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_aiuser1"},
		[]string{"bot_crawler", "crawler_user", "crawler_ai_user", "crawler_verified"},
	},
	{ // GoogleAI-ContentFetcher
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_aiuser2"},
		[]string{"bot_crawler", "crawler_user", "crawler_ai_user", "crawler_verified"},
	},
	// ### CRAWLERS ###
	{
		[]string{"src_net_crawler_apple", "http_user_agent_crawler_applebot"},
		[]string{"bot_crawler", "crawler_search", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_duckduckgo_search", "http_user_agent_crawler_duckduckgo_search"},
		[]string{"bot_crawler", "crawler_search", "crawler_verified"},
	},
	{ // todo: amazon IP-verification
		[]string{"http_user_agent_crawler_amazon_user"},
		[]string{"bot_crawler", "crawler_user", "org_amazon"},
	},
	{ // todo: meta/facebook IP-verification
		[]string{"http_user_agent_crawler_meta_ads"},
		[]string{"bot_crawler", "crawler_socialmedia", "crawler_ads", "org_meta"},
	},
	{ // todo: meta/facebook IP-verification
		[]string{"http_user_agent_crawler_meta_user"},
		[]string{"bot_crawler", "crawler_socialmedia", "crawler_user", "org_meta"},
	},
	{ // todo: meta/facebook IP-verification
		[]string{"http_user_agent_crawler_meta_facebook"},
		[]string{"bot_crawler", "crawler_socialmedia", "org_meta"},
	},
	{
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_search"},
		[]string{"bot_crawler", "crawler_search", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_image"},
		[]string{"bot_crawler", "crawler_search", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_video"},
		[]string{"bot_crawler", "crawler_search", "crawler_verified"},
	},
	{ // todo: user-agent matching
		[]string{"src_net_crawler_google_user1"},
		[]string{"bot_crawler", "crawler_user", "crawler_verified"},
	},
	{ // todo: user-agent matching
		[]string{"src_net_crawler_google_user2"},
		[]string{"bot_crawler", "crawler_user", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_google_special", "http_user_agent_crawler_google_ads1"},
		[]string{"bot_crawler", "crawler_ads", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_google_special", "http_user_agent_crawler_google_ads2"},
		[]string{"bot_crawler", "crawler_ads", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_store"},
		[]string{"bot_crawler", "crawler_ecommerce", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_google_special", "http_user_agent_crawler_google_fallback"},
		[]string{"bot_crawler", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_google_common", "http_user_agent_crawler_google_fallback"},
		[]string{"bot_crawler", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_microsoft_bing", "http_user_agent_crawler_microsoft_bing"},
		[]string{"bot_crawler", "crawler_search", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_microsoft_bing", "http_user_agent_crawler_microsoft_msn"},
		[]string{"bot_crawler", "crawler_search", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_ecom_stripe", "http_user_agent_ecom_stripe"},
		[]string{"bot_crawler", "crawler_ecommerce", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_ahrefs", "http_user_agent_crawler_ahrefs"},
		[]string{"bot_crawler", "crawler_search", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_commoncrawl", "http_user_agent_crawler_commoncrawl"},
		[]string{"bot_crawler", "crawler_ai_data", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_qwant", "http_user_agent_crawler_qwant"},
		[]string{"bot_crawler", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_seekport", "http_user_agent_crawler_seekport"},
		[]string{"bot_crawler", "crawler_verified"},
	},
	// ### MONITORING ###
	{
		[]string{"src_net_crawler_mon_uptimerobot", "http_user_agent_mon_uptimerobot"},
		[]string{"bot_crawler", "crawler_monitoring", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_mon_betterstack", "http_user_agent_mon_betteruptime"},
		[]string{"bot_crawler", "crawler_monitoring", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_mon_solarwinds", "http_user_agent_mon_solarwinds"},
		[]string{"bot_crawler", "crawler_monitoring", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_mon_catchpoint", "http_user_agent_mon_catchpoint"},
		[]string{"bot_crawler", "crawler_monitoring", "crawler_verified"},
	},
	{
		[]string{"src_net_crawler_mon_qualys"},
		[]string{"bot_crawler", "crawler_monitoring"},
	},
	// ### IMPLICIT ###
	{
		[]string{"http_user_agent_scanner"},
		[]string{"bot_scanner"},
	},
	{
		[]string{"fingerprint_scanner_tls_ja4"},
		[]string{"bot_scanner"},
	},
	{
		[]string{"http_user_agent_script"},
		[]string{"bot_script"},
	},
	{
		[]string{"fingerprint_script_tls_ja4"},
		[]string{"bot_script"},
	},
	{
		[]string{"http_user_agent_ai"},
		[]string{"bot_crawler", "crawler_ai_data"},
	},
	{
		[]string{"http_user_agent_software"},
		[]string{"bot_software"},
	},
	{
		[]string{"http_user_agent_crawler"},
		[]string{"bot_crawler"},
	},
	{
		[]string{"http_user_agent_ai"},
		[]string{"bot_crawler"},
	},
}

var FALLBACK_SPOOFED_UA_SUB = []string{
	"google",
	"bing",
	"gptbot",
}
