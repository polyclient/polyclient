package main

import (
	"fmt"
	"log"

	"github.com/polyclient/polyclient/pkg/pluginsdk"
	"github.com/polyclient/polyclient/pkg/utils"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	socketPath := utils.GetSocketPath("db_sqlite")
	log.Printf("Plugin will create socket at: %s", socketPath)

	p := pluginsdk.NewPlugin("db_sqlite", "0.0.1")

	p.RegisterHandler("query", func(payload []byte, metadata map[string]string) ([]byte, error) {
		fmt.Print(string(payload))

		return []byte("take some users"), nil
	})

	if err := pluginsdk.Serve(p); err != nil {
		log.Fatalf("Failed to serve plugin: %v", err)
	}
}
