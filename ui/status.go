package ui

import (
	"github.com/gotk3/gotk3/gtk"
)

func status() *gtk.MenuItem {
	status, err := gtk.MenuItemNewWithLabel("0 received | 0 sent | 0 link")
	onError(err)

	return status
}
