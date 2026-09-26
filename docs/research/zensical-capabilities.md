# Zensical capabilities research

**Ticket:** kfmndev/fltr#8 ("Zensical capabilities research"), part of map #7
**Date:** 2026-09-17
**Sources:** verified against zensical.org/docs (current stable release 0.0.62; the site banner advertises an upcoming 0.1.0 launch on Nov 5, 2026, alongside Studio/Spark and the Material for MkDocs EOL extension to May 5, 2027), not Material for MkDocs docs. Cross-checked against the Zensical backlog/GitHub where the docs are silent.

## 1. Admonitions / callouts

Source: https://zensical.org/docs/authoring/admonitions/

Enable in `zensical.toml`:

```toml
[project.markdown_extensions]
admonition = {}
pymdownx.details = {}       # collapsible blocks
pymdownx.superfences = {}   # nesting arbitrary content
```

### Syntax

- Basic: `!!! <type>` with content indented 4 spaces.
- Custom title: `!!! note "Custom title"`; remove title: `!!! note ""` (not allowed for collapsible).
- **Collapsible:** `??? <type>` (collapsed) / `???+ <type>` (initially expanded). Requires `pymdownx.details`.
- Inline blocks (sidebar): `!!! info inline end "..."` (right) / `!!! info inline "..."` (left); must be declared before the content block.
- Nesting works when `pymdownx.superfences` is enabled.
- Per-type icons via `[project.theme.icon.admonition]` (`<type> = "<icon>"`); 12 types each have a distinct default icon (Octicons/FontAwesome sets selectable).

### Supported admonition types (12)

`note` (default), `abstract`, `info`, `tip`, `success`, `question`, `warning`, `failure`, `danger`, `bug`, `example`, `quote`.

### GitHub callouts (`> [!NOTE]` etc.)

Fully supported via the bundled `pymdownx.quotes` extension with `callouts = true`:

```toml
[project.markdown_extensions.pymdownx.quotes]
callouts = true
```

Marker must be ALL-UPPERCASE (same constraint as GitHub). With this extension, any of the 12 admonition types works as a callout, not just GitHub's five.

### Mapping of GitHub banners to zensical admonition types

| GitHub banner | Zensical type | Notes |
| --- | --- | --- |
| `[!NOTE]` | `note` | direct mapping, styled |
| `[!TIP]` | `tip` | direct mapping, styled |
| `[!WARNING]` | `warning` | direct mapping, styled |
| `[!IMPORTANT]` | `important` | **not** a built-in admonition type; renders with default admonition styling unless custom CSS is added |
| `[!CAUTION]` | `caution` | **not** a built-in admonition type; renders with default admonition styling unless custom CSS is added |

Practical alternative for full styling: rewrite `[!IMPORTANT]` → `!!! abstract` (or `info`) and `[!CAUTION]` → `!!! danger` in the docs, or add small CSS for the two missing types.

## 2. Theming

### Theme variants

Source: https://zensical.org/docs/compatibility/mkdocs/migration/#theme-variant

- Two variants: `modern` (new default look) and `classic` (preserves Material for MkDocs appearance). Same HTML structure in both. Set via `theme.variant = "classic"` (mkdocs.yml) / variant setting in `zensical.toml`.
- `--theme` CLI flag is not supported; use the variant setting.

### Color schemes (dark/light)

Source: https://zensical.org/docs/setup/colors/

- Two built-in schemes: `default` (light) and `slate` (dark). These are the Material equivalents.
- Single scheme: `[project.theme.palette] scheme = "default"`.
- Light/dark toggle: define `[[project.theme.palette]]` as a **list** of palettes, each with `scheme`, `toggle.icon` (e.g. `lucide/sun` / `lucide/moon`), `toggle.name`. Optional `media = "(prefers-color-scheme: dark)"` for system-preference-following, including a 3-toggle automatic light/dark mode with `media = "(prefers-color-scheme)"`.
- Primary/accent colors: `palette.primary` / `palette.accent` with named colors (red … white); `primary = "custom"` + CSS variables (`--md-primary-fg-color` etc.) in an extra stylesheet. Custom named schemes via `[data-md-color-scheme="..."]` CSS + `extra_css`. Slate hue tunable via `--md-hue`.
- Per-palette distinct primary/accent colors supported.

### Fonts

Source: https://zensical.org/docs/setup/fonts/

- `[project.theme] font.text = "Inter"` (regular) and `font.code = "JetBrains Mono"` (mono); any Google Font, auto-loaded.
- `font = false` disables Google Fonts loading (privacy / system fallback). Custom `@font-face` via extra stylesheet; apply with `--md-text-font` / `--md-code-font` CSS variables.

## 3. Project layout & exclude mechanism

### Layout: confirmed

Sources: https://zensical.org/docs/create-your-site/ , https://zensical.org/docs/usage/new/

`zensical new .` creates exactly:

```
.
├─ .github/workflows
│  └─ docs.yml
├─ docs/
│  ├─ index.md
│  └─ markdown.md
└─ zensical.toml
```

Root `zensical.toml` + `docs/` docs_dir + `.github/workflows/docs.yml`. This matches the ticket's assumption. `docs_dir` and `site_dir` are configurable ([setup/basics](https://zensical.org/docs/setup/basics/)); `docs_dir` cannot be `.` (temporary limitation, backlog #101). `site_name` is the only required setting. Zensical also still reads `mkdocs.yml` if you'd rather not migrate yet.

### Exclude mechanism: **NOT YET SUPPORTED** ⚠️

