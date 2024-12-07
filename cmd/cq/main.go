package main

import (
	"io"
	"log"
	"os"

	"github.com/shizhMSFT/cq/internal/version"
	"github.com/shizhMSFT/cq/pkg/cq"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func main() {
	app := &cli.App{
		Name:    "cq",
		Usage:   "Command-line CBOR processor",
		Version: version.GetVersion(),
		Action: func(c *cli.Context) error {
			// Determine the input source
			var source io.Reader
			if args := c.Args(); args.Len() == 0 {
				if term.IsTerminal(int(os.Stdin.Fd())) {
					cli.ShowAppHelpAndExit(c, 2)
				}
				source = os.Stdin
			} else {
				path := args.First()
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				source = file
			}

			// Print the CBOR data
			return cq.Print(source)
		},
		HideHelpCommand: true,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
