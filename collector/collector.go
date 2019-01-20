package collector

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	v2 "github.com/ygurumi/ecs-task-metadata-exporter/apis/v2"
)

type Collector struct {
	client           http.Client
	metadataEndpoint string
	statsEndpoint    string
	metadata         metadataCollector
	stats            statsCollector
}

func New(metadataEndpoint string, statsEndpoint string, timeout time.Duration) Collector {
	return Collector{
		client: http.Client{
			Timeout: timeout,
		},
		metadataEndpoint: metadataEndpoint,
		statsEndpoint:    statsEndpoint,
		metadata:         newMetadataCollector(),
		stats:            newStatsCollector(),
	}
}

func (k Collector) Describe(ch chan<- *prometheus.Desc) {
	k.metadata.Describe(ch)
	k.stats.Describe(ch)
}

func (k Collector) Collect(ch chan<- prometheus.Metric) {
	var taskMetadata v2.TaskMetadata
	if err := readTaskMetadata(k.client, k.metadataEndpoint, &taskMetadata); err != nil {
		log.Println(err)
		return
	}

	var taskStats map[string]v2.ContainerStats
	if err := readTaskStats(k.client, k.statsEndpoint, &taskStats); err != nil {
		log.Println(err)
		return
	}

	k.metadata.setTaskMetadata(ch, &taskMetadata)
	nameMap := map[string]string{}
	for _, containerMetadata := range taskMetadata.Containers {
		k.metadata.setContainerMetadata(ch, &containerMetadata)
		nameMap[containerMetadata.DockerID] = containerMetadata.Name
	}

	for dockerID, containerStats := range taskStats {
		name := nameMap[dockerID]
		k.stats.setStats(ch, dockerID, name, &containerStats)
	}

	k.metadata.Collect(ch)
	k.stats.Collect(ch)
}

func (k Collector) DebugHandler(w http.ResponseWriter, r *http.Request) {
	var taskMetadata interface{}
	if err := readTaskMetadata(k.client, k.metadataEndpoint, &taskMetadata); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%v", err)
		return
	}

	var taskStats interface{}
	if err := readTaskStats(k.client, k.statsEndpoint, &taskStats); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if err := fPrettyPrint(w, []interface{}{
		taskMetadata,
		taskStats,
	}); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%v", err)
		return
	}
}
