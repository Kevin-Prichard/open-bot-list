# Open Bot List - Flagging

Here we define how you can combine the matches to categorize bot-traffic.

For now, it's an abstract configuration that should reflect the matching logic.

With this ruleset we focus on clean matching/categorizing - you might want to modify some of the matching-rules depending on your organizational needs.

## Categorization

We use **Flags** to categorize matches.

Examples:

* `crawler|crawler_search|org_google` => A crawler, is used for search-engines, the organization is Google
* `crawler|crawler_ai_data|org_google` => A crawler, gathers data for AI training, the organization is Google
* `crawler|crawler_search|crawler_ai_user|crawler_user|org_openai` => A crawler, is used for search-engines and user-initiated AI-search, the organization is OpenAI

----

## Specificality

Some traffic might match multiple times.

1. In multiple categories (here seen as separate files)

  You should match from highest to lowest severity.

  Per example:

  * First process `scanner` matches
  * Then `script-bot` matches
  * Then `ai` matches
  * Afterward the `crawler` matches
  * At last the `implicit` matches if no explicit rule matched

  In most cases you might want to stop processing additional rules after the first one matched.

2. Like `Applebot-Extended` might match:

  * `src_net_crawler_applebot & http_user_agent_applebot`
  * `src_net_crawler_applebot & http_user_agent_applebotextended`

  In this case only the first match should count. The more specific the match, the higher it should be ranked (*and processed earlier*).

----

## Logical Operators

We currently support these matching operators:

* `&` for an AND-condition
* `|` for an OR-condition
* `()` to enclose a condition
