# fltr

fltr is a content-based HTTP request filter: it inspects the content of every
request, blocks requests that match any configured block rule, and forwards
the rest to an allow upstream.

## Language

### Rules and matching

**Block Rule**:
A named set of terms. A rule *matches* a request when every one of its terms is
found in the request's searchable content.
_Avoid_: filter, pattern, regex

**Term**:
A single string a rule searches for within the searchable content.
_Avoid_: keyword, word

**Match**:
A rule-level outcome: every term in the rule appears in the searchable content.
Matching is binary — a request either matches a rule or it doesn't, and rules
have no rank or precedence over each other.
_Avoid_: hit, trigger

**Searchable Content**:
The combined text fltr inspects: the request body, the `Title` and `Message`
header values, and the `title` and `message` query parameter values.
_Avoid_: payload, combined request text

### Verdicts and routing

**Allowed**:
The verdict for a request that matches no block rule. Allowed requests are
forwarded to the allow upstream.
_Avoid_: accepted, let through

**Blocked**:
The verdict for a request that matches at least one block rule. Blocked
requests are forwarded to the block upstream when configured, otherwise they
are discarded.
_Avoid_: rejected, denied

**Allow Upstream**:
The HTTP service that receives requests that are allowed.
_Avoid_: default upstream, forward URL

**Block Upstream**:
The optional HTTP service that receives requests that are blocked. When not
configured, blocked requests are discarded instead.
_Avoid_: drop target, deny upstream

**Discarded**:
The fate of a blocked request when no block upstream is configured: it is
dropped and never reaches any upstream.
_Avoid_: dropped, swallowed
