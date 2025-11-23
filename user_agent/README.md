# HTTP User Agents

Differentiation:

* **AI** => Related to AI Search, User-Actions or Training-Data gathering

  Default flag: `crawler_ai_data` or `crawler_ai_user`

* **Crawler** => Organizational crawlers for search engines or cloud-services

  Default flag: `crawler`

* **Script** => Default User-Agents HTTP-Client libraries (often used by dumb script-bots)

  Default flag: `bot_script`

* **Software** => Bots used by end-user-software (not libraries)

  Default flag: `bot_software`

  Some software is more prone to be used for script-bots (or explicitly designed for it) - these will also be flagged as 'bot_script'
