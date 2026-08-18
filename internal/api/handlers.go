// handlers.go 扩展 api 包，增加 M/M/c/K 仿真和网络仿真端点。
package api

import (
	"net/http"

	"des-queue/internal/network"
	"des-queue/internal/sim"
)

// MMcKRequest M/M/c/K 仿真请求。
type MMcKRequest struct {
	Lambda    float64 `json:"lambda"`
	Mu        float64 `json:"mu"`
	Servers   int     `json:"servers"`
	Capacity  int     `json:"capacity"`
	Customers int     `json:"customers"`
	Seed      int64   `json:"seed"`
}

// MMcKResponse M/M/c/K 仿真响应。
type MMcKResponse struct {
	SimResponse
	Rejected   int     `json:"rejected"`
	RejectRate float64 `json:"reject_rate"`
}

// HandleMMcK 处理 M/M/c/K 请求（可被外部注册）。
func HandleMMcK(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpErr(w, 405, "POST only")
		return
	}
	var req MMcKRequest
	if err := readBody(r, &req); err != nil {
		httpErr(w, 400, err.Error())
		return
	}
	if req.Customers == 0 {
		req.Customers = 5000
	}
	if req.Seed == 0 {
		req.Seed = 42
	}
	cfg := sim.MMcKConfig{
		Lambda: req.Lambda, Mu: req.Mu, Servers: req.Servers,
		Capacity: req.Capacity, Customers: req.Customers, Seed: req.Seed,
	}
	st, err := sim.RunMMcK(cfg)
	if err != nil {
		httpErr(w, 422, err.Error())
		return
	}
	writeBody(w, MMcKResponse{
		SimResponse: SimResponse{Rho: st.Rho, L: st.L, W: st.W, Wq: st.Wq, Throughput: st.Throughput, Served: st.Served},
		Rejected:    st.Rejected,
		RejectRate:  st.RejectRate,
	})
}

// TandemRequest 串联网络仿真请求。
type TandemRequest struct {
	Lambda    float64        `json:"lambda"`
	Customers int            `json:"customers"`
	Seed      int64          `json:"seed"`
	Nodes     []NodeRequest  `json:"nodes"`
}

// NodeRequest 节点配置。
type NodeRequest struct {
	Name    string  `json:"name"`
	Servers int     `json:"servers"`
	Mu      float64 `json:"mu"`
}

// TandemResponse 串联网络响应。
type TandemResponse struct {
	Nodes      []network.NodeStats `json:"nodes"`
	TotalW     float64             `json:"total_w"`
	Throughput float64             `json:"throughput"`
}

// HandleTandem 处理串联网络仿真请求。
func HandleTandem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpErr(w, 405, "POST only")
		return
	}
	var req TandemRequest
	if err := readBody(r, &req); err != nil {
		httpErr(w, 400, err.Error())
		return
	}
	nodes := make([]network.Node, len(req.Nodes))
	for i, n := range req.Nodes {
		nodes[i] = network.Node{Name: n.Name, Servers: n.Servers, Mu: n.Mu}
	}
	cfg := network.TandemConfig{
		Nodes: nodes, Lambda: req.Lambda,
		Customers: req.Customers, Seed: req.Seed,
	}
	res, err := network.RunTandem(cfg)
	if err != nil {
		httpErr(w, 422, err.Error())
		return
	}
	writeBody(w, TandemResponse{Nodes: res.Nodes, TotalW: res.TotalW, Throughput: res.Throughput})
}
