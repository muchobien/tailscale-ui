package ui

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dawidd6/go-appindicator"
	"github.com/gotk3/gotk3/gtk"
	"github.com/muchobien/tailscale-ui/daemon"
	"github.com/muchobien/tailscale-ui/utils"
	"gopkg.in/retry.v1"
	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
)

func ui(st *ipnstate.Status) {
	menu, err := gtk.MenuNew()
	onError(err)

	connect := connect()
	status := status()
	disconnect := disconnect()
	separators := buildMenuItemSeparators(5)

	connect.Connect("activate", func() {
		err := daemon.Connect()
		onError(err)
		disconnect.SetSensitive(true)
		connect.SetSensitive(false)
	})

	disconnect.Connect("activate", func() {
		err := daemon.Disconnect()
		onError(err)
		connect.SetSensitive(true)
		disconnect.SetSensitive(false)
	})

	indicator := appindicator.New("dev.muchobien.tailscale-ui", "tailscale-ui", appindicator.CategoryApplicationStatus)
	indicator.SetIconThemePath("/usr/share/icons/hicolor/scalable/apps")
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	menu.Add(about())
	menu.Add(separators[0])
	menu.Add(address(st))
	menu.Add(status)
	menu.Add(separators[1])
	menu.Add(connect)
	menu.Add(disconnect)
	menu.Add(separators[2])
	menu.Add(console())
	menu.Add(separators[3])
	menu.Add(devices())
	menu.Add(separators[4])
	menu.Add(exit())
	menu.ShowAll()

	go func() {
		ctx := context.Background()

		for {
			<-time.After(time.Second)
			st, err := tailscale.Status(ctx)
			onError(err)
			txBytes, rxBytes := int64(0), int64(0)

			for _, peerStatus := range st.Peer {
				if peerStatus.ShareeNode {
					continue
				}

				txBytes += peerStatus.TxBytes
				rxBytes += peerStatus.RxBytes
			}
			status.SetLabel(fmt.Sprintf("%s received | %s sent | %d link", utils.FmtByte(rxBytes), utils.FmtByte(txBytes), len(st.Peers())))
		}
	}()
}

func UI() {
	gtk.Init(nil)
	SetOnError(func(err error) {
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
		ui(status)
		gtk.Main()
	}

	os.Exit(1)
}
