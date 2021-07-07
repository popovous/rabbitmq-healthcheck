# rabbitmq-healthcheck
Development running:
  go run cmd/rmqhc/main.go --fetcher.url="http://admin:xxx@127.0.0.1:15672/api/nodes" --amqp.url="amqp://admin:xxx@127.0.0.1" --listen.addr=":8080"

Building:
  go build -o rabbitmq-healthcheck cmd/rmqhc/main.go

Deb package creating:
  scripts/make_deb.sh