- MkDocs's `exclude_docs` (and `draft_docs`) settings are explicitly listed as **not yet supported** in Zensical: https://zensical.org/docs/compatibility/mkdocs/migration/#unsupported-settings (tracked in https://github.com/zensical/backlog/issues/65).
- Maintainer statement on https://github.com/zensical/zensical/issues/135 : both settings are unsupported; regarded as **low priority**; they'll be tackled during a general config/editing-experience rework.
- The roadmap lists an **"MkDocs `exclude` plugin replacement"** under "Up next" (https://zensical.org/roadmap/#mkdocs-compatibility), planned but not shipped.
- **No documented fallback or exact syntax exists today.** Workarounds for keeping `docs/agents/` out of the public build:
  - Move agent docs out of `docs_dir` entirely (e.g. `docs_dir = "docs/public"`, keep `docs/agents/` as a sibling not under docs_dir), or
  - Keep using `mkdocs.yml` + MkDocs (which supports `exclude_docs`) while building with Zensical only once the feature lands, or
  - A CI step that copies/removes paths into a clean staging dir before `zensical build` (hacky; no first-class option).

## 4. Versioning: transitional only, no native support ⚠️

Source: https://zensical.org/docs/compatibility/mkdocs/mike/

- **No native versioned-docs deployment yet.** Zensical ships a **transitional fork of mike** (squidfunk/mike, based on mike 2.2.0; compatibility fixes only, no new features) installed via `pip install git+https://github.com/squidfunk/mike.git`.
- Existing `extra.version.provider: mike` config and the Material versioning guide settings keep working, including the version selector.
- Native versioning is on the roadmap (https://zensical.org/roadmap/#versioning, covering git-based versions, folder-based versions, selective publishing, and focused rebuilds) with no delivery date. The mike fork is maintained "until Zensical provides native versioning support".
- For fltr this means: versioned docs today = mike + git orphan-branch mechanics; or defer and wait for native support.

## 5. Publishing: official GitHub Actions workflow

Source: https://zensical.org/docs/publish-your-site/#with-github-actions (workflow generated by `zensical new` into `.github/workflows/docs.yml`)

```yaml
name: Documentation
on:
  push:
    branches:
      - master
      - main
permissions:
  contents: read
  pages: write
  id-token: write
jobs:
  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/configure-pages@v6
      - uses: actions/checkout@v7
      - uses: actions/setup-python@v6
        with:
          python-version: 3.x
      - run: pip install zensical
      - run: zensical build --clean
      - uses: actions/upload-pages-artifact@v5
        with:
          path: site
      - uses: actions/deploy-pages@v5
        id: deployment
```

Key facts:

- **Trigger:** push to `master` or `main`.
- **Actions (major versions):** `actions/configure-pages@v6`, `actions/checkout@v7`, `actions/setup-python@v6`, `actions/upload-pages-artifact@v5`, `actions/deploy-pages@v5`.
- **Permissions:** `contents: read`, `pages: write`, `id-token: write`.
- Builds with `pip install zensical` + `zensical build --clean`; deploys the `site/` dir via the official Pages actions.
- **No caching** is recommended right now ("caching functionality will undergo revisions").
- **Required GitHub Pages setting:** repo must be configured to publish **using GitHub Actions** (Settings → Pages → Build and deployment → Source: GitHub Actions). Docs land at `<username>.github.io/<repository>`.

## 6. Assets: logo / favicon

Source: https://zensical.org/docs/setup/logo-and-icons/

- Logo: `[project.theme] logo = "images/logo.png"` (any image incl. `.png`/`.svg`, path relative to `docs/`), or a bundled icon via `[project.theme.icon] logo = "lucide/smile"`. Logo links to `site_url` by default; override with `[project.extra] homepage = "..."`.
- Favicon: `[project.theme] favicon = "images/favicon.png"` (must live in `docs/`).
- Custom icon sets: `custom_dir = "overrides"` + `.icons/` folder, wired through `pymdownx.emoji` options (`custom_icons = ["overrides/.icons"]` with zensical emoji index/generator). 10,000+ bundled icons (Material, FontAwesome, Octicons, Lucide, Simple Icons, etc.).
- Customizable site icons include `logo`, `menu`, `search`, `repo`, `edit`, `view`, `previous`, `next`, `top`, admonition icons, tags.

## 7. Not-yet-supported summary (flagged)

| Capability | Status |
| --- | --- |
| `exclude_docs` / `draft_docs` (MkDocs settings) | ❌ Not supported; backlog #65, low priority; `exclude` plugin replacement "up next" on roadmap |
| Native versioned docs | ❌ Not shipped; transitional mike fork only; native versioning on roadmap, undated |
| `not_in_nav`, `hooks`, `remote_branch`, `remote_name` (MkDocs settings) | ❌ Not supported |
| `mkdocs gh-deploy` | ❌ Not provided; use the Pages workflow / publish guide |
| `docs_dir = "."` | ❌ Temporarily disallowed (backlog #101) |
| CI caching for builds | ⚠️ Officially discouraged for now ("caching will undergo revisions") |
| `[!IMPORTANT]` / `[!CAUTION]` styling | ⚠️ Callout types parse, but no dedicated built-in admonition styling (default styling unless custom CSS / type rewrite) |

## Open questions

1. Whether `important`/`caution` will become first-class admonition types, or whether a type-aliasing option will be added to `pymdownx.quotes`. Neither is documented; answering this would need a backlog search or asking the team.
2. The roadmap's "MkDocs `exclude` plugin replacement" will presumably take a different syntax than `exclude_docs`; its exact shape is unknown until shipped.
3. fltr-specific: whether to theme with `classic` (Material look) or `modern`. That is a design decision, not a research question.
