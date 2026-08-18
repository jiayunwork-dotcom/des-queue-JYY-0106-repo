// des-queue：M/M/1 与 M/M/c 排队系统离散事件仿真 CLI。
package main

import (
	"flag"
	"fmt"
	"os"

	"des-queue/internal/api"
	"des-queue/internal/sim"
	"des-queue/internal/stats"
)

func main() {
	lambda := flag.Float64("lambda", 0.5, "到达率（每单位时间顾客数）")
	mu := flag.Float64("mu", 1.0, "服务率（每服务台）")
	servers := flag.Int("servers", 1, "服务台数量")
	customers := flag.Int("customers", 8000, "仿真顾客数")
	seed := flag.Int64("seed", 42, "随机种子")
	serveMode := flag.Bool("serve", false, "启动 HTTP API 服务器")
	serveAddr := flag.String("addr", ":8080", "HTTP 监听地址")
	webDir := flag.String("web", "web", "前端静态文件目录")
	flag.Parse()

	if *serveMode {
		srv := api.New(api.Config{Addr: *serveAddr, WebDir: *webDir})
		fmt.Fprintf(os.Stderr, "des-queue: serving on %s\n", *serveAddr)
		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	cfg := sim.Config{
		Lambda:    *lambda,
		Mu:        *mu,
		Servers:   *servers,
		Customers: *customers,
		Seed:      *seed,
	}

	var (
		s   sim.Stats
		err error
	)
	if *servers == 1 {
		s, err = sim.RunMM1(cfg)
	} else {
		s, err = sim.RunMMc(cfg, *servers)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var d stats.Diff
	if *servers == 1 {
		t, terr := stats.TheoryMM1(*lambda, *mu)
		if terr == nil {
			d = stats.Compare(s, t)
		}
	}

	fmt.Print(stats.Report(cfg, s, d))
}
