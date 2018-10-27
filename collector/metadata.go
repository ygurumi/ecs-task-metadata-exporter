package collector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ygurumi/ecs-task-metadata-exporter/apis/v2"
)

var (
	taskMetadataDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, TaskSubsystem, "metadata"),
		"task metadata",
		[]string{"cluster", "familiy", "task_arn", "revision", "desired_status", "known_status"},
		nil,
	)

	containerMetadataDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, ContainerSubsystem, "metadata"),
		"container metadata",
		[]string{"docker_id", "docker_name", "image", "image_id", "name", "desired_status", "known_status", "type"},
		nil,
	)

	containerLabelsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, ContainerSubsystem, "labels"),
		"container labels",
		[]string{"docker_id", "key", "value"},
		nil,
	)

	containerLimitsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, ContainerSubsystem, "limits"),
		"container limits",
		[]string{"docker_id", "resource"},
		nil,
	)
)

type MetadataCollector struct {
	client   http.Client
	endpoint string
}

func NewMetadataCollector(endpoint string, timeout time.Duration) MetadataCollector {
	return MetadataCollector{
		client: http.Client{
			Timeout: timeout,
		},
		endpoint: endpoint,
	}
}

func (m MetadataCollector) PutTaskMetadata(ch chan<- prometheus.Metric, t *v2.TaskMetadata) error {
	metric, err := prometheus.NewConstMetric(taskMetadataDesc, prometheus.GaugeValue, 0, t.Cluster, t.Family, t.TaskARN, t.Revision, t.DesiredStatus, t.KnownStatus)
	if err != nil {
		return err
	}

	ch <- metric
	return nil
}

func (m MetadataCollector) PutContainerMetadata(ch chan<- prometheus.Metric, c *v2.ContainerMetadata) error {
	metric, err := prometheus.NewConstMetric(containerMetadataDesc, prometheus.GaugeValue, 0, c.DockerID, c.DockerName, c.Image, c.ImageID, c.Name, c.DesiredStatus, c.KnownStatus, c.Type)
	if err != nil {
		return err
	}

	ch <- metric
	return nil
}

func (m MetadataCollector) PutContainerLabels(ch chan<- prometheus.Metric, c *v2.ContainerMetadata) error {
	for key, value := range c.Labels {
		m, err := prometheus.NewConstMetric(containerLabelsDesc, prometheus.GaugeValue, 0, c.DockerID, key, value)
		if err != nil {
			return err
		}
		ch <- m
	}

	return nil
}

func (m MetadataCollector) PutContainerLimits(ch chan<- prometheus.Metric, c *v2.ContainerMetadata) error {
	for key, value := range c.Limits {
		m, err := prometheus.NewConstMetric(containerLimitsDesc, prometheus.GaugeValue, value, c.DockerID, key)
		if err != nil {
			return err
		}
		ch <- m
	}

	return nil
}

func (m MetadataCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- taskMetadataDesc
	ch <- containerMetadataDesc
	ch <- containerLabelsDesc
	ch <- containerLimitsDesc
}

func (m MetadataCollector) Collect(ch chan<- prometheus.Metric) {
	bs, err := GetHTTPBytes(m.client, m.endpoint)
	if err != nil {
		log.Println(err)
		return
	}

	taskMetadata := v2.TaskMetadata{}
	if err := json.Unmarshal(bs, &taskMetadata); err != nil {
		log.Println(err)
		return
	}

	if err := m.PutTaskMetadata(ch, &taskMetadata); err != nil {
		log.Println(err)
		return
	}

	for _, containerMetadata := range taskMetadata.Containers {
		if err := m.PutContainerMetadata(ch, &containerMetadata); err != nil {
			log.Println(err)
			return
		}
		if err := m.PutContainerLabels(ch, &containerMetadata); err != nil {
			log.Println(err)
			return
		}
		if err := m.PutContainerLimits(ch, &containerMetadata); err != nil {
			log.Println(err)
			return
		}
	}
}

func (m MetadataCollector) DebugHandler(w http.ResponseWriter, r *http.Request) {
	bs, err := GetHTTPBytes(m.client, m.endpoint)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v", string(bs))
}
