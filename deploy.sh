#!/bin/bash
# 在服务器上部署和运行二进制

systemctl stop goosefs-cli2api

# with build
if [$1 = "wb" ]; then
    CGO_ENABLED=1 go build -o /usr/bin/goosefs-cli2api -ldflags "-X main.version=localbuild-beta.$(date +%Y%m%d)"
fi

# system 服务方式
cat <<EOF > /etc/systemd/system/goosefs-cli2api.service
[Unit]
Description=GooseFS CLI to API Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/goosefs-cli2api run
Restart=on-failure
RestartSec=5
User=root
Environment=CGO_ENABLED=1
WorkingDirectory=/root

[Install]
WantedBy=multi-user.target
EOF

# 启动
systemctl daemon-reload
systemctl restart goosefs-cli2api
systemctl status goosefs-cli2api
systemctl enable goosefs-cli2api