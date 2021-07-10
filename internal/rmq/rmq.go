package rmq

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

func isInCluster(info []clusterinfo.Members, hostname string) (isRunning bool, isAlone bool) {
	runningCnt := 0
	for _, v := range info {
		host := parseHostname(v.Name)
		if hostname == host {
			isRunning = v.Running
		}
		if v.Running {
			runningCnt++
		}
	}

	return isRunning, runningCnt == 1
}

type HealthCheckerConfig struct {
	RabbitMQConfig              rabbitmq.Config
	LastClusterInfoFetchTimeout time.Duration
	HealthCheckMaxDuration      time.Duration
}

func NewHealthHandler(config HealthCheckerConfig, ftch fetcher.Fetcher) func(w http.ResponseWriter, r *http.Request) {
	hc := rabbitmq.New(config.RabbitMQConfig)
	return func(w http.ResponseWriter, r *http.Request) {
		lastFetch := ftch.LastSuccessfulFetch()
		if time.Since(lastFetch) > config.LastClusterInfoFetchTimeout {
			log.Printf("[health-check] failed to perform: time passed since last"+
				"successful cluster info fetch %q is more than %q", lastFetch, config.LastClusterInfoFetchTimeout)
			w.WriteHeader(http.StatusInternalServerError)
		}

		cluster := ftch.GetClusterInfo()
		host, err := os.Hostname()
		if err != nil {
			log.Printf("[health-check] failed to get hostname: %s", err)
			host = ""
		}

		isRunning, isAlone := isInCluster(cluster, host)
		if !isRunning {
			log.Println("[health-check] got node not in cluster or not running")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if isAlone {
			log.Println("[health-check] got alone running node in cluster")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		errCh := make(chan error, 1)
		go func() {
			errCh <- hc(r.Context())

			if err != nil {
				log.Printf("[health-check] failed to check rabbitmq status: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				log.Println("[health-check] successfully checked rabbitmq status")
			}
		}()

		select {
		case <-time.After(config.HealthCheckMaxDuration):
			log.Println("[health-check] failed to check rabbitmq status due to timeout")
			w.WriteHeader(http.StatusInternalServerError)
		case <-r.Context().Done():
			log.Printf(
				"[health-check] failed to check rabbitmq status: client disconnected: %s",
				r.Context().Err(),
			)
			w.WriteHeader(499)
		case err := <-errCh:
			if err != nil {
				log.Printf("[health-check] failed to check rabbitmq status: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				log.Println("[health-check] successfully checked rabbitmq status")
			}
		}
	}
}
