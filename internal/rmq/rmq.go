package rmq

import (
	"log"
	"net/http"

	"github.com/hellofresh/health-go/v4/checks/rabbitmq"
)

func NewHealthHandler(config rabbitmq.Config) func(w http.ResponseWriter, r *http.Request) {
	hc := rabbitmq.New(config)
	return func(w http.ResponseWriter, r *http.Request) {
		err := hc(r.Context())
		if err != nil {
			log.Printf("[health-check] failed to check rabbitmq status: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Println("[health-check] successfully checked rabbitmq status")
		}
	}
}
