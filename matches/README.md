# Open Bot List - Matches

### Format

The lists are in CSV-format so they can be easily parsed by many systems.

----

## IP/Network Lists

We only use IP-Lists that are published **official** by the providers.    

Valid formats:

* **JSON**

  Selector => JSON-Query to extract the flat list in [RFC 9535](https://jsonpath.com/) and `jq`-cli-tool format

* **CSV** => Comma-separated values

  Selector => Number of the field to extract

* **NLSV** => New-line separated values 

  Plaintext file, empty lines & lines starting with '#' or ';' or '//' are ignored

* **REGEX-JSON** => Extract JSON via regex - then parse via JSON-Query

  Sadly some providers do for some unknown reason not provide an API.. :'(

----

## HTTP User Agents

These matches should only be used if you have no other choice as the client can easily modify its User-Agent.

Sometimes we are required to check them as the organizations to not provide separate official IP/Network Lists for clean categorization.

Example:

* Validate a Google-Bot by the official IP/Network Lists
* Categorize the validated client by matching its User-Agent (`Google-Extended = AI, Googlebot = Search, etc`)

Some User-Agents might match multiple times - only the first match should be used. (top: specific matches => bottom: general matches)
