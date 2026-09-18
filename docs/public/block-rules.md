---
icon: lucide/filter
---

# Block rules

The block rules file is a JSON object. Its keys name rules, and each value is a non-empty array of non-empty terms:

```json
{
    "credentials": ["password", "secret"],
    "identity": ["private-key"]
}
```

With the example above, a request is blocked when it contains both `password` and `secret`, or when it contains `private-key`, anywhere in its searchable content.

## Matching semantics

- A request matches a rule when **every term** of that rule is found in the searchable content. A request matching no rule is allowed.
- The searchable content is the request body, the `Title`/`Message` header values, and the `title`/`message` query parameter values.
- Matching is case-insensitive by default: both the searchable content and the terms are lowercased before matching. Set `FLTR_CASE_SENSITIVE=true` for case-sensitive matching.
- Nothing is matched in the URL path, only the body, headers, and query parameters listed above.

The default file location is `block_rules.json` in the current directory; `FLTR_BLOCK_RULES_FILE` points to an alternative path. Terms are trimmed when the file is loaded.

!!! danger "Caution"

    The service refuses to start if any rule has no terms or contains a blank term after trimming.

## Try a request

With the example configuration running, send a request that matches the `credentials` rule:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -H 'Title: account details' \
    -d 'password secret'
```

Send a request that matches no rule to route it to the allow upstream:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -d 'hello service'
```
