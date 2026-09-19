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
not forwarded and never reaches any upstream.
_Avoid_: swallowed, ignored

### Health and liveness

**Liveness Endpoint**:
The reserved path `/healthz`, answered by fltr itself with `200` whenever it
can serve HTTP. Liveness is not readiness: the Liveness Endpoint never contacts
an upstream and says nothing about whether requests can be forwarded.
_Avoid_: health check, readiness endpoint

**Reserved Path**:
A path fltr answers itself before content filtering, so it is never inspected,
matched against Block Rules, or forwarded to any upstream.
_Avoid_: special path, internal route

### Forwarding

**Dropped Path**:
The fate of the incoming URL path: fltr does not forward it, so requests reach
the upstream URL exactly as configured.
_Avoid_: discarded path, stripped path, ignored path
