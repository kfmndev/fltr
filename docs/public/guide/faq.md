---
icon: lucide/message-circle-question
---

# FAQ

## Why was my request blocked when I did not expect it?

A [**Block Rule**](../reference/block-rules.md) matches when **every** one of its terms appears anywhere in the [**Searchable Content**](../reference/block-rules.md#matching-semantics). Matching is case-insensitive by default and uses substring search, so a term `password` also matches `PasswordManager`. Set `LOG_LEVEL=debug` to see which rule and terms matched.

## Why did my request reach the upstream at a different path?

The incoming URL path is dropped by design. See [Forwarding](../reference/http-behavior.md#forwarding).

## Why does my upstream not see the client's IP or original URL?

fltr never sets or forwards `X-Forwarded-*` headers. See [Forwarding](../reference/http-behavior.md#forwarding).

## What does `200 Request blocked, discarded` mean?

The request matched a Block Rule and no [**Block Upstream**](../reference/environment-variables.md#upstreams) is configured, so it was discarded. Configure `FLTR_BLOCK_UPSTREAM` to send blocked requests somewhere instead.

## How do I make matching case-sensitive?

Set `FLTR_CASE_SENSITIVE` to `true`. See [Environment variables](../reference/environment-variables.md).

## Why was my request rejected with `413 request too large`?

The request body exceeded `FLTR_MAX_BODY_SIZE`. See [HTTP behavior](../reference/http-behavior.md).
