package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ygurumi/ecs-task-metadata-exporter/apis/v2"
)

type metadataCollector struct {
	taskMetadata      *prometheus.GaugeVec
	containerMetadata *prometheus.GaugeVec
	containerLabels   *prometheus.GaugeVec
	containerLimits   *prometheus.GaugeVec
}

func newMetadataCollector() metadataCollector {
	return metadataCollector{
		taskMetadata: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: taskSubsystem,
				Name:      "metadata",
				Help:      "metadata",
			},
			[]string{"task_arn", "cluster", "familiy", "revision"},
		),

		containerMetadata: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "metadata",
				Help:      "metadata",
			},
			[]string{"docker_id", "name", "docker_name", "image", "image_id", "type"},
		),

		containerLabels: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "labels",
				Help:      "labels",
			},
			[]string{"docker_id", "name", "key", "value"},
		),

		containerLimits: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "limits",
				Help:      "limits",
			},
			[]string{"docker_id", "name", "resource"},
		),
	}
}

func (k metadataCollector) setTaskMetadata(ch chan<- prometheus.Metric, t *v2.TaskMetadata) {
	k.taskMetadata.WithLabelValues(t.TaskARN, t.Cluster, t.Family, t.Revision).Set(0)
}

func (k metadataCollector) setContainerMetadata(ch chan<- prometheus.Metric, c *v2.ContainerMetadata) {
	k.containerMetadata.WithLabelValues(c.DockerID, c.Name, c.DockerName, c.Image, c.ImageID, c.Type).Set(0)

	for key, value := range c.Labels {
		k.containerLabels.WithLabelValues(c.DockerID, c.Name, key, value).Set(0)
	}

	for key, value := range c.Limits {
		k.containerLimits.WithLabelValues(c.DockerID, c.Name, key).Set(value)
	}
}

func (k metadataCollector) Describe(ch chan<- *prometheus.Desc) {
	k.taskMetadata.Describe(ch)
	k.containerMetadata.Describe(ch)
	k.containerLabels.Describe(ch)
	k.containerLimits.Describe(ch)
}

func (k metadataCollector) Collect(ch chan<- prometheus.Metric) {
	k.taskMetadata.Collect(ch)
	k.containerMetadata.Collect(ch)
	k.containerLabels.Collect(ch)
	k.containerLimits.Collect(ch)
}
