# Contributing

Thanks for helping improve fltr. This file covers the basics; the [Development page](https://kfmndev.github.io/fltr/internals/development/) documents the codebase layout, testing, docs build, and release flow.

## Getting set up

Install Go and `prek` to build and test locally. The [Development page](https://kfmndev.github.io/fltr/internals/development/) covers the requirements, setup, and testing commands.

## Making changes

- Follow [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `docs:`, ...). The `commit-msg` hook enforces this.
- Documentation changes live under the public docs sources; build the site with `zensical build --clean` to verify.
- For docs, follow the [terminology convention](https://kfmndev.github.io/fltr/internals/development/#terminology), cross-link related pages, and use the glossary (`CONTEXT.md`) for new user-facing pages.

## Submitting

Open a pull request against `main` with a Conventional-Commits-style title. CI runs the build and tests, static analysis, the link check, and the docs build; keep it green and a maintainer will take it from there.

## Reporting issues

- Bugs: open an issue with steps to reproduce, expected vs actual behavior, and your configuration (redact anything sensitive).
- Vulnerabilities: do **not** open a public issue; see [SECURITY.md](SECURITY.md).

## License

By contributing, you agree that your contributions are licensed under the project's [license](LICENSE.md).
