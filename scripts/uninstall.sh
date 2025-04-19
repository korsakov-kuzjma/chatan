#!/bin/bash

sudo systemctl stop chatan
sudo systemctl disable chatan
sudo rm /etc/systemd/system/chatan.service
sudo systemctl daemon-reload

rm -rf ~/chatan
