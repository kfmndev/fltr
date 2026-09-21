---
icon: lucide/filter
---

# Block Rules

The [**Block Rules**](glossary.md#block-rule) file is a JSON object. Its keys name Block Rules, and each value is a non-empty array of non-empty terms:

```json title="block_rules.json"
{
    "credentials": ["password", "secret"],
    "identity": ["private-key"]
}
```

With the example above, a request is blocked when it contains both `password` and `secret`, or when it contains `private-key`, anywhere in its [**Searchable Content**](glossary.md#searchable-content).

## Matching semantics

- A request matches a Block Rule when **every term** of that rule is found in the Searchable Content. A request matching no Block Rule is allowed.
- The Searchable Content is the request body, the `Title`/`Message` headers, and the `Title`/`title` and `Message`/`message` query parameters. Header names are matched case-insensitively.
- Matching is case-insensitive by default: both the Searchable Content and the terms are lowercased before matching. Set `FLTR_CASE_SENSITIVE=true` for case-sensitive matching.
- Nothing is matched in the URL path, only the body, headers, and query parameters listed above.

The default file location is `block_rules.json` in the current directory; `FLTR_BLOCK_RULES_FILE` points to an alternative path. Terms are trimmed when the file is loaded.

!!! danger "Caution"
    The service refuses to start if any Block Rule has no terms or contains a blank term after trimming.

## Try a request

With the example configuration running, send a request that matches the `credentials` Block Rule:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -H 'Title: account details' \
    -d 'password secret'
```

Send a request that matches no Block Rule to route it to the [**Allow Upstream**](glossary.md#allow-upstream):

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -d 'hello service'
```
