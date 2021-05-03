package ui

import (
	"os/exec"

	"github.com/gotk3/gotk3/gtk"
)

func console() *gtk.MenuItem {
	console, err := gtk.MenuItemNewWithLabel("Admin console...")
	onError(err)

	console.Connect("activate", func() {
		exec.Command("xdg-open", "https://login.tailscale.com/admin/machines").Start()
	})

	return console
}
