// Package network 实现排队网络仿真：串联（tandem）和并联（parallel）排队节点。
// 串联网络中顾客依次经过多个服务站；并联网络中顾客被路由到负载最低的站。
package network

import (
	"errors"
	"math"
	"math/rand"
)

// ErrNoNodes 网络为空。
var ErrNoNodes = errors.New("network: no nodes defined")

// Node 排队网络中的一个服务节点。
type Node struct {
	Name    string
	Servers int     // 服务台数
	Mu      float64 // 单台服务率
}

// TandemConfig 串联排队网络配置。
type TandemConfig struct {
	Nodes     []Node
	Lambda    float64 // 外部到达率
	Customers int
	Seed      int64
}

// TandemResult 串联网络仿真结果。
type TandemResult struct {
	Nodes       []NodeStats
	TotalW      float64 // 整体平均逗留时间
	Throughput  float64
}

// NodeStats 单个节点的统计。
type NodeStats struct {
	Name    string  `json:"name"`
	Rho     float64 `json:"rho"`
	AvgWait float64 `json:"avg_wait"`
	AvgSvc  float64 `json:"avg_svc"`
}

// RunTandem 运行串联排队网络仿真（简化）。
func RunTandem(cfg TandemConfig) (*TandemResult, error) {
	if len(cfg.Nodes) == 0 {
		return nil, ErrNoNodes
	}
	if cfg.Lambda <= 0 || cfg.Customers < 1 {
		return nil, errors.New("network: invalid config")
	}
	r := rand.New(rand.NewSource(cfg.Seed))
	// 简化模型：每个顾客依次通过每个节点。
	// 在每个节点记录等待和服务时间。
	type record struct {
		arriveNode float64
		startSvc   float64
		leavNode   float64
	}
	nNodes := len(cfg.Nodes)
	nodeRecords := make([][]record, nNodes)
	for i := range nodeRecords {
		nodeRecords[i] = make([]record, cfg.Customers)
	}
	// 生成到达时间。
	arrivals := make([]float64, cfg.Customers)
	t := 0.0
	for i := range arrivals {
		u := r.Float64()
		for u == 0 {
			u = r.Float64()
		}
		t += -math.Log(u) / cfg.Lambda
		arrivals[i] = t
	}
	// 遍历每个节点。
	for ni, node := range cfg.Nodes {
		departure := make([]float64, node.Servers)
		for ci := 0; ci < cfg.Customers; ci++ {
			var arriveAt float64
			if ni == 0 {
				arriveAt = arrivals[ci]
			} else {
				arriveAt = nodeRecords[ni-1][ci].leavNode
			}
			// 找最早空闲的服务台。
			earliest := 0
			for s := 1; s < node.Servers; s++ {
				if departure[s] < departure[earliest] {
					earliest = s
				}
			}
			startSvc := math.Max(arriveAt, departure[earliest])
			u := r.Float64()
			for u == 0 {
				u = r.Float64()
			}
			svcTime := -math.Log(u) / node.Mu
			endSvc := startSvc + svcTime
			departure[earliest] = endSvc
			nodeRecords[ni][ci] = record{arriveNode: arriveAt, startSvc: startSvc, leavNode: endSvc}
		}
	}
	// 汇总。
	result := &TandemResult{}
	totalW := 0.0
	for ni, node := range cfg.Nodes {
		var sumWait, sumSvc float64
		var busyTime float64
		lastDepart := 0.0
		for ci := 0; ci < cfg.Customers; ci++ {
			rec := nodeRecords[ni][ci]
			sumWait += rec.startSvc - rec.arriveNode
			sumSvc += rec.leavNode - rec.startSvc
			if rec.leavNode > lastDepart {
				lastDepart = rec.leavNode
			}
		}
		busyTime = sumSvc
		ns := NodeStats{
			Name:    node.Name,
			Rho:     busyTime / (float64(node.Servers) * lastDepart),
			AvgWait: sumWait / float64(cfg.Customers),
			AvgSvc:  sumSvc / float64(cfg.Customers),
		}
		result.Nodes = append(result.Nodes, ns)
		totalW += ns.AvgWait + ns.AvgSvc
	}
	result.TotalW = totalW
	lastNode := nodeRecords[nNodes-1]
	lastDepart := lastNode[cfg.Customers-1].leavNode
	if lastDepart > 0 {
		result.Throughput = float64(cfg.Customers) / lastDepart
	}
	return result, nil
}

// ParallelConfig 并联网络配置。
type ParallelConfig struct {
	Nodes     []Node
	Lambda    float64
	Customers int
	Seed      int64
}

// ParallelResult 并联网络结果。
type ParallelResult struct {
	Nodes      []NodeStats
	AvgW       float64
	Throughput float64
}

// RunParallel 运行并联排队网络（顾客路由到最短队列）。
func RunParallel(cfg ParallelConfig) (*ParallelResult, error) {
	if len(cfg.Nodes) == 0 {
		return nil, ErrNoNodes
	}
	if cfg.Lambda <= 0 || cfg.Customers < 1 {
		return nil, errors.New("network: invalid config")
	}
	r := rand.New(rand.NewSource(cfg.Seed))
	nNodes := len(cfg.Nodes)
	// 每个节点的下次空闲时间。
	nodeAvail := make([]float64, nNodes)
	nodeCount := make([]int, nNodes)
	nodeBusy := make([]float64, nNodes)

	totalWait := 0.0
	t := 0.0
	for ci := 0; ci < cfg.Customers; ci++ {
		u := r.Float64()
		for u == 0 {
			u = r.Float64()
		}
		t += -math.Log(u) / cfg.Lambda
		// Route to node with earliest availability.
		best := 0
		for i := 1; i < nNodes; i++ {
			if nodeAvail[i] < nodeAvail[best] {
				best = i
			}
		}
		startSvc := math.Max(t, nodeAvail[best])
		totalWait += startSvc - t
		u2 := r.Float64()
		for u2 == 0 {
			u2 = r.Float64()
		}
		svcTime := -math.Log(u2) / cfg.Nodes[best].Mu
		nodeAvail[best] = startSvc + svcTime
		nodeCount[best]++
		nodeBusy[best] += svcTime
	}

	result := &ParallelResult{}
	lastDepart := 0.0
	for i := range cfg.Nodes {
		if nodeAvail[i] > lastDepart {
			lastDepart = nodeAvail[i]
		}
	}
	for i, node := range cfg.Nodes {
		rho := 0.0
		if lastDepart > 0 {
			rho = nodeBusy[i] / lastDepart
		}
		result.Nodes = append(result.Nodes, NodeStats{
			Name: node.Name,
			Rho:  rho,
		})
	}
	result.AvgW = totalWait / float64(cfg.Customers)
	if lastDepart > 0 {
		result.Throughput = float64(cfg.Customers) / lastDepart
	}
	return result, nil
}
