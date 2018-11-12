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

type StatsCollector struct {
	client       http.Client
	endpoint     string
	blkioStats   *prometheus.GaugeVec
	cpuStats     *prometheus.GaugeVec
	preCPUStats  *prometheus.GaugeVec
	memoryStats  *prometheus.GaugeVec
	storageStats *prometheus.GaugeVec
	pidsStats    *prometheus.GaugeVec
	network      *prometheus.GaugeVec
}

func NewStatsCollector(endpoint string, timeout time.Duration) StatsCollector {
	return StatsCollector{
		client: http.Client{
			Timeout: timeout,
		},
		endpoint: endpoint,
		blkioStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: ContainerSubsystem,
				Name:      "bulkio_stats",
				Help:      "cpu_stats",
			},
			[]string{"docker_id", "path"},
		),
		cpuStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: ContainerSubsystem,
				Name:      "cpu_stats",
				Help:      "cpu_stats",
			},
			[]string{"docker_id", "path"},
		),
		preCPUStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: ContainerSubsystem,
				Name:      "precpu_stats",
				Help:      "precpu_stats",
			},
			[]string{"docker_id", "path"},
		),
		memoryStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: ContainerSubsystem,
				Name:      "memory_stats",
				Help:      "memory_stats",
			},
			[]string{"docker_id", "path"},
		),
		storageStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: ContainerSubsystem,
				Name:      "storage_stats",
				Help:      "storage_stats",
			},
			[]string{"docker_id", "path"},
		),
		pidsStats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: ContainerSubsystem,
				Name:      "pids_stats",
				Help:      "pids_stats",
			},
			[]string{"docker_id", "path"},
		),
		network: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: ContainerSubsystem,
				Name:      "network",
				Help:      "network",
			},
			[]string{"docker_id", "path"},
		),
	}
}

func (s StatsCollector) SetBulkioStats(ch chan<- prometheus.Metric, dockerID string, c *v2.ContainerStats) {
	for key, arr := range c.BlkioStats {
		for _, value := range arr {
			s.blkioStats.WithLabelValues(dockerID, key+"/"+value.Op).Set(value.Value)
		}
	}
}

func (s StatsCollector) SetCPUStats(ch chan<- prometheus.Metric, dockerID string, c *v2.ContainerStats) {
	metrics := map[string]float64{
		"online_cpus":                   c.CPUStats.OnLineCPUs,
		"system_cpu_usage":              c.CPUStats.SystemCPUUsage,
		"cpu_usage/total_usage":         c.CPUStats.CPUUsage.TotalUsage,
		"cpu_usage/usage_in_kernelmode": c.CPUStats.CPUUsage.UsageInKernelmode,
		"cpu_usage/usage_in_usermode":   c.CPUStats.CPUUsage.UsageInUsermode,
	}

	for key, value := range metrics {
		s.cpuStats.WithLabelValues(dockerID, key).Set(value)
	}

	for i, value := range c.CPUStats.CPUUsage.PerCPUUsage {
		s.cpuStats.WithLabelValues(dockerID, fmt.Sprintf("cpu_usage/percpu_usage/%v", i)).Set(value)
	}
}

func (s StatsCollector) SetPreCPUStats(ch chan<- prometheus.Metric, dockerID string, c *v2.ContainerStats) {
	metrics := map[string]float64{
		"online_cpus":                   c.PreCPUStats.OnLineCPUs,
		"system_cpu_usage":              c.PreCPUStats.SystemCPUUsage,
		"cpu_usage/total_usage":         c.PreCPUStats.CPUUsage.TotalUsage,
		"cpu_usage/usage_in_kernelmode": c.PreCPUStats.CPUUsage.UsageInKernelmode,
		"cpu_usage/usage_in_usermode":   c.PreCPUStats.CPUUsage.UsageInUsermode,
	}

	for key, value := range metrics {
		s.preCPUStats.WithLabelValues(dockerID, key).Set(value)
	}

	for i, value := range c.PreCPUStats.CPUUsage.PerCPUUsage {
		s.preCPUStats.WithLabelValues(dockerID, fmt.Sprintf("cpu_usage/percpu_usage/%v", i)).Set(value)
	}
}

func (s StatsCollector) SetMemoryStats(ch chan<- prometheus.Metric, dockerID string, c *v2.ContainerStats) {
	metrics := map[string]float64{
		"limit":     c.MemoryStats.Limit,
		"max_usage": c.MemoryStats.MaxUsage,
		"usage":     c.MemoryStats.Usage,
	}

	for key, value := range metrics {
		s.memoryStats.WithLabelValues(dockerID, key).Set(value)
	}

	for key, value := range c.MemoryStats.Stats {
		s.memoryStats.WithLabelValues(dockerID, "stats/"+key).Set(value)
	}
}

func (s StatsCollector) SetStorageStats(ch chan<- prometheus.Metric, dockerID string, c *v2.ContainerStats) {
	for key, value := range c.PidsStats {
		s.storageStats.WithLabelValues(dockerID, key).Set(value)
	}
}

func (s StatsCollector) SetPidsStats(ch chan<- prometheus.Metric, dockerID string, c *v2.ContainerStats) {
	for key, value := range c.PidsStats {
		s.pidsStats.WithLabelValues(dockerID, key).Set(value)
	}
}

func (s StatsCollector) SetNetwork(ch chan<- prometheus.Metric, dockerID string, c *v2.ContainerStats) {
	for interfaceName, metrics := range c.Network {
		for key, value := range metrics {
			s.network.WithLabelValues(dockerID, interfaceName+"/"+key).Set(value)
		}
	}
}

func (s StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	s.blkioStats.Describe(ch)
	s.cpuStats.Describe(ch)
	s.preCPUStats.Describe(ch)
	s.memoryStats.Describe(ch)
	s.storageStats.Describe(ch)
	s.pidsStats.Describe(ch)
	s.network.Describe(ch)
}

func (s StatsCollector) Collect(ch chan<- prometheus.Metric) {
	bs, err := GetHTTPBytes(s.client, s.endpoint)
	if err != nil {
		log.Println(err)
		return
	}

	taskStats := map[string]v2.ContainerStats{}
	if err := json.Unmarshal(bs, &taskStats); err != nil {
		log.Println(err)
		return
	}

	for dockerID, containerStats := range taskStats {
		s.SetBulkioStats(ch, dockerID, &containerStats)
		s.SetCPUStats(ch, dockerID, &containerStats)
		s.SetPreCPUStats(ch, dockerID, &containerStats)
		s.SetMemoryStats(ch, dockerID, &containerStats)
		s.SetStorageStats(ch, dockerID, &containerStats)
		s.SetPidsStats(ch, dockerID, &containerStats)
		s.SetNetwork(ch, dockerID, &containerStats)
	}

	s.blkioStats.Collect(ch)
	s.cpuStats.Collect(ch)
	s.preCPUStats.Collect(ch)
	s.memoryStats.Collect(ch)
	s.storageStats.Collect(ch)
	s.pidsStats.Collect(ch)
	s.network.Collect(ch)
}

func (s StatsCollector) DebugHandler(w http.ResponseWriter, r *http.Request) {
	bs, err := GetHTTPBytes(s.client, s.endpoint)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v", string(bs))
}
