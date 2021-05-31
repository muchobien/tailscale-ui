package daemon

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"inet.af/netaddr"
	"tailscale.com/client/tailscale"
	"tailscale.com/ipn"
)

const SocketPath = "/var/run/tailscale-ui.sock"

type exitNodeBody struct {
	ExitNodeIp string
}

func connect(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	st, err := tailscale.Status(ctx)
	if err != nil {
		http.Error(w, "Error fetching current status", http.StatusInternalServerError)
		return
	}

	if st.BackendState == "Running" {
		http.Error(w, "Tailscale was already Running", http.StatusConflict)
		return
	}

	_, err = tailscale.EditPrefs(ctx, &ipn.MaskedPrefs{
		Prefs: ipn.Prefs{
			WantRunning: true,
		},
		WantRunningSet: true,
	})

	if err != nil {
		http.Error(w, "Failed to Connect", http.StatusInternalServerError)
		return
	}
}

func disconnect(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	st, err := tailscale.Status(ctx)
	if err != nil {
		http.Error(w, "Error fetching current status", http.StatusInternalServerError)
		return
	}

	if st.BackendState == "Stopped" {
		http.Error(w, "Tailscale was already stopped", http.StatusConflict)
	}

	_, err = tailscale.EditPrefs(ctx, &ipn.MaskedPrefs{
		Prefs: ipn.Prefs{
			WantRunning: false,
		},
		WantRunningSet: true,
	})

	if err != nil {
		http.Error(w, "Failed to Disconnect", http.StatusInternalServerError)
		return
	}
}

func exitNode(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}

	dec := json.NewDecoder(req.Body)

	var payload exitNodeBody

	if err := dec.Decode(&payload); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	var exitNodeIP netaddr.IP
	if payload.ExitNodeIp != "" {
		var err error
		exitNodeIP, err = netaddr.ParseIP(payload.ExitNodeIp)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	}

	ctx := context.Background()
	_, err := tailscale.EditPrefs(ctx, &ipn.MaskedPrefs{
		Prefs: ipn.Prefs{
			ExitNodeIP: exitNodeIP,
		},
		WantRunningSet: true,
	})
	if err != nil {
		http.Error(w, "Failed to set exitnode", http.StatusInternalServerError)
	}
}

func Listen() {
	os.Remove(SocketPath)

	server := http.Server{}

	http.HandleFunc("/connect", connect)
	http.HandleFunc("/disconnect", disconnect)
	http.HandleFunc("/exit-node", exitNode)

	unixListener, err := net.Listen("unix", SocketPath)
	if err != nil {
		panic(err)
	}
	if err := os.Chmod(SocketPath, 0777); err != nil {
		panic(err)
	}

	log.Println("Listen...")
	server.Serve(unixListener)
}
