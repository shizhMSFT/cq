package main

import (
	"log"
	"os"

	"github.com/shizhMSFT/cq/internal/version"
	"github.com/shizhMSFT/cq/pkg/cq"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "cq",
		Usage:   "Command-line CBOR processor",
		Version: version.GetVersion(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				return cq.Print(os.Stdin)
			}
			path := c.Args().First()
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			return cq.Print(file)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
