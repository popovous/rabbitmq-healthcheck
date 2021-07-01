package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type fetcher struct {
	clusterMembers []ClusterMembers
}

type ClusterMembers struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
}

type CurrentStatus struct {
	NodesCount, NodesRunning uint16
}

func clusterStatus() []ClusterMembers {
	resp, err := http.Get("")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var nodes []ClusterMembers
	if err := json.Unmarshal(body, &nodes); err != nil {
		log.Fatalln(err)
	}
	return nodes
}

func prettyJson(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	cluster := clusterStatus()
	marshalled, err := json.Marshal(cluster)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Fprintf(w, "Cluste nodes:\n")
	marshalled, _ = prettyJson(marshalled)
	//_, _ = w.Write(marshalled)

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

func handleRequests() {
	http.HandleFunc("/", rootPage)
	http.ListenAndServe(":31337", nil)
}

func main() {
	handleRequests()
}
