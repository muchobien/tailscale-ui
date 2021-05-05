package daemon

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"tailscale.com/client/tailscale"
	"tailscale.com/ipn"
)

const SocketPath = "/var/run/tailscale-ui.sock"

func connect(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	st, err := tailscale.Status(ctx)
	if err != nil {
		http.Error(w, "Error fetching current status", http.StatusInternalServerError)
	}
	if st.BackendState == "Running" {
		http.Error(w, "Tailscale was already Running", http.StatusConflict)
	}
	_, err = tailscale.EditPrefs(ctx, &ipn.MaskedPrefs{
		Prefs: ipn.Prefs{
			WantRunning: true,
		},
		WantRunningSet: true,
	})

	if err != nil {
		http.Error(w, "Failed to Connect", http.StatusInternalServerError)
	}
}

func disconnect(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	st, err := tailscale.Status(ctx)
	if err != nil {
		http.Error(w, "Error fetching current status", http.StatusInternalServerError)
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
	}
}

func Listen() {

	os.Remove(SocketPath)

	server := http.Server{}

	http.HandleFunc("/connect", connect)
	http.HandleFunc("/disconnect", disconnect)

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
