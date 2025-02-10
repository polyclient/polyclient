package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/polyclient/polyclient/internal/config"
	"github.com/polyclient/polyclient/internal/plugin"
)

func main() {
	pluginsDirs := []string{"plugins"}
	externalPluginsDirs := []string{}

	userConfigDir, err := os.UserConfigDir()
	if err == nil {
		configPlugins := path.Join(userConfigDir, config.PolyClientConfigDir, config.PolyClientPluginsDir)
		if _, err := os.Stat(configPlugins); err != nil && os.IsNotExist(err) {
			os.MkdirAll(configPlugins, 0755)
		}

		externalPluginsDirs = append(externalPluginsDirs, configPlugins)
	}

	pluginsDirs = append(pluginsDirs, externalPluginsDirs...)

	pm := plugin.NewPluginManager(plugin.NewPluginManagerOptions{
		PluginsDirs: pluginsDirs,
	})

	loadCount, err := pm.LoadPlugins()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load plugins: %w", err))
	}

	if loadCount == 0 {
		log.Default().Println("No plugins found")
	}

	plugins := pm.GetPlugins()
	fmt.Println(plugins)

	result, err := pm.Execute("polyclient-db-sqlite", "query", []byte("SELECT * FROM users"), map[string]string{})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to execute action: %w", err))
	}

	fmt.Println(string(result))

	pm.Shutdown()
}
