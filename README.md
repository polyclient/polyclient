> [!IMPORTANT]
> This project is in early development. Features and documentation are expected to change.

[![release](https://github.com/polyclient/polyclient/actions/workflows/release.yaml/badge.svg)](https://github.com/polyclient/polyclient/actions/workflows/release.yaml)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_shield&issueType=license)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_shield&issueType=security)

---

[**Website**](https://polyclient.pages.dev) · [**Documentation**](https://polyclient.pages.dev/docs) · [**Blog**](https://polyclient.pages.dev/blog) · [**Plugins**](https://polyclient.pages.dev/plugins)

---

# PolyClient

PolyClient is a fast, extensible, and cross-platform database management platform designed to seamlessly manage SQL and NoSQL databases through a unified interface. It offers a powerful set of tools for database administration, integration, and customization.

## Features

### Core

- 💸 NO SaaS model—everything's free and unrestricted.
- 🚫 No telemetry or tracking.
- 💻 Cross-platform (Linux, macOS, Windows).
- 👨🏻‍💻 Powerful CLI to manage your databases without leaving your terminal.
- 🖼️ Fast, intuitive, and extensible GUI with clean, clutter-free design.

### Extensibility and integration

- 🌐 Headless API for building custom database tools:
  - Use PolyClient as a backend for tailored clients.
  - Expose plugin-driven functionality via a structured API.
  - Integrate with web apps, dashboards, and automation systems.
- 🔌 Secure WebAssembly (Wasm) plugin system to extend functionality:
  - Add support for new databases.
  - Create themes and customize UI components.
  - Automate workflows.
  - Integrate LLMs and AI tools.
  - Connect with cloud services and external APIs.

### Data management

- 📤📥 Versatile data import and export capabilities:
  - Export to CSV, TSV, JSON, XML, HTML, YAML, TOML, Markdown.
  - Import from CSV, TSV, JSON, XML, SQL dumps.
  - Smart data mapping for seamless imports across different schemas.
- 💾 Robust backup and restore functionality:
  - Scheduled automated backups.
  - Point-in-time recovery options.
  - Cross-database backup and restore support.

### Security and access control

- 🛡️ Advanced security with dynamic data masking and row-level security.
- 🔒 Role-based access control, user management, and audit logging.

### Query and performance

- 👁️ Visual query builder—no database-specific knowledge needed.
- ⚡️ Fast queries with parallel execution and caching.
- 🗂️ Query history and favorites for quick access to frequently used queries.
- 🧠 AI-powered query optimization.
- 💬 Natural language query interface—talk to your database like a human.

### Monitoring and optimization

- 🏥 Proactive database health monitoring with smart recommendations.
- 📡 Monitor and optimize database performance with observability tools.

### Database comparison and synchronization

- 🔄 Built-in database comparison and synchronization tools.
- 🔍 Schema diff and sync capabilities across multiple databases.

### Miscellaneous

- 📊 Generate unlimited ERDs that you can share with your friends.
- ⏳ Automate tasks with scheduled jobs.
- 🌍 Multilingual support (i18n).
- 📖 Comprehensive documentation, guides, and API references.

## Database support

**Legend:**

- **Support Status:**
  - ✅ Supported: Fully functional and reliable.
  - ⚠️ Limited: Partial functionality, potential issues.
  - 🚧 Developing: Actively being developed, may have changes.
  - 🗓️ Upcoming: Scheduled for future release.
- **First Available Version:** Version `x.y.z` where support was introduced.

| Database        | Type        | Support Status | Version | Reference                                                                                           |
| :-------------- | :---------- | :------------- | :------ | :-------------------------------------------------------------------------------------------------- |
| PostgreSQL      | Relational  | 🗓️ Upcoming    | TBD     | [Source](plugins/postgresql) · [Documentation](https://polyclient.pages.dev/databases/postgresql)   |
| MySQL           | Relational  | 🗓️ Upcoming    | TBD     | [Source](plugins/mysql) · [Documentation](https://polyclient.pages.dev/databases/mysql)             |
| SQLite          | Relational  | 🚧 Developing  | TBD     | [Source](plugins/sqlite) · [Documentation](https://polyclient.pages.dev/databases/sqlite)           |
| SQL Server      | Relational  | 🗓️ Upcoming    | TBD     | [Source](plugins/mssql) · [Documentation](https://polyclient.pages.dev/databases/mssql)             |
| Oracle Database | Relational  | 🗓️ Upcoming    | TBD     | [Source](plugins/oracle) · [Documentation](https://polyclient.pages.dev/databases/oracle)           |
| MariaDB         | Relational  | 🗓️ Upcoming    | TBD     | [Source](plugins/mariadb) · [Documentation](https://polyclient.pages.dev/databases/mariadb)         |
| Firebird        | Relational  | 🗓️ Upcoming    | TBD     | [Source](plugins/firebird) · [Documentation](https://polyclient.pages.dev/databases/firebird)       |
| LibSQL          | Relational  | 🗓️ Upcoming    | TBD     | [Source](plugins/libsql) · [Documentation](https://polyclient.pages.dev/databases/libsql)           |
|                 |             |                |         |                                                                                                     |
| MongoDB         | NoSQL       | 🗓️ Upcoming    | TBD     | [Source](plugins/mongodb) · [Documentation](https://polyclient.pages.dev/databases/mongodb)         |
| Redis           | NoSQL       | 🗓️ Upcoming    | TBD     | [Source](plugins/redis) · [Documentation](https://polyclient.pages.dev/databases/redis)             |
| Cassandra       | NoSQL       | 🗓️ Upcoming    | TBD     | [Source](plugins/cassandra) · [Documentation](https://polyclient.pages.dev/databases/cassandra)     |
|                 |             |                |         |                                                                                                     |
| CockroachDB     | Distributed | 🗓️ Upcoming    | TBD     | [Source](plugins/cockroachdb) · [Documentation](https://polyclient.pages.dev/databases/cockroachdb) |
|                 |             |                |         |                                                                                                     |
| DuckDB          | Analytical  | 🗓️ Upcoming    | TBD     | [Source](plugins/duckdb) · [Documentation](https://polyclient.pages.dev/databases/duckdb)           |
| ClickHouse      | Analytical  | 🗓️ Upcoming    | TBD     | [Source](plugins/clickhouse) · [Documentation](https://polyclient.pages.dev/databases/clickhouse)   |

## Installation

TBD

## Benchmarks

TBD

## Contributing

<details>
    <summary>Local development</summary>

- Clone this repository.
- Install:
  - [Go >= 1.24](https://go.dev)
  - [Task](https://taskfile.dev)
  - [Bun](https://bun.sh)
- Run `task install` to install the project dependencies.
- Run `task --list` to list all the available tasks.

</details>

<details>
    <summary>Remote development</summary>

TBD

</details>

<br>

Refer to the [contributor's guide](CONTRIBUTING.md) for more in-depth information.

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_large&issueType=license)

This project is distributed under multiple licenses following the [SPDX](https://spdx.dev/) standard. The primary codebase is licensed under [GPL-3.0-or-later](LICENSES/GPL-3.0-or-later.txt) with a [Plugin Exception](LICENSES/LicenseRef-PolyClient-Plugin-Exception.txt), while plugins are licensed under [Apache-2.0](LICENSES/Apache-2.0.txt). The website and documentation are licensed under [CC-BY-4.0](LICENSES/CC-BY-4.0.txt). This project follows the [REUSE Specification](https://reuse.software/) to ensure clear and consistent licensing.
