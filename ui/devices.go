package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/tailcfg"
	"tailscale.com/util/dnsname"
)

func devices() *gtk.MenuItem {
	ctx := context.Background()

	clipboard, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	onError(err)

	devices, err := gtk.MenuItemNewWithLabel("Network devices")
	onError(err)

	subMenuDevices, err := gtk.MenuNew()
	onError(err)

	go func() {
		for {
			status, err := tailscale.Status(ctx)
			onError(err)

			grouped := make(map[tailcfg.UserID][]*ipnstate.PeerStatus)

			for user := range status.User {
				grouped[user] = []*ipnstate.PeerStatus{}
			}

			for _, peer := range status.Peer {
				grouped[peer.UserID] = append(grouped[peer.UserID], peer)
			}

			for user, items := range grouped {
				userName, err := gtk.MenuItemNewWithLabel(status.User[user].DisplayName)
				onError(err)

				userDevices, err := gtk.MenuNew()
				onError(err)

				userName.SetSubmenu(userDevices)

				for _, peer := range items {
					hostName := dnsname.SanitizeHostname(peer.HostName)
					dnsName := dnsname.TrimSuffix(peer.DNSName, status.MagicDNSSuffix)
					name := fmt.Sprintf("%s (%s)", dnsName, hostName)

					if strings.EqualFold(dnsName, hostName) || peer.UserID != status.Self.UserID {
						name = dnsName
					}

					item, err := gtk.MenuItemNewWithLabel(name)
					onError(err)

					ip := peer.TailscaleIPs[0].String()
					item.Connect("activate", func() {
						clipboard.SetText(ip)
					})
					userDevices.Add(item)
				}
				subMenuDevices.Add(userName)
			}
			subMenuDevices.ShowAll()

			devices.SetSubmenu(subMenuDevices)

			<-time.After(time.Second * 5)
		}
	}()

	return devices
}
