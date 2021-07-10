package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/popovous/rabbitmq-healthcheck/internal/clusterinfo"
)

const (
	defaultRequestTimeout  = 3 * time.Second
	defaultRefreshInterval = 5 * time.Second
)

type Fetcher interface {
	Start()
	Stop() error
	GetClusterInfo() []clusterinfo.Members
	LastSuccessfulFetch() time.Time
}

type Config struct {
	URL             string
	RefreshInterval time.Duration
	RequestTimeout  time.Duration
}

func (c *Config) withDefaults() *Config {
	if c == nil {
		c = &Config{}
	}
	if c.RequestTimeout == 0 {
		c.RequestTimeout = defaultRequestTimeout
	}
	if c.RefreshInterval == 0 {
		c.RefreshInterval = defaultRefreshInterval
	}

	return c
}

type defaultFetcher struct {
	config              *Config
	data                []clusterinfo.Members
	client              *http.Client
	mu                  sync.RWMutex
	lastSuccessfulFetch time.Time
	onceCloser          sync.Once
	stop                chan struct{}
}

func (d *defaultFetcher) fetch() ([]clusterinfo.Members, error) {
	req, err := http.NewRequest(http.MethodGet, d.config.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare fetcher request: %s", err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from RabbitMQ Management: %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from RabbitMQ Management: %s", err)
	}

	var members []clusterinfo.Members
	err = json.Unmarshal(body, &members)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (d *defaultFetcher) Start() {
	log.Println("starting fetcher")
	t := time.NewTicker(d.config.RefreshInterval)
	go func() {
		defer t.Stop()
		for {
			log.Println("trying to fetch data from RabbitMQ Management")
			members, err := d.fetch()
			if err != nil {
				log.Println(err)
			} else {
				log.Println("successfully fetched data from RabbitMQ Management")

				now := time.Now()
				d.mu.Lock()
				d.data = members
				d.lastSuccessfulFetch = now
				d.mu.Unlock()
			}

			select {
			case <-t.C:
			case <-d.stop:
				return
			}
		}
	}()
}

func (d *defaultFetcher) Stop() error {
	d.onceCloser.Do(func() {
		close(d.stop)
	})

	return nil
}

func (d *defaultFetcher) GetClusterInfo() []clusterinfo.Members {
	d.mu.RLock()
	data := make([]clusterinfo.Members, len(d.data))
	copy(data, d.data)
	d.mu.RUnlock()

	return data
}

func (d *defaultFetcher) LastSuccessfulFetch() time.Time {
	d.mu.RLock()
	defer d.mu.RUnlock()
	tm := d.lastSuccessfulFetch

	return tm
}

func New(config *Config) Fetcher {
	config = config.withDefaults()

	return &defaultFetcher{
		config: config,
		data:   nil,
		client: &http.Client{
			Timeout: config.RequestTimeout,
		},
		stop: make(chan struct{}),
	}
}
