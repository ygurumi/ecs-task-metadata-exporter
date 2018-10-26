package v2

type NetworkMetadata struct {
	NetworkMode   string   `json:"NetworkMode"`
	IPv4Addresses []string `json:"IPv4Addresses"`
}

type ContainerMetadata struct {
	DockerID      string             `json:"DockerId"`
	Name          string             `json:"Name"`
	DockerName    string             `json:"DockerName"`
	Image         string             `json:"Image"`
	ImageID       string             `json:"ImageID"`
	Labels        map[string]string  `json:"Labels"`
	DesiredStatus string             `json:"DesiredStatus"`
	KnownStatus   string             `json:"KnownStatus"`
	Limits        map[string]float64 `json:"Limits"`
	CreatedAt     string             `json:"CreatedAt"`
	StartedAt     string             `json:"StartedAt"`
	Type          string             `json:"Type"`
	Networks      []NetworkMetadata  `json:"Networks"`
}

type TaskMetadata struct {
	Cluster       string              `json:"Cluster"`
	TaskARN       string              `json:"TaskARN"`
	Family        string              `json:"Family"`
	Revision      string              `json:"Revision"`
	DesiredStatus string              `json:"DesiredStatus"`
	KnownStatus   string              `json:"KnownStatus"`
	Containers    []ContainerMetadata `json:"Containers"`
	PullStartedAt string              `json:"PullStartedAt"`
	PullStoppedAt string              `json:"PullStoppedAt"`
}

type BlkioValue struct {
	Major int64   `json:"major"`
	Minor int64   `json:"minor"`
	Op    string  `json:"op"`
	Value float64 `json:"value"`
}

type CPUUsage struct {
	PerCPUUsage       []float64 `json:"percpu_usage"`
	TotalUsage        float64   `json:"total_usage"`
	UsageInKernelmode float64   `json:"usage_in_kernelmode"`
	UsageInUsermode   float64   `json:"usage_in_usermode"`
}

type CPUStats struct {
	CPUUsage       CPUUsage    `json:"cpu_usage"`
	OnLineCPUs     float64     `json:"online_cpus"`
	SystemCPUUsage float64     `json:"system_cpu_usage"`
	ThrottlingData interface{} `json:"throttling_data"`
}

type MemoryStats struct {
	Limit    float64 `json:"limit"`
	MaxUsage float64 `json:"max_usage"`
	Usage    float64 `json:"usage"`
	Stats    map[string]float64
}

type ContainerStats struct {
	BlkioStats   map[string][]BlkioValue       `json:"blkio_stats"`
	CPUStats     CPUStats                      `json:"cpu_stats"`
	PreCPUStats  CPUStats                      `json:"precpu_stats"`
	MemoryStats  MemoryStats                   `json:"memory_stats"`
	StorageStats map[string]float64            `json:"storage_stats"`
	Network      map[string]map[string]float64 `json:"network"`
	NumProcs     float64                       `json:"num_procs"`
	PidsStats    map[string]float64            `json:"pids_status"`
	Read         string                        `json:"read"`
	PreRead      string                        `json:"preread"`
}
