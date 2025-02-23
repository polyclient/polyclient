// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package plugin

import (
	"context"
	"fmt"
	"log"

	"github.com/polyclient/polyclient/runtime/plugin"
	"github.com/urfave/cli/v3"
)

func newLoadCommand() *cli.Command {
	return &cli.Command{
		Name:     "load",
		Usage:    "Load a plugin",
		Category: "Plugins",
		Action: func(context.Context, *cli.Command) error {
			lookupPaths := []string{
				"./plugins",
			}
			pr := plugin.NewPluginRegistry(lookupPaths)

			if err := pr.LoadPlugins(); err != nil {
				return err
			}

			plugin, err := pr.GetWASMPlugin("sqlite")
			if err != nil {
				return err
			}

			_, result, err := plugin.Call("greet", []byte("Juan"))
			if err != nil {
				return fmt.Errorf("failed to call function: %w", err)
			}

			log.Println(string(result))
			return nil
		},
	}
}
