// Package arrival 提供多种到达模式生成器：泊松过程、批量到达、时变到达率。
package arrival

import (
	"math"
	"math/rand"
)

// Generator 到达时间生成器接口。
type Generator interface {
	// Next 返回下一个到达间隔时间。
	Next() float64
	// Name 返回生成器名称。
	Name() string
}

// Poisson 恒定速率泊松过程。
type Poisson struct {
	r    *rand.Rand
	rate float64
}

// NewPoisson 创建泊松到达生成器。
func NewPoisson(seed int64, rate float64) *Poisson {
	return &Poisson{r: rand.New(rand.NewSource(seed)), rate: rate}
}

// Next 返回下一个指数间隔。
func (p *Poisson) Next() float64 {
	u := p.r.Float64()
	for u == 0 {
		u = p.r.Float64()
	}
	return -math.Log(u) / p.rate
}

// Name 实现 Generator。
func (p *Poisson) Name() string { return "poisson" }

// Batch 批量到达：每次到达一批，组间间隔为指数分布。
type Batch struct {
	r         *rand.Rand
	rate      float64
	batchSize int
	remaining int
}

// NewBatch 创建批量到达生成器。
func NewBatch(seed int64, rate float64, batchSize int) *Batch {
	if batchSize < 1 {
		batchSize = 1
	}
	return &Batch{r: rand.New(rand.NewSource(seed)), rate: rate, batchSize: batchSize}
}

// Next 同批内间隔为 0，不同批间隔为指数分布。
func (b *Batch) Next() float64 {
	if b.remaining > 0 {
		b.remaining--
		return 0
	}
	b.remaining = b.batchSize - 1
	u := b.r.Float64()
	for u == 0 {
		u = b.r.Float64()
	}
	return -math.Log(u) / b.rate
}

// Name 实现 Generator。
func (b *Batch) Name() string { return "batch" }

// TimeVarying 时变到达率：速率随时间按正弦波变化。
type TimeVarying struct {
	r         *rand.Rand
	baseRate  float64
	amplitude float64
	period    float64
	clock     float64
}

// NewTimeVarying 创建时变到达生成器。baseRate 为平均速率，amplitude 为波动幅度。
func NewTimeVarying(seed int64, baseRate, amplitude, period float64) *TimeVarying {
	return &TimeVarying{
		r: rand.New(rand.NewSource(seed)), baseRate: baseRate,
		amplitude: amplitude, period: period,
	}
}

// Next 使用 thinning 算法生成非齐次泊松过程的到达间隔。
func (tv *TimeVarying) Next() float64 {
	maxRate := tv.baseRate + tv.amplitude
	for {
		u := tv.r.Float64()
		for u == 0 {
			u = tv.r.Float64()
		}
		interval := -math.Log(u) / maxRate
		tv.clock += interval
		currentRate := tv.baseRate + tv.amplitude*math.Sin(2*math.Pi*tv.clock/tv.period)
		if currentRate < 0 {
			currentRate = 0
		}
		acceptance := currentRate / maxRate
		if tv.r.Float64() < acceptance {
			return interval
		}
	}
}

// Name 实现 Generator。
func (tv *TimeVarying) Name() string { return "time-varying" }

// Deterministic 固定间隔到达（D/M/1 模型）。
type Deterministic struct {
	interval float64
}

// NewDeterministic 创建固定间隔生成器。
func NewDeterministic(interval float64) *Deterministic {
	return &Deterministic{interval: interval}
}

// Next 返回固定间隔。
func (d *Deterministic) Next() float64 { return d.interval }

// Name 实现 Generator。
func (d *Deterministic) Name() string { return "deterministic" }

// Erlang k 阶 Erlang 分布到达间隔。
type Erlang struct {
	r    *rand.Rand
	k    int
	rate float64
}

// NewErlang 创建 Erlang-k 到达生成器。
func NewErlang(seed int64, k int, rate float64) *Erlang {
	if k < 1 {
		k = 1
	}
	return &Erlang{r: rand.New(rand.NewSource(seed)), k: k, rate: rate}
}

// Next 返回 k 个指数随机变量之和。
func (e *Erlang) Next() float64 {
	sum := 0.0
	for i := 0; i < e.k; i++ {
		u := e.r.Float64()
		for u == 0 {
			u = e.r.Float64()
		}
		sum += -math.Log(u) / e.rate
	}
	return sum
}

// Name 实现 Generator。
func (e *Erlang) Name() string { return "erlang" }
