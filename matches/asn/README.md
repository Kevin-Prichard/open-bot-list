# Open Bot List - ASN's

For now see: [OXL IP-Abuse Database Lists](https://github.com/O-X-L/risk-db-lists/tree/main/asn)

You can look-up the ASN of IPs by using a GeoIP-Database: (*offline recommended for performance*)

* Our [GeoIP-ASN Database](https://github.com/O-X-L/geoip-asn)
* [IPInfo](https://ipinfo.io/)
* [MaxMind](https://maxmind.com)

If your system lacks GeoIP-lookup capabilities - you are able to translate AS-numbers to network-ranges via:

* Our [GeoIP-ASN Database](https://github.com/O-X-L/geoip-asn) ([API](https://geoip.oxl.app))

* The `whois` cli-tool: `whois -h whois.radb.net -- '-i origin AS<NUMBER>' | grep ^route | awk '{gsub("(route:|route6:)","");print}' | awk '{gsub(/ /,""); print}'`
