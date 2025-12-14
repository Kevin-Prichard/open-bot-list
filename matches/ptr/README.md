# Open Bot List - PTR's

Some organizations only provide us one or more PTR's (*reverse DNS*) that can be used to validate if an IP-address belongs to them.

Example: perform a PTR-lookup for the IP-address `1.1.1.1` ([toolbox.googleapps.com/apps/dig](https://toolbox.googleapps.com/apps/dig/#PTR/))

These DNS-lookups should be performed asynchronously as the query can take relatively long.

For now, we will not cover how you can integrate those PTR-checks in your WAF/proxy.

The [Log-Flagger application](https://github.com/O-X-L/open-bot-list/blob/latest/apps/log_flagger/README.md) can query the PTR's for existing logs.
