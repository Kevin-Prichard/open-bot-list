# HTTP User Agents

All matches should be done case-insensitive!

## Differentiation

* **AI** => Related to AI Search, User-Actions or Training-Data gathering

  File match: `http_user_agent_ai`

* **Crawler** => Organizational crawlers for search engines or cloud-services

  File match: `http_user_agent_crawler`

* **Script** => Default User-Agents HTTP-Client libraries (often used by dumb script-bots) and crawler-libraries without that do not belong to organizations

  File match: `http_user_agent_script`

* **Software** => Bots used by end-user-software (not libraries)

  File match: `http_user_agent_software`

  Some software is more prone to be used for script-bots (or explicitly designed for it) - these are also listed in 'script'

* **Scanner** => End-user-software and libraries built to scan applications for vulnerabilities (or even attack them)

  File match: `http_user_agent_scanner`
