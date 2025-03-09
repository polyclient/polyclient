# SQLite Plugin

This plugin provides SQLite database access for PolyClient.

> [!NOTE]
> This is a built-in plugin, which means it is loaded automatically when PolyClient starts. If you need to unload it for some reason, use the following CLI command:
>
> ```sh
> polyclient plugin unload --id "sqlite"
> ```
>
> To load the plugin again, use the appropriate CLI command for your operating system:
>
> **Unix:**
>
> ```sh
> polyclient plugin load --path "$POLYCLIENT_PLUGINS_DIR/sqlite"
> ```
>
> **Windows:**
>
> ```powershell
> polyclient plugin load --path "%POLYCLIENT_PLUGINS_DIR%\sqlite"
> ```

## Documentation

<https://polyclient.pages.dev/plugins/sqlite>

## License

[Apache-2.0](LICENSE.txt)
