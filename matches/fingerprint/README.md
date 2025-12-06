# Open Bot List - Fingerprints

**ToDo**: Check all JA4-fingerprints for false-positives.

Be aware that multiple different crawlers and bots can match the same JA4 TLS-fingerprint. That mainly is the case if they use the same underlying client-library.

## Differentiation

* **Crawler** => Organizational crawlers for search engines or cloud-services

  File match: `fingerprint_crawler`

  Kind match for JA4: `fingerprint_tls_ja4_crawler`

* **Script** => Default User-Agents HTTP-Client libraries (often used by dumb script-bots) and crawler-libraries without that do not belong to organizations

  File match: `fingerprint_script`

  Kind match for JA4: `fingerprint_tls_ja4_script`

* **Scanner** => End-user-software and libraries built to scan applications for vulnerabilities (or even attack them)

  File match: `fingerprint_scanner`

  Kind match for JA4: `fingerprint_tls_ja4_scanner`
