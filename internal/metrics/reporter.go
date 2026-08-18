// reporter.go 提供指标汇总报告生成和格式化输出。
package metrics

import (
	"fmt"
	"strings"
)

// Report 汇总指标报告。
type Report struct {
	AvgQueueLen     float64 `json:"avg_queue_len"`
	MaxQueueLen     int     `json:"max_queue_len"`
	P95QueueLen     float64 `json:"p95_queue_len"`
	AvgUtilization  float64 `json:"avg_utilization"`
	PeakUtilization float64 `json:"peak_utilization"`
	SampleCount     int     `json:"sample_count"`
}

// GenerateReport 从 Recorder 生成汇总报告。
func GenerateReport(r *Recorder, totalServers int) Report {
	return Report{
		AvgQueueLen:     r.AvgQueueLen(),
		MaxQueueLen:     r.MaxQueueLen(),
		P95QueueLen:     r.QueueLenPercentile(95),
		AvgUtilization:  r.AvgUtilization(totalServers),
		PeakUtilization: r.PeakUtilization(totalServers),
		SampleCount:     r.Count(),
	}
}

// FormatReport 生成人类可读报告文本。
func FormatReport(rpt Report) string {
	var sb strings.Builder
	sb.WriteString("=== Queue Metrics Report ===\n")
	sb.WriteString(fmt.Sprintf("  Samples:          %d\n", rpt.SampleCount))
	sb.WriteString(fmt.Sprintf("  Avg Queue Len:    %.2f\n", rpt.AvgQueueLen))
	sb.WriteString(fmt.Sprintf("  Max Queue Len:    %d\n", rpt.MaxQueueLen))
	sb.WriteString(fmt.Sprintf("  P95 Queue Len:    %.2f\n", rpt.P95QueueLen))
	sb.WriteString(fmt.Sprintf("  Avg Utilization:  %.2f%%\n", rpt.AvgUtilization*100))
	sb.WriteString(fmt.Sprintf("  Peak Utilization: %.2f%%\n", rpt.PeakUtilization*100))
	return sb.String()
}

// SLACheck 检查是否满足 SLA 目标。
type SLAResult struct {
	Target       string  `json:"target"`
	Threshold    float64 `json:"threshold"`
	Actual       float64 `json:"actual"`
	Met          bool    `json:"met"`
}

// CheckSLA 根据阈值检查多个 SLA 指标。
func CheckSLA(rpt Report, maxWaitTarget, utilizationTarget float64) []SLAResult {
	var results []SLAResult
	results = append(results, SLAResult{
		Target:    "avg_queue_len",
		Threshold: maxWaitTarget,
		Actual:    rpt.AvgQueueLen,
		Met:       rpt.AvgQueueLen <= maxWaitTarget,
	})
	results = append(results, SLAResult{
		Target:    "utilization",
		Threshold: utilizationTarget,
		Actual:    rpt.AvgUtilization,
		Met:       rpt.AvgUtilization <= utilizationTarget,
	})
	return results
}
