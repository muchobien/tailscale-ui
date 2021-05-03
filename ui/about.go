package ui

import "github.com/gotk3/gotk3/gtk"

func about() *gtk.MenuItem {
	about, err := gtk.MenuItemNewWithLabel("About Tailscale-ui")
	onError(err)

	return about
}
