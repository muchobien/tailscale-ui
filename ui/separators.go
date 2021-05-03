package ui

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func buildMenuItemSeparators(number int) []*gtk.SeparatorMenuItem {

	separators := make([]*gtk.SeparatorMenuItem, number)

	for i := 0; i < number; i++ {
		sep, err := gtk.SeparatorMenuItemNew()
		if err != nil {
			log.Fatal(err)
		}

		separators[i] = sep
	}

	return separators
}
