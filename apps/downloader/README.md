# OXL Open Bot List - Downloader Application

This application can be used to download and prepare the Open-Bot-List data for your local integration.

This repository also contains an application for downloading and parsing these infos so you can easily use/implement them in you systems.

----

## How it works

* Downloads the [match-information](https://github.com/O-X-L/open-bot-list/tree/latest/matches) from 'raw.githubusercontent.com'
* Parses the match-files
* Generates simple lookup-files for your services

----

## Setup

It is a simple single binary file that can be executed. No install required.

Download: [Releases](https://github.com/O-X-L/open-bot-list/releases)

Or build [with Go](https://go.dev/doc/install): `cd downloader/ && go build -o './open-bot-list-downloader' ./cmd/main.go`

----

## Usage

```bash
rath@gate:~ ./open-bot-list-downloader

OXL Open-Bot-List Downloader v1.0
> © OXL IT Services / Rath Pascal
> git.oxl.at/open-bot-list
> License: GPLv3

Usage of build/downloader:
  -output-dir string
        Output directory for processed data (required).
  -runtime-dir string
        Runtime directory to store downloaded manifest files. (default "/tmp/oxl-open-bot-list")
```

----

## Output

Currently, these kinds of files are generated:

* **User-Agents / Fingerprints:**
  * Map-files: maps `match_name` to the value to match
  * Lists: the whole file is a `match_name` - all values inside it need to be checked 

* **IP/Network-Lists:**
  * IPv4 Networks in CIDR
  * IPv6 Networks in CIDR
  * IPv4 & IPv6 Networks in CIDR combined

**Roadmap**:

* ASN Lists ([Categorization: Hosting providers, ISPs, Education](https://github.com/O-X-L/risk-db-lists/tree/main/asn))
* PTR's for crawler-verification

### IP/Network-List Validation & Sanitization

These processing steps are performed:

* Bogon networks are skipped (*internal ranges and so on*)
* Invalid IPs/Networks are skipped
* You will see a warning if many IPs are included in a list (*>10_000 for now*)
* Where possible - Networks are summarized (*p.e. 10.0.0.0/24, 10.0.1.0/24, 10.0.1.1/32 = 10.0.0.0/23*)

### Example

<details>

```bash
rath@gate:~ head -n 5 /tmp/oxl-open-bot-list-out/http_user_agent_ai_sub.map 
http_user_agent_crawler_google_aidata1 Google-Extended
http_user_agent_crawler_google_aidata2 Google-LLM-Research
http_user_agent_crawler_google_aiuser1 Google-CloudVertexBot
http_user_agent_crawler_google_aiuser2 GoogleAI-ContentFetcher
http_user_agent_crawler_google_aidata3 Gemini-Deep-Research

rath@gate:~ head -n 5 /tmp/oxl-open-bot-list-out/http_user_agent_ai_sub.lst 
Applebot-Extended
Amazonbot
Perplexity-User
AliyunSecBot
Google-CloudVertexBot

rath@gate:~ head -n 5 /tmp/oxl-open-bot-list-out/src_net_crawler_google_common_net4.lst 
34.22.85.0/27
34.64.82.64/28
34.65.242.112/28
34.80.50.80/28
34.88.194.0/28

rath@gate:~ head -n 5 /tmp/oxl-open-bot-list-out/src_net_crawler_google_common_net6.lst 
2001:4860:4801:2::/64
2001:4860:4801:c::/64
2001:4860:4801:f::/64
2001:4860:4801:10::/64
2001:4860:4801:12::/63

rath@gate:~ head -n 5 /tmp/oxl-open-bot-list-out/fingerprint_script_tls_ja4.lst 
t13d1909h2_9dc949149365_97f8aa674fd9
t13d181100_85036bcba153_24695f2957a7
t12d3805h1_10ed599f3404_aaf95bb78ec9
t12d1909h2_ab14d9cb224d_70a28de75618
t13d9112h2_0d5420ba6086_78b3e9c34d1f

rath@gate:~ head -n 5 /tmp/oxl-open-bot-list-out/fingerprint_script.map 
t13i4311h1_c7886603b240_b26ce05bbdd6 python-requests
t13i3111h1_e8f1e7e78f70_d41ae481755e Python aiohttp
t13i3111h1_e8f1e7e78f70_b26ce05bbdd6 python-requests
t13i181000_85036bcba153_d41ae481755e Python-urllib
t13i1712h1_ab0a1bf427ad_ecd0401ec68b Python aiohttp

rath@gate:~ tree /tmp/oxl-open-bot-list-out
/tmp/oxl-open-bot-list-out
├── fingerprint_crawler.lst
├── fingerprint_crawler.map
├── fingerprint_crawler_tls_ja4.lst
├── fingerprint_scanner.lst
├── fingerprint_scanner.map
├── fingerprint_scanner_tls_ja4.lst
├── fingerprint_script.lst
├── fingerprint_script.map
├── fingerprint_script_tls_ja4.lst
├── http_user_agent_ai_sub.lst
├── http_user_agent_ai_sub.map
├── http_user_agent_crawler_sub.lst
├── http_user_agent_crawler_sub.map
├── http_user_agent_ecommerce_sub.lst
├── http_user_agent_ecommerce_sub.map
├── http_user_agent_monitoring_sub.lst
├── http_user_agent_monitoring_sub.map
├── http_user_agent_random_sub.lst
├── http_user_agent_random_sub.map
├── http_user_agent_scanner_sub.lst
├── http_user_agent_scanner_sub.map
├── http_user_agent_script_sub.lst
├── http_user_agent_script_sub.map
├── http_user_agent_software_sub.lst
├── http_user_agent_software_sub.map
├── src_net_cdn_bunnyway_all.lst
├── src_net_cdn_bunnyway_net4.lst
├── src_net_cdn_bunnyway_net6.lst
├── src_net_cdn_cloudflare_all.lst
├── src_net_cdn_cloudflare_net4.lst
├── src_net_cdn_cloudflare_net6.lst
├── src_net_cdn_fastly_all.lst
├── src_net_cdn_fastly_net4.lst
├── src_net_cdn_fastly_net6.lst
├── src_net_crawler_ahrefs_all.lst
├── src_net_crawler_ahrefs_net4.lst
├── src_net_crawler_ahrefs_net6.lst
├── src_net_crawler_apple_all.lst
├── src_net_crawler_apple_net4.lst
├── src_net_crawler_apple_net6.lst
├── src_net_crawler_commoncrawl_all.lst
├── src_net_crawler_commoncrawl_net4.lst
├── src_net_crawler_commoncrawl_net6.lst
├── src_net_crawler_duckduckgo_aiuser_all.lst
├── src_net_crawler_duckduckgo_aiuser_net4.lst
├── src_net_crawler_duckduckgo_aiuser_net6.lst
├── src_net_crawler_duckduckgo_search_all.lst
├── src_net_crawler_duckduckgo_search_net4.lst
├── src_net_crawler_duckduckgo_search_net6.lst
├── src_net_crawler_ecom_stripe_all.lst
├── src_net_crawler_ecom_stripe_net4.lst
├── src_net_crawler_ecom_stripe_net6.lst
├── src_net_crawler_google_common_all.lst
├── src_net_crawler_google_common_net4.lst
├── src_net_crawler_google_common_net6.lst
├── src_net_crawler_google_special_all.lst
├── src_net_crawler_google_special_net4.lst
├── src_net_crawler_google_special_net6.lst
├── src_net_crawler_google_user1_all.lst
├── src_net_crawler_google_user1_net4.lst
├── src_net_crawler_google_user1_net6.lst
├── src_net_crawler_google_user2_all.lst
├── src_net_crawler_google_user2_net4.lst
├── src_net_crawler_google_user2_net6.lst
├── src_net_crawler_microsoft_bing_all.lst
├── src_net_crawler_microsoft_bing_net4.lst
├── src_net_crawler_microsoft_bing_net6.lst
├── src_net_crawler_mistralai_all.lst
├── src_net_crawler_mistralai_net4.lst
├── src_net_crawler_mistralai_net6.lst
├── src_net_crawler_mon_betterstack_all.lst
├── src_net_crawler_mon_betterstack_net4.lst
├── src_net_crawler_mon_betterstack_net6.lst
├── src_net_crawler_mon_catchpoint_all.lst
├── src_net_crawler_mon_catchpoint_net4.lst
├── src_net_crawler_mon_catchpoint_net6.lst
├── src_net_crawler_mon_qualys_all.lst
├── src_net_crawler_mon_qualys_net4.lst
├── src_net_crawler_mon_qualys_net6.lst
├── src_net_crawler_mon_solarwinds_all.lst
├── src_net_crawler_mon_solarwinds_net4.lst
├── src_net_crawler_mon_solarwinds_net6.lst
├── src_net_crawler_mon_uptimerobot_all.lst
├── src_net_crawler_mon_uptimerobot_net4.lst
├── src_net_crawler_mon_uptimerobot_net6.lst
├── src_net_crawler_openai_aidata_all.lst
├── src_net_crawler_openai_aidata_net4.lst
├── src_net_crawler_openai_aidata_net6.lst
├── src_net_crawler_openai_search_all.lst
├── src_net_crawler_openai_search_net4.lst
├── src_net_crawler_openai_search_net6.lst
├── src_net_crawler_openai_user_all.lst
├── src_net_crawler_openai_user_net4.lst
├── src_net_crawler_openai_user_net6.lst
├── src_net_crawler_perplexity_aidata_all.lst
├── src_net_crawler_perplexity_aidata_net4.lst
├── src_net_crawler_perplexity_aidata_net6.lst
├── src_net_crawler_perplexity_user_all.lst
├── src_net_crawler_perplexity_user_net4.lst
├── src_net_crawler_perplexity_user_net6.lst
├── src_net_crawler_qwant_all.lst
├── src_net_crawler_qwant_net4.lst
├── src_net_crawler_qwant_net6.lst
├── src_net_crawler_seekport_all.lst
├── src_net_crawler_seekport_net4.lst
├── src_net_crawler_seekport_net6.lst
├── src_net_crawler_telegram_all.lst
├── src_net_crawler_telegram_net4.lst
├── src_net_crawler_telegram_net6.lst
├── src_net_proxy_tor_all.lst
├── src_net_proxy_tor_net4.lst
├── src_net_proxy_tor_net6.lst
├── src_net_vpn_apple_privacyrelay_all.lst
├── src_net_vpn_apple_privacyrelay_net4.lst
└── src_net_vpn_apple_privacyrelay_net6.lst
```

</details>
