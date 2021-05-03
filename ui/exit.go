package ui

import "github.com/gotk3/gotk3/gtk"

func exit() *gtk.MenuItem {
	exit, err := gtk.MenuItemNewWithLabel("Exit")
	onError(err)
	exit.Connect("activate", gtk.MainQuit)

	return exit
}
