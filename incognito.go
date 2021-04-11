package main

import (
	"flag"
	"github.com/MarshallRawson/incognito/cli"
	"github.com/MarshallRawson/incognito/front_end"
)

func main() {

	command_line := flag.Bool("cli", false, "use the command line interface")
	landing_page := flag.String("landing_page", "", "landing_page_for_gui")

	flag.Parse()

	if (*command_line) == true {
		cli.Run()
	} else {
		front_end.Run(*landing_page)
	}
}
