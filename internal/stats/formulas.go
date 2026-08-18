// formulas.go 扩展 stats 包，提供 M/M/c 和 M/M/c/K 理论公式。
package stats

import "math"

// TheoryMMc 计算 M/M/c 稳态理论值（Erlang-C 公式）。
func TheoryMMc(lambda, mu float64, c int) (*MMcTheory, error) {
	if lambda <= 0 || mu <= 0 || c < 1 {
		return nil, ErrInvalid
	}
	rho := lambda / (float64(c) * mu)
	if rho >= 1 {
		return nil, ErrUnstable
	}
	erlangC := erlangCProb(lambda, mu, c)
	wq := erlangC / (float64(c) * mu * (1 - rho))
	w := wq + 1/mu
	lq := lambda * wq
	l := lambda * w
	return &MMcTheory{Rho: rho, Lq: lq, L: l, Wq: wq, W: w, ErlangC: erlangC}, nil
}

// MMcTheory M/M/c 理论值。
type MMcTheory struct {
	Rho     float64
	Lq      float64 // 平均排队人数
	L       float64 // 平均系统人数
	Wq      float64 // 平均排队时间
	W       float64 // 平均逗留时间
	ErlangC float64 // Erlang-C 概率
}

// erlangCProb 计算 Erlang-C 概率（等待概率）。
func erlangCProb(lambda, mu float64, c int) float64 {
	a := lambda / mu // 提供的话务量
	rho := a / float64(c)
	// 分子：(a^c / c!) * (1 / (1-rho))
	numerator := math.Pow(a, float64(c)) / factorial(c) / (1 - rho)
	// 分母：sum_{k=0}^{c-1} a^k/k! + numerator
	sum := 0.0
	for k := 0; k < c; k++ {
		sum += math.Pow(a, float64(k)) / factorial(k)
	}
	return numerator / (sum + numerator)
}

// factorial 计算阶乘（整数）。
func factorial(n int) float64 {
	if n <= 1 {
		return 1
	}
	f := 1.0
	for i := 2; i <= n; i++ {
		f *= float64(i)
	}
	return f
}

// ErrInvalid 参数不合法。
var ErrInvalid = errNew("stats: invalid parameters")

// ErrUnstable 系统不稳定。
var ErrUnstable = errNew("stats: system is unstable (rho >= 1)")

func errNew(msg string) error { return &staticError{msg} }

type staticError struct{ msg string }

func (e *staticError) Error() string { return e.msg }

// ErlangBLoss 计算 Erlang-B 丢失概率（M/M/c/c 系统）。
func ErlangBLoss(a float64, c int) float64 {
	if c < 1 {
		return 1
	}
	// 递推公式：B(a, 0) = 1; B(a, k) = (a * B(a, k-1)) / (k + a * B(a, k-1))
	b := 1.0
	for k := 1; k <= c; k++ {
		b = (a * b) / (float64(k) + a*b)
	}
	return b
}

// MM1Theory 包装原有 TheoryMM1 为新 API 格式。
func MM1Theory(lambda, mu float64) (*MMcTheory, error) {
	t, err := TheoryMM1(lambda, mu)
	if err != nil {
		return nil, err
	}
	return &MMcTheory{Rho: t.Rho, Lq: t.Rho * t.Rho / (1 - t.Rho), L: t.L, Wq: t.Wq, W: t.W, ErlangC: t.Rho}, nil
}

// LittlesLaw 验证 Little 定律：L = λW。
func LittlesLaw(lambda, l, w float64) float64 {
	return math.Abs(l - lambda*w)
}
