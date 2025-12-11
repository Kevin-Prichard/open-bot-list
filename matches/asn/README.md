# Open Bot List - ASN's

Information base-source: [OXL IP-Abuse Database Lists](https://github.com/O-X-L/risk-db-lists/tree/main/asn)

----

## Differentiation

ASN-Flagging can be tough because one organization can provide many services and only use one ASN.

### Main Kind

An ASN can only be flagged as one of these:

* **Hosting** => server hosting for 3rd-parties, VPS, Serverless functions like Cloud-provider

  File match: `src_asn_hosting`

* **ISP** => provides Internet-Access to end-users and businesses

  File match: `src_asn_isp`

* **Education** => Universities and so on

  File match: `src_asn_edu`

### Services

These have to be major services.

Optimally a service that the ASN-owner-organization (*ASN-org*) itself provides.

* **Cloud** => the ASN-org does more than just server hosting/VPS

  File match: `src_asn_cloud`

* **CDN** => the ASN-org provides a CDN service

  File match: `src_asn_cdn`

* **VPN** => only if the ASN-org itself provides a VPN or large 'players' are using this ASN - some end-users running a VPN on a hosting-provider is not a valid

  File match: `src_asn_vpn`

* **Proxy** => same as VPN - only if the ASN is specialized in it

  File match: `src_asn_proxy`

* **Crawler** => the ASN-org itself provides services that crawl websites to gather information - examples: Search engine, Social Media

  File match: `src_asn_crawler`

* **Scanner** => the ASN-org itself provides services that probe/scan networks and hosts to gather information - example: Attack surface checker

  File match: `src_asn_scanner`

----

## How to Use


You can lookup the ASN of IPs by using a GeoIP-Database: (*offline recommended for performance*)

* Our [GeoIP-ASN Database](https://github.com/O-X-L/geoip-asn)
* [IPInfo](https://ipinfo.io/)
* [MaxMind](https://maxmind.com)

If your system lacks GeoIP-lookup capabilities - you are able to translate AS-numbers to network-ranges via:

* Our [GeoIP-ASN Database](https://github.com/O-X-L/geoip-asn) ([API](https://geoip.oxl.app))

* The `whois` cli-tool:

  ```bash
  sudo apt install whois
  whois -h whois.radb.net -- '-i origin AS<NUMBER>' | grep -E '^(route:|route6:)' | tr -d ' ' | cut -d ':' -f2-
  ```
