package collector

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ygurumi/ecs-task-metadata-exporter/apis/v2"
)

const (
	namespace          string = "ecs"
	containerSubsystem string = "container"
	taskSubsystem      string = "task"
)

type Collector struct {
	client            http.Client
	metadataEndpoint  string
	statsEndpoint     string
	taskMetadata      *prometheus.GaugeVec
	containerMetadata *prometheus.GaugeVec
	containerLabels   *prometheus.GaugeVec
	containerLimits   *prometheus.GaugeVec
	blkioStats        *prometheus.GaugeVec
	cpuStats          *prometheus.GaugeVec
	preCPUStats       *prometheus.GaugeVec
	memoryStats       *prometheus.GaugeVec
	storageStats      *prometheus.GaugeVec
	pidsStats         *prometheus.GaugeVec
	network           *prometheus.GaugeVec
}

func New(metadataEndpoint string, statsEndpoint string, timeout time.Duration) Collector {
	return Collector{
		client: http.Client{
			Timeout: timeout,
		},

		metadataEndpoint: metadataEndpoint,

		statsEndpoint: statsEndpoint,

		taskMetadata: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: taskSubsystem,
				Name:      "metadata",
				Help:      "metadata",
			},
			[]string{"cluster", "familiy", "task_arn", "revision", "desired_status", "known_status"},
		),

		containerMetadata: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "metadata",
				Help:      "metadata",
			},
			[]string{"docker_id", "docker_name", "image", "image_id", "name", "desired_status", "known_status", "type"},
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

		blkioStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "bulkio_stats",
				Help:      "bulkio_stats",
			},
			[]string{"docker_id", "name", "path"},
		),

		cpuStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "cpu_stats",
				Help:      "cpu_stats",
			},
			[]string{"docker_id", "name", "path"},
		),

		preCPUStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "precpu_stats",
				Help:      "precpu_stats",
			},
			[]string{"docker_id", "name", "path"},
		),

		memoryStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "memory_stats",
				Help:      "memory_stats",
			},
			[]string{"docker_id", "name", "path"},
		),

		storageStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "storage_stats",
				Help:      "storage_stats",
			},
			[]string{"docker_id", "name", "path"},
		),

		pidsStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "pids_stats",
				Help:      "pids_stats",
			},
			[]string{"docker_id", "name", "path"},
		),

		network: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: containerSubsystem,
				Name:      "network",
				Help:      "network",
			},
			[]string{"docker_id", "name", "path"},
		),
	}
}

func (k Collector) setTaskMetadata(ch chan<- prometheus.Metric, t *v2.TaskMetadata) {
	k.taskMetadata.WithLabelValues(t.Cluster, t.Family, t.TaskARN, t.Revision, t.DesiredStatus, t.KnownStatus).Set(0)
}

func (k Collector) setContainerMetadata(ch chan<- prometheus.Metric, c *v2.ContainerMetadata) {
	k.containerMetadata.WithLabelValues(c.DockerID, c.DockerName, c.Image, c.ImageID, c.Name, c.DesiredStatus, c.KnownStatus, c.Type).Set(0)
}

func (k Collector) setContainerLabels(ch chan<- prometheus.Metric, c *v2.ContainerMetadata) {
	for key, value := range c.Labels {
		k.containerLabels.WithLabelValues(c.DockerID, c.Name, key, value).Set(0)
	}
}

func (k Collector) setContainerLimits(ch chan<- prometheus.Metric, c *v2.ContainerMetadata) {
	for key, value := range c.Limits {
		k.containerLimits.WithLabelValues(c.DockerID, c.Name, key).Set(value)
	}
}

func (k Collector) setBulkioStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	for key, arr := range c.BlkioStats {
		for _, value := range arr {
			k.blkioStats.WithLabelValues(dockerID, name, key+"/"+value.Op).Set(value.Value)
		}
	}
}

func (k Collector) setCPUStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	metrics := map[string]float64{
		"online_cpus":                   c.CPUStats.OnLineCPUs,
		"system_cpu_usage":              c.CPUStats.SystemCPUUsage,
		"cpu_usage/total_usage":         c.CPUStats.CPUUsage.TotalUsage,
		"cpu_usage/usage_in_kernelmode": c.CPUStats.CPUUsage.UsageInKernelmode,
		"cpu_usage/usage_in_usermode":   c.CPUStats.CPUUsage.UsageInUsermode,
	}

	for key, value := range metrics {
		k.cpuStats.WithLabelValues(dockerID, name, key).Set(value)
	}

	for i, value := range c.CPUStats.CPUUsage.PerCPUUsage {
		k.cpuStats.WithLabelValues(dockerID, name, fmt.Sprintf("cpu_usage/percpu_usage/%v", i)).Set(value)
	}
}

