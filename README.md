# SHIELD-WALL WAF Lists

<p align="center">
    <a title="Support this Project (Donate, Support-Licenses)" href="https://shop.oxl.app/collections/open-source">
        <img src="https://files.oxl.at/img/badge-oss-support.svg" alt="Support Badge (Donate, Support-Licenses)"/>
    </a>
</p>

This repository is used to collect information that can be used to categorize & match traffic.

----

## SHIELD-WALL WAF Project

Check-out the main [SHIELD-WALL WAF Project](https://github.com/O-X-L/shieldwall-waf) | [demo.waf.shield-wall.net](https://demo.waf.shield-wall.net)

---

## Categorization

We use **Flags** to categorize matches.

Examples:

* `crawler|crawler_search|org_google` => A crawler, is used for search-engines, the organization is Google
* `crawler|crawler_ai_data|org_google` => A crawler, gathers data for AI training, the organization is Google
* `crawler|crawler_search|crawler_ai_search|crawler_user|org_openai` => A crawler, is used for search-engines and user-initiated AI-search, the organization is OpenAI

### Format

The lists are in CSV-format so they can be easily parsed by many systems.

----

## IP/Network Lists

We only use IP-Lists that are published **official** by the providers.    

Valid formats:

* **JSON**

  Selector => JSON-Query to extract the flat list in [RFC 9535](https://jsonpath.com/) and `jq`-cli-tool format

* **CSV** => Comma-separated values

  Selector => Number of the field to extract

* **NLSV** => New-line separated values 

  Plaintext file, empty lines & lines starting with '#' or ';' or '//' are ignored

* **HTML** => Embedded inside HTML

  Sadly some providers do for some unknown reason not provide an API.. :'(

----

## HTTP User Agents

These matches should only be used if you have no other choice as the client can easily modify its User-Agent.

Sometimes we are required to check them as the organizations to not provide separate official IP/Network Lists for clean categorization.

Example:

* Validate a Google-Bot by the official IP/Network Lists
* Categorize the validated client by matching its User-Agent (`Google-Extended = AI, Googlebot = Search, etc`)

Some User-Agents might match multiple times - only the first match should be used. (top: specific matches => bottom: general matches)
