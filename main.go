package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/muchobien/tailscale-ui/ui"
	"gopkg.in/retry.v1"
	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
)

func main() {
	gtk.Init(nil)
	ui.SetOnError(func(err error) {
		if err != nil {
			log.Fatalln(err)
		}
	})

	ctx := context.Background()
	strategy := retry.LimitTime(30*time.Second,
		retry.Exponential{
			Initial: 10 * time.Millisecond,
			Factor:  1.5,
		},
	)

	var status *ipnstate.Status
	var err error

	for a := retry.Start(strategy, nil); a.Next(); {
		status, err = tailscale.Status(ctx)
		if err == nil {
			break
		}
	}

	if err == nil {
		ui.UI(status)
		gtk.Main()
	}

	os.Exit(1)
}
