// balance.go 提供负载均衡策略和路由算法。
package network

// RoutingPolicy 路由策略枚举。
type RoutingPolicy int

const (
	RouteShortestQueue RoutingPolicy = iota // 最短队列
	RouteRoundRobin                         // 轮询
	RouteRandom                             // 随机
	RouteLeastLoaded                        // 最低负载
)

// String 返回策略名。
func (rp RoutingPolicy) String() string {
	switch rp {
	case RouteShortestQueue:
		return "shortest-queue"
	case RouteRoundRobin:
		return "round-robin"
	case RouteRandom:
		return "random"
	case RouteLeastLoaded:
		return "least-loaded"
	default:
		return "unknown"
	}
}

// LoadBalancer 负载均衡器。
type LoadBalancer struct {
	policy RoutingPolicy
	nodes  []Node
	next   int // round-robin 计数器
}

// NewLoadBalancer 创建负载均衡器。
func NewLoadBalancer(policy RoutingPolicy, nodes []Node) *LoadBalancer {
	return &LoadBalancer{policy: policy, nodes: nodes}
}

// Route 根据策略返回应该路由到的节点索引。
// loads 表示每个节点当前排队人数。
func (lb *LoadBalancer) Route(loads []int) int {
	switch lb.policy {
	case RouteShortestQueue:
		return shortestIdx(loads)
	case RouteRoundRobin:
		idx := lb.next
		lb.next = (lb.next + 1) % len(lb.nodes)
		return idx
	case RouteLeastLoaded:
		return leastLoadedIdx(loads, lb.nodes)
	default:
		return 0
	}
}

// shortestIdx 找队列最短的节点。
func shortestIdx(loads []int) int {
	min := 0
	for i := 1; i < len(loads); i++ {
		if loads[i] < loads[min] {
			min = i
		}
	}
	return min
}

// leastLoadedIdx 按 load/capacity 比选择。
func leastLoadedIdx(loads []int, nodes []Node) int {
	min := 0
	minRatio := float64(loads[0]) / float64(nodes[0].Servers)
	for i := 1; i < len(loads); i++ {
		ratio := float64(loads[i]) / float64(nodes[i].Servers)
		if ratio < minRatio {
			min = i
			minRatio = ratio
		}
	}
	return min
}

// BalanceReport 负载均衡报告。
type BalanceReport struct {
	Policy string         `json:"policy"`
	Distribution []float64 `json:"distribution"` // 每个节点分配的比例
}

// ComputeDistribution 根据路由次数计算分配比例。
func ComputeDistribution(routeCount []int) []float64 {
	total := 0
	for _, c := range routeCount {
		total += c
	}
	dist := make([]float64, len(routeCount))
	if total == 0 {
		return dist
	}
	for i, c := range routeCount {
		dist[i] = float64(c) / float64(total)
	}
	return dist
}
