package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// ServerStatusController 控制器结构体
type ServerStatusController struct{}

// ServerStatus 服务器状态结构体
type ServerStatus struct {
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsage    float64 `json:"memory_usage"`
	DiskUsage      float64 `json:"disk_usage"`
	BandwidthUsage float64 `json:"bandwidth_usage"`
	NetworkLatency float64 `json:"network_latency"`
}

// GetServerStatus 获取服务器状态信息
func (s *ServerStatusController) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 获取服务器状态
	status, err := s.fetchServerStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将数据转换为 JSON 格式并写入响应
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// fetchServerStatus 实际获取服务器状态数据
func (s *ServerStatusController) fetchServerStatus() (ServerStatus, error) {
	// 获取 CPU 使用率
	cpuPercentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return ServerStatus{}, err
	}
	cpuUsage := cpuPercentages[0]

	// 获取内存使用率
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return ServerStatus{}, err
	}
	memoryUsage := vmStat.UsedPercent

	// 获取磁盘使用率
	diskStat, err := disk.Usage("/")
	if err != nil {
		return ServerStatus{}, err
	}
	diskUsage := diskStat.UsedPercent

	// 获取带宽使用情况（假设获取所有网卡的总带宽使用率）
	netIOs, err := net.IOCounters(false)
	if err != nil {
		return ServerStatus{}, err
	}
	bandwidthUsage := float64(netIOs[0].BytesSent+netIOs[0].BytesRecv) / (1024 * 1024) // 转换为 MB

	// 模拟网络延迟（根据需要替换为实际数据）
	networkLatency := 20.0 // 单位：毫秒

	// 返回服务器状态
	return ServerStatus{
		CPUUsage:       cpuUsage,
		MemoryUsage:    memoryUsage,
		DiskUsage:      diskUsage,
		BandwidthUsage: bandwidthUsage,
		NetworkLatency: networkLatency,
	}, nil
}
