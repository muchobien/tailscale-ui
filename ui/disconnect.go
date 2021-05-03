package ui

import "github.com/gotk3/gotk3/gtk"

func disconnect() *gtk.MenuItem {
	disconnect, err := gtk.MenuItemNewWithLabel("Disconnect")
	onError(err)

	return disconnect
}
