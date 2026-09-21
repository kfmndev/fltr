---
icon: lucide/message-circle-question
---

# FAQ

## Why was my request blocked when I did not expect it?

A [**Block Rule**](../reference/block-rules.md) matches when **every** one of its terms appears anywhere in the [**Searchable Content**](../reference/block-rules.md#matching-semantics). Matching is case-insensitive by default and uses substring search, so a term `password` also matches `PasswordManager`. Set `LOG_LEVEL=debug` to see which rule and terms matched.

## Why did my request reach the upstream at a different path?

The incoming URL path is dropped by design. fltr forwards requests to the upstream URL exactly as configured. See [Forwarding behavior](how-it-works.md#forwarding-behavior).

## Why does my upstream not see the client's IP or original URL?

fltr never sets or forwards `X-Forwarded-For`, `X-Forwarded-Host`, or `X-Forwarded-Proto`, and the reverse proxy strips any such headers a proxy in front of fltr supplied. The upstream sees the hop from fltr, not the original client. See [Forwarding behavior](how-it-works.md#forwarding-behavior).

## What does `200 Request blocked, discarded` mean?

The request matched a Block Rule, and no [**Block Upstream**](../reference/environment-variables.md#upstreams) is configured, so it was discarded. fltr returns `200` with that body instead of forwarding it. Configure `FLTR_BLOCK_UPSTREAM` to send blocked requests somewhere.

## How do I make matching case-sensitive?

Set the environment variable `FLTR_CASE_SENSITIVE` to `true`. Invalid values fall back to the case-insensitive default. See [Environment variables](../reference/environment-variables.md).

## Why was my request rejected with `413 request too large`?

The request body exceeded `FLTR_MAX_BODY_SIZE` (default 10 MB). Raise it if the contents are legitimate. See [HTTP behavior](../reference/http-behavior.md).
