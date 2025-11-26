# Open Bot List

<p align="center">
    <a title="Support this Project (Donate, Support-Licenses)" href="https://shop.oxl.app/collections/open-source">
        <img src="https://files.oxl.at/img/badge-oss-support.svg" alt="Support Badge (Donate, Support-Licenses)"/>
    </a>
</p>

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
  * Separating different bot-categories like `script bots`, `hidden bots`, `search-engine crawlers`, `AI-data crawlers`, `AI-user crawlers`, `social-media crawlers`, `crawlers for ADs`, `crawlers for ecommerce`, ...
    * Matching clear script-bots by their User-Agent (*dumb script-kiddies*)
    * Matching 'hidden' bots by their client-fingerprints (*[JA4](https://github.com/O-X-L/haproxy-ja4-fingerprint), etc.*)
    * ... *to be extended* ...

* **PTR-checks**
  * Some organizations only supply us with a PTR-match to validate if a crawler-IP is theirs (*no simple IP-list lookups*)

* **Traffic Flagging**
  * We provide you with abstract configuration that shows how the matches can be combined
  * Practical configuration examples for proxy-services will be added later on

----

## Motivation

We are working on building a [FOSS WAF-platform](https://github.com/O-X-L/shieldwall-waf) (*and centrally manageable network-firewalls*) which require such a collection of bot-related information.

With our [IP-Abuse Reporting-System & Databases](https://github.com/O-X-L/risk-db) we have already started to collect information for it.

As the mindset of Open-Source is at the core of our being - we want to transparently share it with the whole world.

----

### SHIELD-WALL WAF Project

This information-collection is part of our [SHIELD-WALL WAF Project](https://github.com/O-X-L/shieldwall-waf).

Check-out the demo: [demo.waf.shield-wall.net](https://demo.waf.shield-wall.net)
