package k8s

import (
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

type ClustersListRequest struct {
	Sources []models.Condition `form:"sources" url:"sources"`
}

type ClustersListResponse struct {
	Clusters   []ClustersListResult `json:"clusters"`
	TotalCount int                  `json:"totalCount"`
}

type ClustersListResult struct {
	Name                            string         `json:"name"`
	Env                             string         `json:"env"`
	CpuUsage                        *float64       `json:"cpuUsage"`
	CpuLimit                        *float64       `json:"cpuLimit"`
	CpuAllocatable                  *float64       `json:"cpuAllocatable"`
	CpuRequest                      *float64       `json:"cpuRequest"`
	CpuUsageAllocatablePercent      *float64       `json:"cpuUsageAllocatablePercent"`
	CpuRequestAllocatablePercent    *float64       `json:"cpuRequestAllocatablePercent"`
	CpuUsageRequestPercent          *float64       `json:"cpuUsageRequestPercent"`
	CpuUsageLimitPercent            *float64       `json:"cpuUsageLimitPercent"`
	CpuLimitAllocatablePercent      *float64       `json:"cpuLimitAllocatablePercent"`
	MemoryUsage                     *float64       `json:"memoryUsage"`
	MemoryLimit                     *float64       `json:"memoryLimit"`
	MemoryAllocatable               *float64       `json:"memoryAllocatable"`
	MemoryRequest                   *float64       `json:"memoryRequest"`
	MemoryUsageAllocatablePercent   *float64       `json:"memoryUsageAllocatablePercent"`
	MemoryRequestAllocatablePercent *float64       `json:"memoryRequestAllocatablePercent"`
	MemoryUsageRequestPercent       *float64       `json:"memoryUsageRequestPercent"`
	MemoryUsageLimitPercent         *float64       `json:"memoryUsageLimitPercent"`
	MemoryLimitAllocatablePercent   *float64       `json:"memoryLimitAllocatablePercent"`
	NodesCount                      int            `json:"nodesCount"`
	Pods                            map[string]int `json:"pods"`
	IssueCount                      int            `json:"issueCount"`
	CreationTimestamp               time.Time      `json:"creationTimestamp"`
}

type WorkloadsListRequest struct {
	Sources    []models.Condition `form:"sources" url:"sources"`
	Conditions []models.Condition `form:"conditions" url:"conditions"`
	SortBy     string             `form:"sortBy" url:"sortBy" binding:"oneof=envType env cluster namespace workload kind ready podsCount p50 p95 p99 rps errorRate cpuLimit cpuUsage memoryLimit memoryUsage issueCount"`
	Order      string             `form:"order" url:"order" binding:"oneof=asc desc"`
	Limit      uint32             `form:"limit" url:"limit"`
	Skip       uint32             `form:"skip" url:"skip"`
}

type WorkloadsListResponse struct {
	Total     uint32              `json:"total"`
	Workloads []WorkloadsListItem `json:"workloads"`
}

type WorkloadsListItem struct {
	UID             string   `json:"uid"`
	EnvType         string   `json:"envType"`
	Env             string   `json:"env"`
	Cluster         string   `json:"cluster"`
	Namespace       string   `json:"namespace"`
	Workload        string   `json:"workload"`
	Kind            string   `json:"kind"`
	ResourceVersion int64    `json:"resourceVersion"`
	Ready           bool     `json:"ready"`
	PodsCount       *uint32  `json:"podsCount"`
	P50             *float64 `json:"p50"`
	P95             *float64 `json:"p95"`
	P99             *float64 `json:"p99"`
	RPS             *float64 `json:"rps"`
	ErrorRate       *float64 `json:"errorRate"`
	CPULimit        *float64 `json:"cpuLimit"`
	CPUUsage        *float64 `json:"cpuUsage"`
	MemoryLimit     *float64 `json:"memoryLimit"`
	MemoryUsage     *float64 `json:"memoryUsage"`
	IssueCount      uint32   `json:"issueCount"`
}
