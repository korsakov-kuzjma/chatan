#!/bin/bash

# Установка зависимостей
sudo apt update
sudo apt install -y git golang

# Настройка проекта
cd ~/chatan
cp config.yaml.example config.yaml
go mod tidy
go build -o chatan cmd/chatan/main.go

# Создание systemd службы
sudo tee /etc/systemd/system/chatan.service >/dev/null <<EOF
[Unit]
Description=ChatAN Telegram Bot
After=network.target

[Service]
User=kkorsakov
WorkingDirectory=/home/kkorsakov/chatan
ExecStart=/home/kkorsakov/chatan/chatan
Restart=always

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable chatan
