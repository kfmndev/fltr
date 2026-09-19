# Health endpoint on the proxy port

fltr exposes a liveness probe at `GET /healthz` on the same listener as the proxy
(`FLTR_ADDR`, default `:8080`), answered by a small handler wrapped around the
content filter so the request is never inspected, matched against Block Rules, or
forwarded. We chose the same port over a dedicated health listener because it
requires no extra configuration and the probe verifies the exact socket that
serves traffic; the tradeoff is that a `/healthz` short-circuited ahead of the
filter cannot detect a hung filter chain, which is acceptable because this probe
is liveness only: it asserts that the process can serve HTTP, not that requests
can be forwarded.

## Considered Options

- **Dedicated health port (`FLTR_HEALTH_PORT`).** Rejected: it probes a different
  socket, so it can report healthy while the proxy listener is dead, and it adds
  a port and a configuration variable to document for no additional safety.
