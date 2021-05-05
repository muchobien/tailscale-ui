#!/bin/bash
systemctl daemon-reload
systemctl enable tailscale-ui-daemon
systemctl start tailscale-ui-daemon

sudo -H -u $SUDO_USER bash -c "XDG_RUNTIME_DIR=/run/user/$(id -u  $SUDO_USER) systemctl --user daemon-reload"
sudo -H -u $SUDO_USER bash -c "XDG_RUNTIME_DIR=/run/user/$(id -u  $SUDO_USER) systemctl --user enable tailscale-ui"
sudo -H -u $SUDO_USER bash -c "XDG_RUNTIME_DIR=/run/user/$(id -u  $SUDO_USER) systemctl --user start tailscale-ui"
exit 0
