# PolyClient Documentation

This is the official documentation for PolyClient. For users, it contains guides for installing, configuring, and using PolyClient. For developers, it covers the core components, database management, advanced features, plugin system, integrations, and contributing guidelines.

## Structure

PolyClient’s documentation is organized into sections for easy navigation:

- **Getting Started** – Installation, configuration, and first steps.
- **Core Components** – Overview of the CLI, GUI, API, and architecture.
- **Database Management** – Handling connections, queries, and schema management.
- **Advanced Features** – Security, automation, performance monitoring, and backup strategies.
- **Plugin System** – Authoring, installing, publishing, and managing plugins.
- **Integration** – API usage, third-party service integration, and custom extensions.
- **Developer Reference** – Building from source, API documentation, and internal architecture.
- **Benchmarks** – Performance comparisons and testing methodologies.
- **Contributing** – Guidelines for submitting code, documentation, and testing contributions.
- **Troubleshooting** – Common issues, debugging, and community support.

## Running the documentation locally

To preview the documentation on your machine, run:

```sh
task dev:docs
```

This starts a local development server with real-time updates.

> [!NOTE]
> You need to have [mdBook](https://rust-lang.github.io/mdBook/) and [mdBook-mermaid](https://crates.io/crates/mdbook-mermaid/) installed.

## Building the documentation

To generate a static build, use:

```sh
task build:docs
```

This compiles the documentation into a deployable format with all pages and assets correctly structured.

## Accessing the documentation online

The latest version is always available at:

<https://polyclient.pages.dev/docs>

## Contributing to the documentation

- **Quick fixes**: Small changes, such as fixing typos or broken links, can be submitted directly via pull requests.
- **Major changes**: For larger updates, open an issue on GitHub to discuss your proposal before implementation.

There are plans to support additional languages. There's a [bash script](../scripts/make-docs.sh) that can be used to generate the entire documentation directory structure and placeholder files for a specific language, but this should be discussed beforehand, as it adds maintenance overhead.

To use the script, run:

```sh
./scripts/make-docs.sh <lang>

# For example:
./scripts/make-docs.sh es # Spanish
./scripts/make-docs.sh it # Italian
./scripts/make-docs.sh ja # Japanese
```

## License

[Apache-2.0](LICENSE.txt)
