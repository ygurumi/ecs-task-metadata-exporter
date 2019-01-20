package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	v2 "github.com/ygurumi/ecs-task-metadata-exporter/apis/v2"
)

type statsCollector struct {
	blkioStats   *prometheus.GaugeVec
	cpuStats     *prometheus.GaugeVec
	preCPUStats  *prometheus.GaugeVec
	memoryStats  *prometheus.GaugeVec
	storageStats *prometheus.GaugeVec
	pidsStats    *prometheus.GaugeVec
	network      *prometheus.GaugeVec
}

func newStatsCollector() statsCollector {
	return statsCollector{
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

func (k statsCollector) setBulkioStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	for key, arr := range c.BlkioStats {
		for _, value := range arr {
			k.blkioStats.WithLabelValues(dockerID, name, key+"/"+value.Op).Set(value.Value)
		}
	}
}

func (k statsCollector) setCPUStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
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

func (k statsCollector) setPreCPUStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
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

func (k statsCollector) setMemoryStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
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

func (k statsCollector) setStorageStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	for key, value := range c.PidsStats {
		k.storageStats.WithLabelValues(dockerID, name, key).Set(value)
	}
}

func (k statsCollector) setPidsStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	for key, value := range c.PidsStats {
		k.pidsStats.WithLabelValues(dockerID, name, key).Set(value)
	}
}

func (k statsCollector) setNetwork(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	for interfaceName, metrics := range c.Network {
		for key, value := range metrics {
			k.network.WithLabelValues(dockerID, name, interfaceName+"/"+key).Set(value)
		}
	}
}

func (k statsCollector) setStats(ch chan<- prometheus.Metric, dockerID string, name string, c *v2.ContainerStats) {
	k.setBulkioStats(ch, dockerID, name, c)
	k.setCPUStats(ch, dockerID, name, c)
	k.setPreCPUStats(ch, dockerID, name, c)
	k.setMemoryStats(ch, dockerID, name, c)
	k.setStorageStats(ch, dockerID, name, c)
	k.setPidsStats(ch, dockerID, name, c)
	k.setNetwork(ch, dockerID, name, c)
}

func (k statsCollector) Describe(ch chan<- *prometheus.Desc) {
	k.blkioStats.Describe(ch)
	k.cpuStats.Describe(ch)
	k.preCPUStats.Describe(ch)
	k.memoryStats.Describe(ch)
	k.storageStats.Describe(ch)
	k.pidsStats.Describe(ch)
	k.network.Describe(ch)
}

func (k statsCollector) Collect(ch chan<- prometheus.Metric) {
	k.blkioStats.Collect(ch)
	k.cpuStats.Collect(ch)
	k.preCPUStats.Collect(ch)
	k.memoryStats.Collect(ch)
	k.storageStats.Collect(ch)
	k.pidsStats.Collect(ch)
	k.network.Collect(ch)
}
