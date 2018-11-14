package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ygurumi/ecs-task-metadata-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultMetadataEndpoint string = "http://169.254.170.2/v2/metadata"
	defaultStatsEndpoint    string = "http://169.254.170.2/v2/stats"
	defaultTimeoutString    string = "500ms"
	defaultAddress          string = ":9887"
	defaultPath             string = "/metrics"
)

func main() {
	var metadataEndpoint, statsEndpoint, timeoutString, address, path string
	flag.StringVar(&metadataEndpoint, "metadata", defaultMetadataEndpoint, "Task metadata endpoint to scrape.")
	flag.StringVar(&statsEndpoint, "stats", defaultStatsEndpoint, "Docker stats endpoint to scrape.")
	flag.StringVar(&timeoutString, "timeout", defaultTimeoutString, "Per-scrape timeout.")
	flag.StringVar(&address, "address", defaultAddress, "Address on which to expose metrics.")
	flag.StringVar(&path, "path", defaultPath, "Path under which to expose metrics.")
	flag.Parse()

	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		log.Fatal(err)
	}

	collector := collector.New(metadataEndpoint, statsEndpoint, timeout)
	if err := prometheus.Register(collector); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintf(w, "%v", path)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "NotFound")
		}
	})

	http.HandleFunc("/debug", collector.DebugHandler)
	http.Handle(path, promhttp.Handler())

	log.Fatal(http.ListenAndServe(address, nil))
}
