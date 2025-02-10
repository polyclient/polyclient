package main

import (
	"log"

	"github.com/polyclient/polyclient/pkg/pluginsdk"
)

func main() {
	plugin, err := pluginsdk.NewPlugin()
	if err != nil {
		log.Fatal("failed to initialize plugin:", err)
	}

	plugin.RegisterHandler("query", func(payload []byte) ([]byte, error) {
		return []byte("[{\"id\": 1, \"name\": \"John\"}, {\"id\": 2, \"name\": \"Jane\"}]"), nil
	})

	pluginsdk.Serve(plugin)
}
