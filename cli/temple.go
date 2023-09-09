package main

import (
	"log"
	"os"

	"github.com/ronyv89/temple/generators/ruby"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "temple",
		Usage: "generate templates for your favorite languages and frameworks",
		Action: func(*cli.Context) error {
			ruby.NewRailsProject()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
