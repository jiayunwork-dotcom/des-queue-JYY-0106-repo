// Package metrics 提供仿真过程中的瞬时指标记录器。在仿真运行时按固定
// 间隔采样系统状态（队列长度、忙服务台数、等待时间），用于分析系统随时
// 间的瞬时行为而非仅输出稳态均值。
package metrics

import (
	"math"
	"sort"
)

// Sample 一个采样点。
type Sample struct {
	Time      float64 `json:"time"`
	QueueLen  int     `json:"queue_len"`
	BusyCount int     `json:"busy_count"`
	IdleCount int     `json:"idle_count"`
	InSystem  int     `json:"in_system"`
}

// Recorder 瞬时指标记录器。
type Recorder struct {
	Samples  []Sample
	Interval float64 // 采样间隔
}

// NewRecorder 以采样间隔 interval 创建记录器。
func NewRecorder(interval float64) *Recorder {
	if interval <= 0 {
		interval = 1.0
	}
	return &Recorder{Interval: interval}
}

// Record 记录一个采样点。
func (r *Recorder) Record(s Sample) {
	r.Samples = append(r.Samples, s)
}

// Count 返回记录的采样数。
func (r *Recorder) Count() int { return len(r.Samples) }

// AvgQueueLen 平均队列长度。
func (r *Recorder) AvgQueueLen() float64 {
	if len(r.Samples) == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range r.Samples {
		sum += float64(s.QueueLen)
	}
	return sum / float64(len(r.Samples))
}

// MaxQueueLen 最大队列长度。
func (r *Recorder) MaxQueueLen() int {
	max := 0
	for _, s := range r.Samples {
		if s.QueueLen > max {
			max = s.QueueLen
		}
	}
	return max
}

// AvgUtilization 平均利用率。
func (r *Recorder) AvgUtilization(totalServers int) float64 {
	if len(r.Samples) == 0 || totalServers == 0 {
		return 0
	}
	sum := 0.0
	for _, s := range r.Samples {
		sum += float64(s.BusyCount)
	}
	return sum / float64(len(r.Samples)) / float64(totalServers)
}

// PeakUtilization 峰值利用率。
func (r *Recorder) PeakUtilization(totalServers int) float64 {
	if len(r.Samples) == 0 || totalServers == 0 {
		return 0
	}
	max := 0
	for _, s := range r.Samples {
		if s.BusyCount > max {
			max = s.BusyCount
		}
	}
	return float64(max) / float64(totalServers)
}

// QueueLenPercentile 队列长度第 p 百分位数。
func (r *Recorder) QueueLenPercentile(p float64) float64 {
	if len(r.Samples) == 0 {
		return 0
	}
	vals := make([]int, len(r.Samples))
	for i, s := range r.Samples {
		vals[i] = s.QueueLen
	}
	sort.Ints(vals)
	idx := p / 100.0 * float64(len(vals)-1)
	lo := int(math.Floor(idx))
	hi := int(math.Ceil(idx))
	if hi >= len(vals) {
		hi = len(vals) - 1
	}
	frac := idx - float64(lo)
	return float64(vals[lo])*(1-frac) + float64(vals[hi])*frac
}

// TimeSeries 返回时间-队列长度时间序列。
type TimeSeriesPoint struct {
	Time  float64 `json:"time"`
	Value float64 `json:"value"`
}

// QueueLenTimeSeries 返回队列长度随时间的变化。
func (r *Recorder) QueueLenTimeSeries() []TimeSeriesPoint {
	out := make([]TimeSeriesPoint, len(r.Samples))
	for i, s := range r.Samples {
		out[i] = TimeSeriesPoint{Time: s.Time, Value: float64(s.QueueLen)}
	}
	return out
}

// UtilizationTimeSeries 返回利用率随时间的变化。
func (r *Recorder) UtilizationTimeSeries(totalServers int) []TimeSeriesPoint {
	out := make([]TimeSeriesPoint, len(r.Samples))
	for i, s := range r.Samples {
		util := 0.0
		if totalServers > 0 {
			util = float64(s.BusyCount) / float64(totalServers)
		}
		out[i] = TimeSeriesPoint{Time: s.Time, Value: util}
	}
	return out
}
