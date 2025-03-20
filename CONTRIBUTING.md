# Contributor's guide

This page provides information on reporting issues and guidelines for open-source contributors. Before contributing, review the PolyClient [documentation](https://polyclient.pages.dev) for project context.

## License

By contributing, you agree to license your contributions under the project's applicable licenses. See the [License section](README.md#license) in the README for details.

## How to contribute

We welcome community contributions. Here's how you can help:

- **Feature requests:** Share ideas via [feedback discussions](https://github.com/polyclient/polyclient/discussions). Accepted ideas become issues.
- **Documentation:** Fix typos or links with pull requests. For major changes, submit an issue first.
- **Bug discovery:** Report new bugs ([check existing issues](https://github.com/polyclient/polyclient/issues)). Fix them or leave them for others.
- **Plugin authoring:** Extend functionality by creating plugins. [Plugin docs](https://polyclient.pages.dev/plugins).

## Contribution Workflow

1. Someone open an issue to report bugs or propose changes.
2. Maintainers label issues (e.g., `duplicate`, `blocked`, `needs owner`, `needs info`, `wontfix`, `breaking change`) to prioritize work.
3. Volunteers express interest in solving an issue.
4. Maintainers assign issues to owners.
5. Owners submit a pull request with changes.
6. Maintainers review and merge approved pull requests.

## Good first issues

If you're new to the codebase, consider starting with issues labeled as [good first issue](https://github.com/polyclient/polyclient/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22+-label%3A%22blocked+by+upstream%22). Before you start, check for existing pull requests or assignments. Notify maintainers with a comment on the issue to avoid duplicate work.

## Coding Style

We follow [Effective Go](https://go.dev/doc/effective_go) principles, ensuring consistency with:

- **Go:** `golangci-lint` ([config](.golangci.toml)) for linting and `gofmt` for standardized formatting.
- **TypeScript:** `Biome` ([config](gui/biome.json)) for linting and formatting, maintaining TypeScript code quality.
- **Cross-editor consistency:** `EditorConfig` ([config](.editorconfig)) to enforce uniform styles across different editors and IDEs.

## Git workflow

We use trunk-based development (<https://trunkbaseddevelopment.com>).

1. Developers work on short-lived branches off main, integrating changes frequently.
2. Create a pull request to merge into main. At least one maintainer must approve.
3. Approved pull requests are merged into main, triggering automated workflows.

```mermaid
gitGraph
    commit
    commit
    branch feature/branch1
    checkout feature/branch1
    commit
    commit
    checkout main
    merge feature/branch1
    commit
    commit
```
