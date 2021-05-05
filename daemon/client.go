package daemon

import (
	"context"
	"net"
	"net/http"
)

func client() http.Client {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", SocketPath)
			},
		},
	}

	return httpc
}

func Connect() error {
	c := client()
	_, err := c.Get("http://unix/connect")

	return err
}

func Disconnect() error {
	c := client()
	_, err := c.Get("http://unix/disconnect")

	return err
}
