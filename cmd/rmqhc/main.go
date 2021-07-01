package main

import (
	"flag"
	"html/template"
	"net/http"

	"github.com/popovous/rabbitmq-healthcheck/internal/fetcher"
)

var (
	fetcherURL = flag.String("fetcher.url", "", "Remote Cloud URL.")
)

type CurrentStatus struct {
	NodesCount, NodesRunning uint16
}

func newRootPageHandler(fetcher fetcher.Fetcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cluster := fetcher.GetClusterInfo()

		nodesRunning := 0
		nodesCount := 0

		for _, v := range cluster {
			//fmt.Fprintf(w, v.Name+"\n")
			nodesCount++
			if v.Running {
				nodesRunning++
			}
		}

		if nodesRunning == 1 {
			w.WriteHeader(http.StatusInternalServerError)
		}

		htmlValues := CurrentStatus{uint16(nodesCount), uint16(nodesRunning)}
		tmpl, _ := template.ParseFiles("templates/main.html")
		tmpl.Execute(w, htmlValues)
	}
}

func handleRequests(f fetcher.Fetcher) {
	http.HandleFunc("/", newRootPageHandler(f))
	http.ListenAndServe(":31337", nil)
}

func main() {
	flag.Parse()

	f := fetcher.New(&fetcher.Config{
		URL: *fetcherURL,
	})
	f.Start()

	handleRequests(f)
}
