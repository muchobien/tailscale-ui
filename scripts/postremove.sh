#!/bin/bash
systemctl disable tailscale-ui-daemon 
  
sudo -H -u $SUDO_USER bash -c "XDG_RUNTIME_DIR=/run/user/$(id -u  $SUDO_USER) systemctl --user stop tailscale-ui"
sudo -H -u $SUDO_USER bash -c "XDG_RUNTIME_DIR=/run/user/$(id -u  $SUDO_USER) systemctl --user disable tailscale-ui"
sudo -H -u $SUDO_USER bash -c "XDG_RUNTIME_DIR=/run/user/$(id -u  $SUDO_USER) systemctl --user daemon-reload"
