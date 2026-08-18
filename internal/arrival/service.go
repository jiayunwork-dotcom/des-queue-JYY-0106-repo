// service.go 扩展 arrival 包，提供服务时间分布生成器。
package arrival

import (
	"math"
	"math/rand"
)

// ServiceGenerator 服务时间生成器接口。
type ServiceGenerator interface {
	Next() float64
	Name() string
}

// ExpService 指数分布服务时间。
type ExpService struct {
	r    *rand.Rand
	rate float64
}

// NewExpService 创建指数服务时间生成器。
func NewExpService(seed int64, mu float64) *ExpService {
	return &ExpService{r: rand.New(rand.NewSource(seed)), rate: mu}
}

// Next 返回下一个服务时间。
func (e *ExpService) Next() float64 {
	u := e.r.Float64()
	for u == 0 {
		u = e.r.Float64()
	}
	return -math.Log(u) / e.rate
}

// Name 实现接口。
func (e *ExpService) Name() string { return "exponential" }

// UniformService 均匀分布服务时间。
type UniformService struct {
	r   *rand.Rand
	min float64
	max float64
}

// NewUniformService 创建均匀分布服务时间生成器。
func NewUniformService(seed int64, min, max float64) *UniformService {
	return &UniformService{r: rand.New(rand.NewSource(seed)), min: min, max: max}
}

// Next 返回 [min, max] 内均匀分布的服务时间。
func (u *UniformService) Next() float64 {
	return u.min + u.r.Float64()*(u.max-u.min)
}

// Name 实现接口。
func (u *UniformService) Name() string { return "uniform" }

// LogNormalService 对数正态分布服务时间。
type LogNormalService struct {
	r     *rand.Rand
	mu    float64 // 对数均值
	sigma float64 // 对数标准差
}

// NewLogNormalService 创建对数正态服务时间生成器。
func NewLogNormalService(seed int64, mu, sigma float64) *LogNormalService {
	return &LogNormalService{r: rand.New(rand.NewSource(seed)), mu: mu, sigma: sigma}
}

// Next 使用 Box-Muller 生成正态样本后取指数。
func (l *LogNormalService) Next() float64 {
	u1 := l.r.Float64()
	for u1 == 0 {
		u1 = l.r.Float64()
	}
	u2 := l.r.Float64()
	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return math.Exp(l.mu + l.sigma*z)
}

// Name 实现接口。
func (l *LogNormalService) Name() string { return "lognormal" }

// ConstantService 固定服务时间（D/M/1 的服务端）。
type ConstantService struct {
	duration float64
}

// NewConstantService 创建固定时间服务生成器。
func NewConstantService(duration float64) *ConstantService {
	return &ConstantService{duration: duration}
}

// Next 返回固定时间。
func (c *ConstantService) Next() float64 { return c.duration }

// Name 实现接口。
func (c *ConstantService) Name() string { return "constant" }

// WeibullService Weibull 分布服务时间。
type WeibullService struct {
	r     *rand.Rand
	shape float64
	scale float64
}

// NewWeibullService 创建 Weibull 分布服务时间生成器。
func NewWeibullService(seed int64, shape, scale float64) *WeibullService {
	return &WeibullService{r: rand.New(rand.NewSource(seed)), shape: shape, scale: scale}
}

// Next 返回 Weibull 分布样本。
func (w *WeibullService) Next() float64 {
	u := w.r.Float64()
	for u == 0 {
		u = w.r.Float64()
	}
	return w.scale * math.Pow(-math.Log(u), 1.0/w.shape)
}

// Name 实现接口。
func (w *WeibullService) Name() string { return "weibull" }
