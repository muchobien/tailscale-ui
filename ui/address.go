package ui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"tailscale.com/ipn/ipnstate"
)

func address(status *ipnstate.Status) *gtk.MenuItem {
	clipboard, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	onError(err)

	address, err := gtk.MenuItemNewWithLabel("My address: " + status.Self.TailscaleIPs[0].String())
	onError(err)

	address.Connect("activate", func() {
		clipboard.SetText(status.Self.TailscaleIPs[0].String())
	})

	return address
}
