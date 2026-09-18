# Contributing

Thanks for helping improve fltr. This file covers the basics; the [Development page](https://kfmndev.github.io/fltr/development/) documents the codebase layout, testing, docs build, and release flow.

## Requirements

- [Go 1.25](https://go.dev) or newer
- [prek](https://github.com/j178/prek)

## Getting set up

```sh
go build ./...
prek install --hook-type pre-commit --hook-type commit-msg
```

The hooks run formatting, `go mod tidy`, `go vet`, a compile check, and the test suite on every commit; the `commit-msg` hook enforces Conventional Commits.

## Making changes

- Build and test locally before pushing:

  ```sh
  go test ./...
  ```

- Follow [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `docs:`, ...).
- Documentation changes live under the public docs sources; build the site with `zensical build --clean` to verify.
- Link every behavior you document, and keep new user-facing pages using the project's glossary (see `CONTEXT.md`).

## Submitting

Open a pull request against `main` with a Conventional-Commits-style title. CI runs the build and tests, static analysis, and the link check; keep it green and a maintainer will take it from there.

## Reporting issues

- Bugs: open an issue with steps to reproduce, expected vs actual behavior, and your configuration (redact anything sensitive).
- Vulnerabilities: do **not** open a public issue — see [SECURITY.md](SECURITY.md).

## License

By contributing, you agree that your contributions are licensed under the project's [license](LICENSE.md).
