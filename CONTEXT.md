# fltr

fltr is a content-based HTTP request filter: it inspects the content of every
request, blocks requests that match any configured Block Rule, and forwards
the rest to an Allow Upstream.

## Language

### Rules and matching

**Block Rule**:
A named set of Terms. A rule *matches* a request when every one of its Terms is
found in the request's Searchable Content.
_Avoid_: filter, pattern, regex

**Term**:
A single string a rule searches for within the Searchable Content.
_Avoid_: keyword, word

**Match**:
A rule-level outcome: every Term in the rule appears in the Searchable Content. Matching is binary.
A request either matches a rule or it doesn't, and rules have no rank or
precedence over each other.
_Avoid_: hit, trigger

**Searchable Content**:
The combined text fltr inspects: the request body, the `Title`/`title` and `Message`/`message` header values, and the `Title`/`title` and `Message`/`message` query parameter values.
_Avoid_: payload, combined request text

### Verdicts and routing

**Allowed**:
The verdict for a request that matches no Block Rule. Such a request is
forwarded to the Allow Upstream.
_Avoid_: accepted, let through

**Blocked**:
The verdict for a request that matches at least one Block Rule. Such a request
is forwarded to the Block Upstream when configured. Otherwise it is Discarded.
_Avoid_: rejected, denied

**Allow Upstream**:
The HTTP service that receives Allowed requests.
_Avoid_: default upstream, forward URL

**Block Upstream**:
The optional HTTP service that receives Blocked requests. When not configured,
Blocked requests are Discarded instead.
_Avoid_: drop target, deny upstream

**Discarded**:
A Blocked request with no Block Upstream configured is not forwarded and never
reaches any upstream.
_Avoid_: swallowed, ignored

### Health and liveness

**Liveness Endpoint**:
The Reserved Path `/healthz`, answered by fltr itself with `200` whenever it
can serve HTTP. Liveness is not readiness. The Liveness Endpoint never contacts
an upstream and says nothing about whether requests can be forwarded.
_Avoid_: health check, readiness endpoint

**Reserved Path**:
A path fltr answers itself before content filtering, so it is never inspected,
matched against Block Rules, or forwarded to any upstream.
_Avoid_: special path, internal route

### Forwarding

**Dropped Path**:
The incoming URL path is not forwarded, so requests reach the upstream URL
exactly as configured.
_Avoid_: discarded path, stripped path, ignored path
