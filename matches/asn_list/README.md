# Open Bot List - ASN Lists

We gather commonly used ASN-Lists from serious providers that can help to categorize traffic.

As we cannot control the contents of these lists - you have to keep some things mind:

* In most cases you might not want to drop traffic because only of an ASN-match.
* It can be dangerous to perform actions only with matching the source-ASN => you should always at least match two factors (*like user-agent*)

**WARNING**: These have usage limitations! Make sure to read their usage policy that is linked in the CSV-file!

## Differentiation

* **Malicious** => ASN-Lists that flag malicious traffic

  File match: `src_asnlist_malicious`
