---
icon: lucide/route
---

# How it works

**fltr** is a content-based HTTP request filter. It inspects the [**Searchable Content**](block-rules.md#matching-semantics) of every request, matches it against the [**Block Rules**](block-rules.md), and routes the request by its verdict. [**Matching**](glossary.md#match) is binary. The Block Rules have no rank or precedence over each other.

```mermaid
flowchart LR
    R(<b>📦 Request</b>) --> BR(<b>📋 Block Rules</b>)
    BR -- No Rules match ---> A(<b>✅ Allowed</b>)
    BR -- Any Rule matches ---> B(<b>❌ Blocked</b>)

    classDef blue fill:#4287ff1F,stroke:#4287ff
    classDef green fill:#00c7531F,stroke:#00c753
    classDef red fill:#ff1a471F,stroke:#ff1a47
    class BR blue
    class A green
    class B red
```

## Startup

Before fltr starts listening, it:

1. Reads the configuration and validates each upstream URL (only `http` and `https` with a hostname are accepted).
2. Sends an HTTP `HEAD` request to each configured upstream. Any HTTP response counts as reachable, including error status codes. Network and TLS failures abort the startup.
3. Loads the Block Rules from `block_rules.json`, or the file named by `FLTR_BLOCK_RULES_FILE`. A Block Rule with no terms, or a blank term after trimming, aborts startup.

## Request flow

For each incoming request:

```mermaid
flowchart TD
    A(<b>Incoming request</b>)
    B(<b>Liveness Check</b><br/>GET / HEAD → 200<br/>Other methods → 405)
    A -- Reserved Path /healthz? --> B
    A -- Other Request --> C(<b>Read</b><br>Body ≤ <i>FLTR_MAX_BODY_SIZE</i>?)
    C -- yes --> D(<b>Assemble</b><br>Searchable Content)
    C -- no --> E(<b>Rejected</b><br>413 request too large)
    D --> F(<b>Match</b><br>Block Rules)
    F -- No Rules match --> G(<b>Allowed</b><br>Forward to Allow Upstream)
    F -- Any Rule matches --> H(<i>FLTR_BLOCK_UPSTREAM</i> set?)
    H -- yes --> I(<b>Blocked</b><br>Forward to Block Upstream)
    H -- no --> J(<b>Discarded</b><br/>200 Request blocked, discarded)

    classDef grey fill:#ffffff1F,stroke:#ffffff
    classDef red fill:#ff1a471F,stroke:#ff1a47
    classDef green fill:#00c7531F,stroke:#00c753
    classDef blue fill:#4287ff1F,stroke:#4287ff
    class B grey
    class C,D,F,H blue
    class E,I,J red
    class G green
```

1. **Health**: a request to the [**Reserved Path**](glossary.md#reserved-path) `/healthz` is answered by fltr itself and never reaches the filter or an upstream. `GET` and `HEAD` return `200`; any other method returns `405`.
2. **Read**: the body is read up to `FLTR_MAX_BODY_SIZE` (default 10 MB). A larger body is rejected with `413 request too large`.
3. **Assemble**: the Searchable Content is built from the request body, the `Title`/`Message` headers, and the `Title`/`title` and `Message`/`message` query parameters. Nothing in the URL path is inspected.
4. **Match**: fltr checks each Block Rule. A rule *matches* when every term of that rule is found in the Searchable Content. A request matching any Block Rule is blocked. A request matching none is allowed.
5. **Route** by verdict:
    - **Allowed**: forwarded to the [**Allow Upstream**](environment-variables.md#upstreams).
    - **Blocked**: forwarded to the [**Block Upstream**](environment-variables.md#upstreams) when `FLTR_BLOCK_UPSTREAM` is configured. Otherwise discarded, and fltr responds with `200 Request blocked, discarded`.

See [HTTP behavior](http-behavior.md#responses) for a table containing all possible responses.

## Forwarding behavior

fltr forwards the method, end-to-end headers, and body, **drops** the incoming URL path, and combines the upstream's query parameters with the incoming ones. It never adds `X-Forwarded-*` headers. See [HTTP behavior](http-behavior.md#forwarding) for more detail.

## Where each variable fits

- `FLTR_BLOCK_RULES_FILE`: startup (rules loading).
- `FLTR_MAX_BODY_SIZE`: step 2 (body reading).
- `FLTR_CASE_SENSITIVE`: step 4 (matching).
- `FLTR_ALLOW_UPSTREAM` and `FLTR_BLOCK_UPSTREAM`: step 5 (routing). See [Startup](#startup) for reachability semantics.
- `LOG_LEVEL` and `LOG_FORMAT`: logging only; they do not affect filtering or routing.
