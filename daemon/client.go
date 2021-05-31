package daemon

import (
	"bytes"
	"context"
	"encoding/json"
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

var c http.Client = client()

func Connect() error {
	_, err := c.Get("http://unix/connect")

	return err
}

func Disconnect() error {
	_, err := c.Get("http://unix/disconnect")

	return err
}

func ExitNode(ip string) error {
	values := map[string]string{"exitNodeIp": ip}
	jsonData, err := json.Marshal(values)
	if err != nil {
		return err
	}

	_, err = c.Post("http://unix/exit-node", "application/json; charset=UTF-8", bytes.NewBuffer(jsonData))

	return err
}
