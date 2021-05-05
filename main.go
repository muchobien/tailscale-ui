package main

import (
	"flag"

	"github.com/muchobien/tailscale-ui/daemon"
	"github.com/muchobien/tailscale-ui/ui"
)

func main() {
	flag.Parse()
	if flag.Arg(0) == "daemon" {
		daemon.Listen()
	} else {
		ui.UI()
	}
}
