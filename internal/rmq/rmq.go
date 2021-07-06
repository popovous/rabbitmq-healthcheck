package rmq

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/popovous/rabbitmq-healthcheck/internal/clusterinfo"

	"github.com/popovous/rabbitmq-healthcheck/internal/fetcher"

	"github.com/hellofresh/health-go/v4/checks/rabbitmq"
)

func parseHostname(nodeName string) string {
	sp := strings.Split(nodeName, "@")
	if len(sp) == 0 || len(sp) == 1 {
		return ""
	}
	return sp[1]
}

func isInCluster(info []clusterinfo.Members, hostname string) bool {
	for _, v := range info {
		host := parseHostname(v.Name)
		if hostname == host {
			return true
		}
	}
	return false
}

func NewHealthHandler(config rabbitmq.Config, ftch fetcher.Fetcher) func(w http.ResponseWriter, r *http.Request) {
	hc := rabbitmq.New(config)
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := ftch.GetClusterInfo()
		host, err := os.Hostname()
		if err != nil {
			log.Printf("[health-check] failed to get hostname")
			host = ""
		}

		inCluster := isInCluster(cluster, host)
		if !inCluster {
			log.Printf("[health-check] got node not in cluster")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = hc(r.Context())
		if err != nil {
			log.Printf("[health-check] failed to check rabbitmq status: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Println("[health-check] successfully checked rabbitmq status")
		}
	}
}
