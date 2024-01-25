package main

import (
	"flag"
	"log"
	"os"

	"github.com/posener/complete/v2"
)

func main() {
	// config
	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// cli completion
	cmd := &complete.Command{
		Sub: c.Commands.Completions(),
	}

	cmd.Complete("console")

	// flags
	flag.Usage = CliUsage(
		c.Title,
		c.Usage,
		c.Commands,
	)
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// run
	entry := c.Commands.Find(flag.Arg(0))
	if entry == nil {
		flag.Usage()
		os.Exit(0)
	}

	if err := entry.Run(flag.Args()[1:]); err != nil {
		log.Fatal(err)
	}
}
