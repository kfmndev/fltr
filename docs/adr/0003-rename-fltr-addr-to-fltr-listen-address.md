# Rename `FLTR_ADDR` to `FLTR_LISTEN_ADDRESS`

The listen-address setting's environment variable is renamed from `FLTR_ADDR`
to `FLTR_LISTEN_ADDRESS`, matching its canonical flag `--listen-address`. The
old name is removed outright: there is no deprecated alias and no compatibility
fallback. This is a deliberate breaking change, accepted because the
application is still young and the inconsistency between `FLTR_ADDR` and
`--listen-address` would otherwise be permanent. The cleaner alternative,
keeping `FLTR_ADDR` working while warning, was rejected as carrying a second
name and a deprecation path with no real deployment base to protect.

## Consequences

- Existing configurations that set `FLTR_ADDR` stop taking effect; they must be
  updated to `FLTR_LISTEN_ADDRESS`.
- Documentation, Docker defaults, and the `.env` example must use the new name.
