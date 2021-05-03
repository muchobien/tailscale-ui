#!/bin/bash

sudo -H -u $SUDO_USER bash -c "XDG_RUNTIME_DIR=/run/user/$(id -u  $SUDO_USER) systemctl --user disable tailscale-ui"
