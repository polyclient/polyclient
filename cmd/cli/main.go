package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "polyclient",
		Usage: "Run the CLI client",
		Action: func(c *cli.Context) error {
			fmt.Println("Running the CLI client")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
