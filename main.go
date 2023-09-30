package main

import (
	"flag"

	"github.com/blubooks/blubooks-cli/cmd"
)

func main() {

	server := flag.Bool("server", false, "Start Server")

	flag.Parse()

	if *server {
		cmd.Server()
		return
	}
	cmd.Server()
}
