> [!IMPORTANT]
> This project is in early development. Features and documentation are expected to change.

[![release](https://github.com/polyclient/polyclient/actions/workflows/release.yaml/badge.svg)](https://github.com/polyclient/polyclient/actions/workflows/release.yaml)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_shield&issueType=license)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=shield&issueType=security)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_shield&issueType=security)

# PolyClient

PolyClient is an open-source, cross-platform client for managing SQL and NoSQL databases. It provides a single interface to handle various database systems efficiently. Its core features include:

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

Refer to the [contributor's guide](CONTRIBUTING.md) for more in-depth information.

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpolyclient%2Fpolyclient?ref=badge_large&issueType=license)

This project is licensed under the [GNU General Public License v3.0](LICENSE-GPL), with a [runtime library linkage exception](COPYING.RUNTIME). Built-in plugin modules are licensed under the [Apache License v2.0](LICENSE-APACHE), The runtime library linkage exception allows you to build and distribute your own linked modules under any terms you choose, including proprietary licenses.
