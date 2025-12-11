# Open Bot List - ASN's

Information base-source: [OXL IP-Abuse Database Lists](https://github.com/O-X-L/risk-db-lists/tree/main/asn)

## Guidelines

ASN-Flagging can be tough because one organization can provide many services and only use one ASN.

An ASN can be defined as:

* **One Kind**:
  * **ISP**
    (*provides Internet-Access to end-users and businesses*)
  * **Hosting**
    (*server hosting for 3rd-parties, VPS, Serverless functions like Cloud-provider*)
  * **Education**
    (*Universities and so on*)

* **Multiple Services**:
  These have to be major services. Optimally a service that the ASN-owner-organization (*ASN-org*) itself provides.

  * **Cloud**
    (*if more than just server hosting/VPS*)
  * **CDN**
    (*the ASN-org provides a CDN service*)
  * **VPN**
    (*only if the ASN-org itself provides a VPN or large 'players' are using this ASN - some end-users running a VPN on a hosting-provider is not a valid*)
  * **Proxy**
    (*same as VPN - only if the ASN is specialized in it*)
  * **Crawler**
    (*the ASN-org itself provides services that crawl websites to gather information - examples: Search engine, Social Media*)
  * **Scanner**
    (*the ASN-org itself provides services that probe/scan networks and hosts to gather information - example: Attack surface checker*)

----

## How to Use


You can look-up the ASN of IPs by using a GeoIP-Database: (*offline recommended for performance*)

* Our [GeoIP-ASN Database](https://github.com/O-X-L/geoip-asn)
* [IPInfo](https://ipinfo.io/)
* [MaxMind](https://maxmind.com)

If your system lacks GeoIP-lookup capabilities - you are able to translate AS-numbers to network-ranges via:

* Our [GeoIP-ASN Database](https://github.com/O-X-L/geoip-asn) ([API](https://geoip.oxl.app))

* The `whois` cli-tool: `whois -h whois.radb.net -- '-i origin AS<NUMBER>' | grep ^route | awk '{gsub("(route:|route6:)","");print}' | awk '{gsub(/ /,""); print}'`
