#!/bin/sh
DIR=`mktemp -d`
PRJ_NAME=rabbitmq-healthcheck
RRJ_VER=`date +%H%M%d%m%Y`
mkdir -p ${DIR}/${PRJ_NAME}/usr/bin ${DIR}/${PRJ_NAME}/DEBIAN ${DIR}/${PRJ_NAME}/etc/systemd/system

cat << EOF >> ${DIR}/${PRJ_NAME}/DEBIAN/control
Source: rabbitmq-healthcheck
Description: Advanced RabbitMQ HealthCheck
Architecture: amd64
Package: rabbitmq-healthcheck
Homepage: https://github.com/popovous/rabbitmq-healthcheck
Multi-Arch: foreign
Maintainer: $(git config --global user.email)
Version: 0.1
Priority: optional
EOF

cat << EOF >> ${DIR}/${PRJ_NAME}/etc/systemd/system/${PRJ_NAME}.service
[Unit]
Description=Advanced RabbitMQ HealthCheck

[Service]
ExecStart=/usr/bin/rabbitmq-healthcheck --fetcher.url="http://admin:xxx@127.0.0.1:15672/api/nodes" --amqp.url="amqp://admin:xxx@127.0.0.1" --listen.addr=":8080"
User=nobody

[Install]
WantedBy=multi-user.target
EOF

go build -o ${DIR}/${PRJ_NAME}/usr/bin/rabbitmq-healthcheck cmd/rmqhc/main.go
dpkg-deb -Z xz -S extreme -b ${DIR}/${PRJ_NAME} ${DIR}/${PRJ_NAME}-${RRJ_VER}.deb
dpkg-deb -c ${DIR}/${PRJ_NAME}-${RRJ_VER}.deb
dpkg-deb -W ${DIR}/${PRJ_NAME}-${RRJ_VER}.deb
dpkg-deb -I ${DIR}/${PRJ_NAME}-${RRJ_VER}.deb
