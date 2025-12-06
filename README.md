# Open Bot List

<p align="center">
    <a title="Support this Project (Donate, Support-Licenses)" href="https://shop.oxl.app/collections/open-source">
        <img src="https://files.oxl.at/img/badge-oss-support.svg" alt="Support Badge (Donate, Support-Licenses)"/>
    </a>
</p>

[![Lint CSV Configs](https://github.com/O-X-L/open-bot-list/actions/workflows/csv_lint.yml/badge.svg)](https://github.com/O-X-L/open-bot-list/actions/workflows/csv_lint.yml)
[![Lint Downloader](https://github.com/O-X-L/open-bot-list/actions/workflows/dl_lint.yml/badge.svg)](https://github.com/O-X-L/open-bot-list/actions/workflows/dl_lint.yml)
[![Unit Test Downloader](https://github.com/O-X-L/open-bot-list/actions/workflows/dl_unit_test.yml/badge.svg)](https://github.com/O-X-L/open-bot-list/actions/workflows/dl_unit_test.yml)
[![Functional Tests for Downloader](https://github.com/O-X-L/open-bot-list/actions/workflows/dl_functional_test.yml/badge.svg)](https://github.com/O-X-L/open-bot-list/actions/workflows/dl_functional_test.yml)

This repository is used to collect information that can be used to categorize & match traffic.

We will auto-generate full lists in plaintext and JSON later on!

----

## Contribute

Contributions are very welcome.

If you:

* know of official IP-Lists we missed
* found other missing/incorrect information

..feel free to either [open a ticket](https://github.com/O-X-L/open-bot-list/issues) or [email us directly](mailto://contact+openbotlist@OXL.at)

---

## How it works

To transparently match & categorize bots we need to combine:

* **Traffic Matches**
  * Matching the source-IP with IP- or ASN-Lists
    * Separating different kinds of bots by their HTTP User-Agent (*if they use the same IP-range*)
    * Categorizing the source-IP into hosting/vpn/isp/proxy/isp-cgnat (*not that easy.. (; *)
  * Separating different bot-categories like:
    `script bots`, `hidden bots`, `search-engine crawlers`, `AI-data crawlers`, `AI-user crawlers`, `social-media crawlers`, `crawlers for ADs`, `crawlers for ecommerce`, and so on

    * Matching clear script-bots by their User-Agent (*dumb script-kiddies*)
    * Matching 'hidden' bots by their client-fingerprints (*[JA4](https://github.com/O-X-L/haproxy-ja4-fingerprint), etc.*)
    * ... *to be extended* ...

* **PTR-checks**
  * Some organizations only supply us with a PTR-match to validate if a crawler-IP is theirs (*no simple IP-list lookups*)

* **Traffic Flagging**
  * We provide you with abstract configuration that shows how the matches can be combined
  * Practical configuration examples for proxy-services will be added later on

----

## Downloader Application

This repository also contains an application for downloading and parsing these infos so you can easily use/implement them in you systems.

Download: [Releases](https://github.com/O-X-L/open-bot-list/releases)

Or build [with Go](https://go.dev/doc/install): `cd downloader/ && go build -o './open-bot-list-downloader' ./cmd/main.go`

### Usage

```bash
rath@gate:~ ./open-bot-list-downloader

OXL Open-Bot-List Downloader
> © OXL IT Services / Rath Pascal
> git.oxl.at/open-bot-list
> License: BSD-3-Clause

Usage of build/downloader:
  -output-dir string
        Output directory for processed data (required).
  -runtime-dir string
        Runtime directory to store downloaded manifest files. (default "/tmp/oxl-open-bot-list")
```

### Output

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

#### Example

<details>

```bash
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

----

## Motivation

We are working on building a [FOSS WAF-platform](https://github.com/O-X-L/shieldwall-waf) (*and centrally manageable network-firewalls*) which require such a collection of bot-related information.

With our [IP-Abuse Reporting-System & Databases](https://github.com/O-X-L/risk-db) we have already started to collect information for it.

As the mindset of Open-Source is at the core of our being - we want to transparently share it with the whole world.

----

### SHIELD-WALL WAF Project

This information-collection is part of our [SHIELD-WALL WAF Project](https://github.com/O-X-L/shieldwall-waf).

Check-out the demo: [demo.waf.shield-wall.net](https://demo.waf.shield-wall.net)
