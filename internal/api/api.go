// Package api 提供 des-queue 的 HTTP API 和前端文件服务。
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"des-queue/internal/sim"
	"des-queue/internal/stats"
)

// Server HTTP API 服务器。
type Server struct {
	mux    *http.ServeMux
	addr   string
	webDir string
}

// Config 服务器配置。
type Config struct {
	Addr   string
	WebDir string
}

// New 创建 API 服务器。
func New(cfg Config) *Server {
	s := &Server{mux: http.NewServeMux(), addr: cfg.Addr, webDir: cfg.WebDir}
	s.routes()
	return s
}

// Handler 返回 HTTP handler。
func (s *Server) Handler() http.Handler { return s.mux }

// ListenAndServe 启动服务器。
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/simulate", s.handleSimulate)
	s.mux.HandleFunc("/api/theory", s.handleTheory)
	s.mux.HandleFunc("/api/health", s.handleHealth)
	if s.webDir != "" {
		s.mux.Handle("/", http.FileServer(http.Dir(s.webDir)))
	}
}

// SimRequest 仿真请求。
type SimRequest struct {
	Lambda    float64 `json:"lambda"`
	Mu        float64 `json:"mu"`
	Servers   int     `json:"servers"`
	Customers int     `json:"customers"`
	Seed      int64   `json:"seed"`
}

// SimResponse 仿真响应。
type SimResponse struct {
	Rho        float64 `json:"rho"`
	L          float64 `json:"l"`
	W          float64 `json:"w"`
	Wq         float64 `json:"wq"`
	Throughput float64 `json:"throughput"`
	Served     int     `json:"served"`
	RhoErr     float64 `json:"rho_err_pct"`
	WqErr      float64 `json:"wq_err_pct"`
}

// TheoryRequest 理论值请求。
type TheoryRequest struct {
	Lambda float64 `json:"lambda"`
	Mu     float64 `json:"mu"`
}

// TheoryResponse 理论值响应。
type TheoryResponse struct {
	Rho float64 `json:"rho"`
	L   float64 `json:"l"`
	W   float64 `json:"w"`
	Wq  float64 `json:"wq"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpErr(w, 405, "POST only")
		return
	}
	var req SimRequest
	if err := readBody(r, &req); err != nil {
		httpErr(w, 400, err.Error())
		return
	}
	if req.Customers == 0 {
		req.Customers = 8000
	}
	if req.Seed == 0 {
		req.Seed = 42
	}
	if req.Servers == 0 {
		req.Servers = 1
	}
	cfg := sim.Config{Lambda: req.Lambda, Mu: req.Mu, Servers: req.Servers, Customers: req.Customers, Seed: req.Seed}
	var st sim.Stats
	var err error
	if req.Servers == 1 {
		st, err = sim.RunMM1(cfg)
	} else {
		st, err = sim.RunMMc(cfg, req.Servers)
	}
	if err != nil {
		httpErr(w, 422, err.Error())
		return
	}
	resp := SimResponse{Rho: st.Rho, L: st.L, W: st.W, Wq: st.Wq, Throughput: st.Throughput, Served: st.Served}
	if req.Servers == 1 {
		th, terr := stats.TheoryMM1(req.Lambda, req.Mu)
		if terr == nil {
			d := stats.Compare(st, th)
			resp.RhoErr = d.RhoErr
			resp.WqErr = d.WqErrPct
		}
	}
	writeBody(w, resp)
}

func (s *Server) handleTheory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpErr(w, 405, "POST only")
		return
	}
	var req TheoryRequest
	if err := readBody(r, &req); err != nil {
		httpErr(w, 400, err.Error())
		return
	}
	th, err := stats.TheoryMM1(req.Lambda, req.Mu)
	if err != nil {
		httpErr(w, 422, err.Error())
		return
	}
	writeBody(w, TheoryResponse{Rho: th.Rho, L: th.L, W: th.W, Wq: th.Wq})
}

func readBody(r *http.Request, v interface{}) error {
	b, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return json.Unmarshal(b, v)
}

func writeBody(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func httpErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
