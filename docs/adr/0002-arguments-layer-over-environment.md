# Arguments layer over environment variables

fltr resolves every setting from three sources in a fixed order: an explicitly
set command-line argument wins, otherwise a non-empty environment variable,
otherwise the built-in default. Environment variables remain a first-class,
fully supported source — arguments *layer over* them rather than replacing
them, and no environment variable is deprecated in favour of a flag. This lets
the same binary serve two postures: flags for interactive and ad-hoc runs, and
environment variables for containers and orchestrators, where Docker and the
s6 service keep configuring fltr entirely through the environment.

The tradeoff is two configuration surfaces to document and keep in agreement,
and a precedence rule every setting must honour. We chose it because the
alternative — arguments or environment, not both — would strand either existing
deployments or the interactive use case. A configuration *file* is deliberately
not a source: it would add a format, a mount point, and a third precedence
layer for no benefit the two existing sources do not already cover.

## Considered Options

- **Arguments replace environment variables.** Rejected: it would break Docker
  and existing deployments, which have no natural place to pass arguments.
- **Environment variables win over arguments.** Rejected: the flag is the more
  specific, per-invocation input, and users expect the explicit flag to win.
- **A configuration file as a third source.** Rejected for now: a separate
  effort, and not needed while environment variables already cover deployments.

## Consequences

- Every new setting must be resolvable from both sources, or explicitly
  documented as one-source-only.
- Invalid values fail loudly from either source; only setting neither source
  falls back to the default.
