// Package stats 提供 M/M/1 理论公式计算与仿真结果对照。
package stats

import (
	"errors"
	"fmt"

	"des-queue/internal/sim"
)

// Theory 是 M/M/1 稳态理论值。
type Theory struct {
	Rho float64 // 利用率 lambda/mu
	L   float64 // 平均系统人数 rho/(1-rho)
	W   float64 // 平均逗留时间 1/(mu-lambda)
	Wq  float64 // 平均等待时间 lambda/(mu*(mu-lambda))
}

// Diff 是仿真值与理论值的相对误差（百分比）。
type Diff struct {
	RhoErr  float64 // rho 相对误差 %
	WqErrPct float64 // Wq 相对误差 %
}

// TheoryMM1 计算 M/M/1 稳态理论值；参数非法或不稳态返回 error。
func TheoryMM1(lambda, mu float64) (Theory, error) {
	if lambda <= 0 || mu <= 0 {
		return Theory{}, errors.New("theory: lambda and mu must be > 0")
	}
	if lambda >= mu {
		return Theory{}, errors.New("theory: unstable, requires lambda < mu")
	}
	rho := lambda / mu
	return Theory{
		Rho: rho,
		L:   rho / (1 - rho),
		W:   1 / (mu - lambda),
		Wq:  lambda / (mu * (mu - lambda)),
	}, nil
}

// Compare 计算仿真统计 s 相对理论值 t 的相对误差百分比。
func Compare(s sim.Stats, t Theory) Diff {
	return Diff{
		RhoErr:   pct(s.Rho, t.Rho),
		WqErrPct: pct(s.Wq, t.Wq),
	}
}

// pct 计算 a 相对 b 的相对误差百分比。
func pct(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return (a - b) / b * 100
}

// Report 生成多行文本报告：配置、仿真结果、理论值与误差对照。
func Report(c sim.Config, s sim.Stats, d Diff) string {
	return fmt.Sprintf(`M/M/c 排队系统离散事件仿真报告
================================
参数
  lambda    = %.4f
  mu        = %.4f
  servers   = %d
  customers = %d
  seed      = %d

仿真结果
  rho        = %.4f
  L          = %.4f
  W          = %.4f
  Wq         = %.4f
  throughput = %.4f
  served     = %d

理论对照（M/M/1）
  rho 相对误差 = %+.2f%%
  Wq  相对误差 = %+.2f%%
`,
		c.Lambda, c.Mu, c.Servers, c.Customers, c.Seed,
		s.Rho, s.L, s.W, s.Wq, s.Throughput, s.Served,
		d.RhoErr, d.WqErrPct)
}
