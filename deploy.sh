#!/bin/bash
# 在服务器上部署和运行二进制

# system 服务方式
cat <<EOF > /etc/systemd/system/goosefs-cli2api.service
[Unit]
Description=GooseFS CLI to API Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/goosefs-cli2api
Restart=on-failure
RestartSec=5
User=root
Environment=GOOSEFS_ENV=production
WorkingDirectory=/root

[Install]
WantedBy=multi-user.target
EOF

# 启动
systemctl daemon-reload
systemctl restart goosefs-cli2api
systemctl status goosefs-cli2api
systemctl enable goosefs-cli2api