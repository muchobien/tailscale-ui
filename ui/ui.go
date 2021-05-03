package ui

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/dawidd6/go-appindicator"
	"github.com/gotk3/gotk3/gtk"
	"github.com/muchobien/tailscale-ui/utils"
	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
)

func UI(st *ipnstate.Status) {
	menu, err := gtk.MenuNew()
	onError(err)

	connect := connect()
	status := status()
	disconnect := disconnect()
	separators := buildMenuItemSeparators(5)

	connect.Connect("activate", func() {
		exec.Command("sudo", "tailscale", "up").Start()
		disconnect.SetSensitive(true)
	})

	disconnect.Connect("activate", func() {
		exec.Command("sudo", "tailscale", "down").Start()
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
