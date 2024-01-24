package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

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

	entry := c.Commands.Find(flag.Arg(0))
	if entry == nil {
		flag.Usage()
		os.Exit(0)
	}

	if err := entry.Run(flag.Args()[1:]); err != nil {
		log.Fatal(err)
	}
}
