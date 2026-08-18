// topology.go 提供排队网络拓扑分析：瓶颈检测、容量规划建议。
package network

import (
	"math"
	"sort"
)

// BottleneckResult 瓶颈分析结果。
type BottleneckResult struct {
	NodeName   string  `json:"node_name"`
	Rho        float64 `json:"rho"`
	IsBottleneck bool  `json:"is_bottleneck"`
}

// FindBottleneck 找到利用率最高的节点（瓶颈）。
func FindBottleneck(stats []NodeStats) *BottleneckResult {
	if len(stats) == 0 {
		return nil
	}
	best := &BottleneckResult{NodeName: stats[0].Name, Rho: stats[0].Rho}
	for _, s := range stats[1:] {
		if s.Rho > best.Rho {
			best.NodeName = s.Name
			best.Rho = s.Rho
		}
	}
	best.IsBottleneck = best.Rho > 0.8
	return best
}

// CapacityRecommendation 容量规划建议。
type CapacityRecommendation struct {
	NodeName       string `json:"node_name"`
	CurrentRho     float64 `json:"current_rho"`
	TargetRho      float64 `json:"target_rho"`
	RequiredServers int    `json:"required_servers"`
}

// RecommendCapacity 根据目标利用率建议每个节点需要的最少服务台数。
func RecommendCapacity(nodes []Node, lambda float64, targetRho float64) []CapacityRecommendation {
	if targetRho <= 0 || targetRho >= 1 {
		targetRho = 0.7
	}
	var recs []CapacityRecommendation
	for _, n := range nodes {
		currentRho := lambda / (float64(n.Servers) * n.Mu)
		required := int(math.Ceil(lambda / (n.Mu * targetRho)))
		if required < 1 {
			required = 1
		}
		recs = append(recs, CapacityRecommendation{
			NodeName:       n.Name,
			CurrentRho:     currentRho,
			TargetRho:      targetRho,
			RequiredServers: required,
		})
	}
	return recs
}

// SortNodesByRho 按利用率降序排列节点统计。
func SortNodesByRho(stats []NodeStats) []NodeStats {
	sorted := make([]NodeStats, len(stats))
	copy(sorted, stats)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Rho > sorted[j].Rho })
	return sorted
}

// TotalCapacity 计算网络总服务能力（所有节点服务率之和）。
func TotalCapacity(nodes []Node) float64 {
	total := 0.0
	for _, n := range nodes {
		total += float64(n.Servers) * n.Mu
	}
	return total
}

// NetworkLoad 计算网络负载因子：lambda / 总容量。
func NetworkLoad(nodes []Node, lambda float64) float64 {
	cap := TotalCapacity(nodes)
	if cap == 0 {
		return math.Inf(1)
	}
	return lambda / cap
}