func (k Collector) setPreCPUStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	metrics := map[string]float64{
		"online_cpus":                   c.PreCPUStats.OnLineCPUs,
		"system_cpu_usage":              c.PreCPUStats.SystemCPUUsage,
		"cpu_usage/total_usage":         c.PreCPUStats.CPUUsage.TotalUsage,
		"cpu_usage/usage_in_kernelmode": c.PreCPUStats.CPUUsage.UsageInKernelmode,
		"cpu_usage/usage_in_usermode":   c.PreCPUStats.CPUUsage.UsageInUsermode,
	}

	for key, value := range metrics {
		k.preCPUStats.WithLabelValues(dockerID, name, key).Set(value)
	}

	for i, value := range c.PreCPUStats.CPUUsage.PerCPUUsage {
		k.preCPUStats.WithLabelValues(dockerID, name, fmt.Sprintf("cpu_usage/percpu_usage/%v", i)).Set(value)
	}
}

func (k Collector) setMemoryStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	metrics := map[string]float64{
		"limit":     c.MemoryStats.Limit,
		"max_usage": c.MemoryStats.MaxUsage,
		"usage":     c.MemoryStats.Usage,
	}

	for key, value := range metrics {
		k.memoryStats.WithLabelValues(dockerID, name, key).Set(value)
	}

	for key, value := range c.MemoryStats.Stats {
		k.memoryStats.WithLabelValues(dockerID, name, "stats/"+key).Set(value)
	}
}

func (k Collector) setStorageStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	for key, value := range c.PidsStats {
		k.storageStats.WithLabelValues(dockerID, name, key).Set(value)
	}
}

func (k Collector) setPidsStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	for key, value := range c.PidsStats {
		k.pidsStats.WithLabelValues(dockerID, name, key).Set(value)
	}
}

func (k Collector) setNetwork(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	for interfaceName, metrics := range c.Network {
		for key, value := range metrics {
			k.network.WithLabelValues(dockerID, name, interfaceName+"/"+key).Set(value)
		}
	}
}

func (k Collector) Describe(ch chan<- *prometheus.Desc) {
	k.taskMetadata.Describe(ch)
	k.containerMetadata.Describe(ch)
	k.containerLabels.Describe(ch)
	k.containerLimits.Describe(ch)
	k.blkioStats.Describe(ch)
	k.cpuStats.Describe(ch)
	k.preCPUStats.Describe(ch)
	k.memoryStats.Describe(ch)
	k.storageStats.Describe(ch)
	k.pidsStats.Describe(ch)
	k.network.Describe(ch)
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

	k.setTaskMetadata(ch, &taskMetadata)

	nameMap := map[string]string{}
	for _, containerMetadata := range taskMetadata.Containers {
		k.setContainerMetadata(ch, &containerMetadata)
		k.setContainerLabels(ch, &containerMetadata)
		k.setContainerLimits(ch, &containerMetadata)
		nameMap[containerMetadata.DockerID] = containerMetadata.Name
	}

	for dockerID, containerStats := range taskStats {
		name := nameMap[dockerID]
		k.setBulkioStats(ch, dockerID, name, &containerStats)
		k.setCPUStats(ch, dockerID, name, &containerStats)
		k.setPreCPUStats(ch, dockerID, name, &containerStats)
		k.setMemoryStats(ch, dockerID, name, &containerStats)
		k.setStorageStats(ch, dockerID, name, &containerStats)
		k.setPidsStats(ch, dockerID, name, &containerStats)
		k.setNetwork(ch, dockerID, name, &containerStats)
	}

	k.taskMetadata.Collect(ch)
	k.containerMetadata.Collect(ch)
	k.containerLabels.Collect(ch)
	k.containerLimits.Collect(ch)
	k.blkioStats.Collect(ch)
	k.cpuStats.Collect(ch)
	k.preCPUStats.Collect(ch)
	k.memoryStats.Collect(ch)
	k.storageStats.Collect(ch)
	k.pidsStats.Collect(ch)
	k.network.Collect(ch)
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
