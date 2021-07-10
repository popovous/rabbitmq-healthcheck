package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/hellofresh/health-go/v4/checks/rabbitmq"
	"github.com/popovous/rabbitmq-healthcheck/internal/rmq"

	"github.com/popovous/rabbitmq-healthcheck/internal/fetcher"
)

var (
	fetcherURL = flag.String("fetcher.url", "", "RabbitMQ Managment URL.")
	amqpDSN    = flag.String("amqp.url", "", "AMQP URL.")
	listenAddr = flag.String("listen.addr", "", "AMQP URL.")
)

func handleRequests(f fetcher.Fetcher) {
	http.HandleFunc("/", rmq.NewHealthHandler(rmq.HealthCheckerConfig{
		RabbitMQConfig: rabbitmq.Config{
			DialTimeout: 200 * time.Millisecond,
			DSN:         *amqpDSN,
		},
		LastClusterInfoFetchTimeout: 60 * time.Second,
		HealthCheckMaxDuration:      200 * time.Millisecond,
	}, f))
	http.ListenAndServe(*listenAddr, nil)
}

func main() {
	flag.Parse()

	f := fetcher.New(&fetcher.Config{
		URL: *fetcherURL,
	})
	f.Start()

	handleRequests(f)
}
