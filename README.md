> [!IMPORTANT]
> This project is currently in a very early stage of development. As such, features and documentation are subject to change.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_shield&issueType=license)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_shield&issueType=security)

# PolyClient

PolyClient is a cross-platform, open-source, and extensible client designed for managing and querying a variety of SQL and NoSQL databases. It provides a unified interface to simplify database interaction, regardless of the specific database system. PolyClient offers a range of features to support common database workflows:

- **Database management**: Manage multiple databases in a single interface.
- **Query execution**: Run SQL and NoSQL queries with syntax highlighting and code completion with support for **parallel query execution**.
- **Data visualization**: Visualize data in tables, charts, graphs, and maps.
- **Data Import/Export**: import/export data in various formats (CSV, JSON, XML, etc.) with zero restrictions.
- **ERD generation**: Auto-generate Entity-Relationship Diagrams (ERDs).
- **AI assistance**: AI-driven query generation, completion, and optimization.
- **Plugin extensibility**: Easily integrate any functionality you want to the core that isn't provided by PolyClient.
- **Cross-Platform support**: Available on Linux, macOS, and Windows.
- **Open source**: Free to use, modify, and distribute.

## Plugin Support

PolyClient's powerful plugin architecture decouples core functionality, enabling third parties to extend it (like adding support for a new database or AI provider). Learn how to author plugins [here](https://polyclient.pages.dev/en/plugins/).

**Legend**:

- 🟢 Supported – Fully functional
- 🟡 In Progress – Actively being developed
- 🟠 Partial – Some features may be limited
- 🔴 Planned – Not yet implemented
- **First supported version**: Version where support was introduced.

### Databases

| Database                                                 | Status     | First supported version | Plugin                                                                                                  |
| -------------------------------------------------------- | ---------- | ----------------------- | ------------------------------------------------------------------------------------------------------- |
| [MariaDB](https://mariadb.org/)                          | 🔴 Planned | TBD                     | [Source](plugins/db_mariadb) · [Documentation](https://polyclient.pages.dev/en/databases/mariadb)       |
| [MongoDB](https://www.mongodb.com/)                      | 🔴 Planned | TBD                     | [Source](plugins/db_mongodb) · [Documentation](https://polyclient.pages.dev/en/databases/mongodb)       |
| [MySQL](https://www.mysql.com/)                          | 🔴 Planned | TBD                     | [Source](plugins/db_mysql) · [Documentation](https://polyclient.pages.dev/en/databases/mysql)           |
| [PostgreSQL](https://postgresql.org/)                    | 🔴 Planned | TBD                     | [Source](plugins/db_postgresql) · [Documentation](https://polyclient.pages.dev/en/databases/postgresql) |
| [Redis](https://redis.io/)                               | 🔴 Planned | TBD                     | [Source](plugins/db_redis) · [Documentation](https://polyclient.pages.dev/en/databases/redis)           |
| [SQL Server](https://www.microsoft.com/en-us/sql-server) | 🔴 Planned | TBD                     | [Source](plugins/db_mssql) · [Documentation](https://polyclient.pages.dev/en/databases/mssql)           |
| [SQLite](https://www.sqlite.org/)                        | 🔴 Planned | TBD                     | [Source](plugins/db_sqlite) · [Documentation](https://polyclient.pages.dev/en/databases/sqlite)         |

### AI Providers

| AI provider                            | Status     | First supported version | Plugin                                                                                         |
| -------------------------------------- | ---------- | ----------------------- | ---------------------------------------------------------------------------------------------- |
| [Anthropic](https://www.anthropic.com) | 🔴 Planned | TBD                     | [Source](plugins/ai_anthropic) · [Documentation](https://polyclient.pages.dev/en/ai/anthropic) |
| [DeepSeek](https://deepseek.ai/)       | 🔴 Planned | TBD                     | [Source](plugins/ai_deepseek) • [Documentation](https://polyclient.pages.dev/en/ai/deepseek)   |
| [Google AI](https://ai.google.com)     | 🔴 Planned | TBD                     | [Source](plugins/ai_google) • [Documentation](https://polyclient.pages.dev/en/ai/google)       |
| [LM Studio](https://lmstudio.ai/)      | 🔴 Planned | TBD                     | [Source](plugins/ai_lmstudio) • [Documentation](https://polyclient.pages.dev/en/ai/lmstudio)   |
| [Mistral](https://mistral.ai/)         | 🔴 Planned | TBD                     | [Source](plugins/ai_mistral) • [Documentation](https://polyclient.pages.dev/en/ai/mistral)     |
| [Ollama](https://ollama.ai/)           | 🔴 Planned | TBD                     | [Source](plugins/ai_ollama) • [Documentation](https://polyclient.pages.dev/en/ai/ollama)       |
| [OpenAI](https://openai.com/)          | 🔴 Planned | TBD                     | [Source](plugins/ai_openai) • [Documentation](https://polyclient.pages.dev/en/ai/openai)       |

## Third-party plugins

> [!NOTE]
> Submit a pull request if you created a plugin and want it to be listed here.

## Installation

To be added

## Contributing

<details>
    <summary>Local development</summary>

- Clone this repository.
- Install:
  - [Go](https://go.dev/)
  - [Task](https://taskfile.dev/)
  - [Bun](https://bun.sh/)
- Run `task install` to install the project dependencies.
- Run `task --list` to list the available commands.

</details>

Refer to the [contributor's guide](CONTRIBUTING.md) for more in-depth information.

## License

[![image](https://www.gnu.org/graphics/gplv3-with-text-136x68.png)](https://www.gnu.org/licenses/gpl-3.0.html)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_large&issueType=license)

PolyClient's internals are licensed under the [GNU General Public License v3.0 (GPLv3)](LICENSE-GPL). Reusable packages and plugins are licensed under the [MIT license](LICENSE-MIT).

### References

- [GPLv3 license in plain English](https://www.tldrlegal.com/license/gnu-general-public-license-v3-gpl-3)
- [MIT license in plain English](https://www.tldrlegal.com/license/mit-license)
