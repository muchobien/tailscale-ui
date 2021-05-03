package ui

import (
	"github.com/gotk3/gotk3/gtk"
)

func connect() *gtk.MenuItem {
	connect, err := gtk.MenuItemNewWithLabel("Connect")
	onError(err)

	connect.SetSensitive(false)

	return connect
}
