#!/bin/bash
systemctl stop goosefs-cli2api
cp bin/goosefs-cli2api-linux-amd64 /usr/bin/goosefs-cli2api
systemctl restart goosefs-cli2api
systemctl status goosefs-cli2api