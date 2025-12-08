# OXL Open Bot List - Log-Flagger Application

This application can be used to apply traffic-flags **on logs in CSV-format**.

Enriching your log-data can be useful to perform ad-hoc/asynchronous traffic analysis - per example to detect malicious bot-traffic that you could limit or block.

In an enterprise-grade environment you can include such flagging in your log-pipelines. This application is not (yet?) designed to be used in such.

----

## Usage

```bash
rath@gate:~ ./log_flagger 

OXL Open-Bot-List Log-Flagger v1.0
> © OXL IT Services / Rath Pascal
> git.oxl.at/open-bot-list
> License: GPLv3

Usage of build/log_flagger:
  -csv-field-client-ip int
        Index of the CSV-Field containing the timestamp. (default 0)
  -csv-field-fp-ja4 int
        Index of the CSV-Field containing the JA4 client-fingerprint. (default 1)
  -csv-field-user-agent int
        Index of the CSV-Field containing the User-Agent. (default 2)
  -data-dir string
        Path to the directory containing the 'open-bot-list' data (required). See: https://github.com/O-X-L/open-bot-list/tree/latest?tab=readme-ov-file#downloader-application
  -debug
        Enable debug output.
  -input-file string
        Path to the input CSV-file to process (required).
  -output-file string
        Path to the output CSV-file to be written (required).
```

----

## Example

**Input Data**:

<details>

```
rath@gate:~ head ./logs-input.csv 
"timestamp","domain","client_ip","status","fingerprint_ja4","user_agent"
"2025-11-12T13:52:11.105000+01:00","b2466.test.oxl.app","185.132.187.239","301","t13d1812h1_85036bcba153_d41ae481755e","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave/1.63.162 Chrome/122.0.6261.111 Safari/537.36"
"2025-11-12T13:52:11.109000+01:00","14b32.test.oxl.app","109.243.69.100","200","q13d0312h3_55b375c5d22e_151122171f7d","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.1 Safari/605.1.15"
"2025-11-12T13:52:11.144000+01:00","7c8f2.test.oxl.app","213.55.221.232","200","q13d0311h3_55b375c5d22e_f2a83c8e78ae","Python requests someversion"
"2025-11-12T13:52:11.168000+01:00","7c8f2.test.oxl.app","66.249.83.65","200","t13d181300_e8a523a41297_43ade6aba3df","Mozilla/5.0 (X11; Linux x86_64; Storebot-Google/1.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
"2025-11-12T13:52:11.207000+01:00","3ebcc.test.oxl.app","52.205.141.124","200","t12d430700_1ce71f0edbb1_269097225c9f","Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Amazonbot/0.1; +https://developer.amazon.com/support/amazonbot) Chrome/119.0.6045.214 Safari/537.36"
"2025-11-12T13:52:11.222000+01:00","3ebcc.test.oxl.app","114.119.156.140","200","t13d311100_e8f1e7e78f70_d41ae481755e","Mozilla/5.0 (Linux; Android 7.0;) AppleWebKit/537.36 (KHTML, like Gecko) Mobile Safari/537.36 (compatible; PetalBot;+https://webmaster.petalsearch.com/site/petalbot)"
"2025-11-12T13:52:11.234000+01:00","b2466.test.oxl.app","185.132.187.72","200","t13d1812h1_85036bcba153_d41ae481755e","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave/1.63.162 Chrome/122.0.6261.111 Safari/537.36"
"2025-11-12T13:52:11.240000+01:00","eb84e.test.oxl.app","116.204.8.140","200","t13i140900_cbb2034c60b8_e7c285222651","This is a very serious browser!"
"2025-11-12T13:52:11.243000+01:00","17d5b.test.oxl.app","66.249.77.224","200","t13d181300_e8a523a41297_43ade6aba3df","Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
```

</details>

**Output Data**: (*with flags*)

<details>

```
rath@gate:~ head ./logs-output.csv 
timestamp,domain,client_ip,status,fingerprint_ja4,user_agent,flags
2025-11-12T13:52:11.105000+01:00,b2466.test.oxl.app,185.132.187.239,301,t13d1812h1_85036bcba153_d41ae481755e,"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave/1.63.162 Chrome/122.0.6261.111 Safari/537.36",
2025-11-12T13:52:11.109000+01:00,14b32.test.oxl.app,109.243.69.100,200,q13d0312h3_55b375c5d22e_151122171f7d,"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.1 Safari/605.1.15",
2025-11-12T13:52:11.144000+01:00,7c8f2.test.oxl.app,213.55.221.232,200,q13d0311h3_55b375c5d22e_f2a83c8e78ae,Python requests someversion,bot|bot_script
2025-11-12T13:52:11.168000+01:00,7c8f2.test.oxl.app,66.249.83.65,200,t13d181300_e8a523a41297_43ade6aba3df,"Mozilla/5.0 (X11; Linux x86_64; Storebot-Google/1.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",bot|bot_crawler|crawler_user|crawler_verified|fingerprint_crawler|org_google
2025-11-12T13:52:11.207000+01:00,3ebcc.test.oxl.app,52.205.141.124,200,t12d430700_1ce71f0edbb1_269097225c9f,"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Amazonbot/0.1; +https://developer.amazon.com/support/amazonbot) Chrome/119.0.6045.214 Safari/537.36",bot|bot_crawler|crawler_ai_data|org_amazon
2025-11-12T13:52:11.222000+01:00,3ebcc.test.oxl.app,114.119.156.140,200,t13d311100_e8f1e7e78f70_d41ae481755e,"Mozilla/5.0 (Linux; Android 7.0;) AppleWebKit/537.36 (KHTML, like Gecko) Mobile Safari/537.36 (compatible; PetalBot;+https://webmaster.petalsearch.com/site/petalbot)",bot|bot_crawler
2025-11-12T13:52:11.234000+01:00,b2466.test.oxl.app,185.132.187.72,200,t13d1812h1_85036bcba153_d41ae481755e,"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave/1.63.162 Chrome/122.0.6261.111 Safari/537.36",
2025-11-12T13:52:11.240000+01:00,eb84e.test.oxl.app,116.204.8.140,200,t13i140900_cbb2034c60b8_e7c285222651,This is a very serious browser!,bot|bot_script|fingerprint_script
2025-11-12T13:52:11.243000+01:00,17d5b.test.oxl.app,66.249.77.224,200,t13d181300_e8a523a41297_43ade6aba3df,Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html),bot|bot_crawler|crawler_verified|fingerprint_crawler|org_google
```

</details>

----

## Ruleset / Config

The flagging-logic can be found in these files:

* [Match Lists](https://github.com/O-X-L/open-bot-list/blob/latest/log_flagger/internal/config/lists.go)

  These define which match-lists to load and what flags to apply if a request matches it.

* [Enrichment](https://github.com/O-X-L/open-bot-list/blob/latest/log_flagger/internal/enrich.go)

  Here you can see the processing steps:

  * Check if the request has any matches in the open-bot-list config
  * Apply flags specific to match-files
  * Process ruleset and apply flags if rule matches (only one match)
  * Custom fallback flags (*might be migrated to open-bot-list config later on*)

* [Ruleset](https://github.com/O-X-L/open-bot-list/blob/latest/log_flagger/internal/config/ruleset.go)

  This ruleset is the equivalent of the [Open-Bot-List Flagging-Config](https://github.com/O-X-L/open-bot-list/tree/latest/flagging)

  It checks for previously set flags and applies new flags if all required ones exist.

  This allows us to cover much of the flagging-logic required.
