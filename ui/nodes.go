package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"inet.af/netaddr"
	"tailscale.com/client/tailscale"
	"tailscale.com/tailcfg"
	"tailscale.com/util/dnsname"
)

var (
	ipv4default = netaddr.MustParseIPPrefix("0.0.0.0/0")
	ipv6default = netaddr.MustParseIPPrefix("::/0")
)

func isExitNode(node *tailcfg.Node) bool {
	for _, ip := range node.AllowedIPs {
		ipp := netaddr.MustParseIPPrefix(ip.String())
		if ipp == ipv4default || ipp == ipv6default {
			return true
		}
	}

	return false
}

func nodes() *gtk.MenuItem {
	ctx := context.Background()

	nodes, err := gtk.MenuItemNewWithLabel("Exit node")
	onError(err)

	subMenuNodes, err := gtk.MenuNew()
	onError(err)

	go func() {
		for {
			separators := buildMenuItemSeparators(1)
			status, err := tailscale.Status(ctx)
			onError(err)

			subMenuNodes.GetChildren().Foreach(func(item interface{}) {
				subMenuNodes.Remove(item.(*gtk.Widget))
			})

			var group *gtk.RadioMenuItem

			none, err := gtk.RadioMenuItemNewWithLabelFromWidget(group, "None")
			onError(err)

			group = none

			none.SetActive(true)

			subMenuNodes.Add(none)
			subMenuNodes.Add(separators[0])

			group.Connect("toggled", func() {

			})

			for _, peer := range status.Peer {
				node, err := tailscale.WhoIs(ctx, fmt.Sprintf("%s:0", peer.TailscaleIPs[0]))
				onError(err)

				if isExitNode(node.Node) {
					hostName := dnsname.SanitizeHostname(peer.HostName)
					nodeMenu, err := gtk.RadioMenuItemNewWithLabelFromWidget(group, hostName)
					onError(err)

					group = nodeMenu

					group.SetActive(peer.ExitNode)

					group.Connect("toggled", func() {

					})

					subMenuNodes.Add(nodeMenu)

				}

			}

			nodes.SetSubmenu(subMenuNodes)
			subMenuNodes.ShowAll()
			<-time.After(time.Minute)
		}
	}()

	return nodes
}
