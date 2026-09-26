# Liveness Endpoint on the proxy port

fltr exposes its Liveness Endpoint at `GET /healthz` on the same listener as the
proxy (`FLTR_ADDR`, default `:8080`). A small handler wraps the content filter,
so the request is never inspected, matched against Block Rules, or forwarded. We
chose the same port over a dedicated listener (`FLTR_HEALTH_PORT`) because it
requires no extra configuration and the Liveness Endpoint verifies the exact
socket that serves traffic. The tradeoff is that a request short-circuited ahead
of the filter cannot detect a hung filter chain. That is acceptable because
liveness is not readiness. It asserts that the process can serve HTTP, not that
requests can be forwarded.

## Considered options

- **Dedicated health port (`FLTR_HEALTH_PORT`).** Rejected: it probes a different
  socket, so it can report healthy while the proxy listener is dead, and it adds
  a port and a configuration variable to document for no additional safety.
