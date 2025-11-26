# Open Bot List - Flagging

Here we define how you can combine the matches to categorize bot-traffic.

For now, it's an abstract configuration that should reflect the matching logic.

With this ruleset we focus on clean matching/categorizing - you might want to modify some of the matching-rules depending on your organizational needs.

## Specificality

Some traffic might match multiple times.

Like `Applebot-Extended` might match:

* `src_net_crawler_applebot & http_user_agent_applebot`
* `src_net_crawler_applebot & http_user_agent_applebotextended`

In that case only the first match should count. The more specific the match, the higher it should be ranked (*and processed earlier*).

## Logical Operators

We currently support these matching operators:

* `&` for an AND-condition
* `|` for an OR-condition
* `()` to enclose a condition
