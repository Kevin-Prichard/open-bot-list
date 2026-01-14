# Open Bot List - IP/Network Lists

We gather official IP/Network-Lists that can be used to identify bots.

As we cannot control the contents of these lists - you have to keep some things mind:

* It can be dangerous to perform actions only with matching the source-IP => you should always at least match two factors (*like user-agent*)
* It can make sense to run some sanity-checks once downloaded (*p.e. they do not include internal networks*)

## Differentiation

* **AI** => Related to AI Search, User-Actions or Training-Data gathering

  File match: `src_net_ai`

* **CDN** => Used by CDN proxies

  File match: `src_net_cdn`

* **Crawler** => Organizational crawlers for search engines or cloud-services

  File match: `src_net_crawler`

* **eCommerce** => Organizational crawlers for eCommerce-related services

  File match: `src_net_ecommerce`

* **Malicious** => IP-Lists that flag malicious traffic

  **WARNING**: These have usage limitations! Make sure to read their usage policy that is linked in the CSV-file!

  File match: `src_net_malicious`

* **Monitoring** => Organizational crawlers for Monitoring-related services

  File match: `src_net_monitoring`

* **Proxy** => IPs and Networks officially used by/for privacy-proxies (*like Tor-network*)

  File match: `src_net_proxy`

* **VPN** => IPs and Networks officially used by/for VPNs

  File match: `src_net_vpn`
