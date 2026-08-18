// Package sim 实现 M/M/1 与 M/M/c 排队系统的离散事件仿真。
package sim

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"des-queue/internal/event"
)

// Config 是一次仿真的输入参数。
type Config struct {
	Lambda    float64 // 到达率
	Mu        float64 // 单服务台服务率
	Servers   int     // 服务台数
	Customers int     // 仿真处理的到达顾客总数
	Seed      int64   // 随机种子（保证结果可复现）
}

// Stats 是仿真输出的统计量。
type Stats struct {
	Rho        float64 // 利用率（M/M/c 为平均忙服务台比例）
	L          float64 // 平均系统人数（含服务中）
	W          float64 // 平均逗留时间
	Wq         float64 // 平均排队等待时间
	Throughput float64 // 吞吐量（单位时间完成服务数）
	Served     int     // 完成服务的顾客数
}

// ExpRand 是生成指数分布间隔的确定性序列器。
type ExpRand struct {
	r    *rand.Rand
	rate float64
}

// NewExpRand 以固定 seed 构造速率为 rate 的指数间隔序列器。
func NewExpRand(seed int64, rate float64) *ExpRand {
	return &ExpRand{r: rand.New(rand.NewSource(seed)), rate: rate}
}

// Next 返回下一个指数分布间隔。
func (e *ExpRand) Next() float64 {
	u := e.r.Float64()
	return -math.Log(1 - u) / e.rate
}

// Validate 检查配置合法性：lambda/mu 必须为正，servers/customers 至少为 1；
// M/M/1（servers==1）还要求 lambda < mu，否则系统不稳态。
func Validate(c Config) error {
	if c.Lambda <= 0 {
		return fmt.Errorf("lambda must be > 0, got %g", c.Lambda)
	}
	if c.Mu <= 0 {
		return fmt.Errorf("mu must be > 0, got %g", c.Mu)
	}
	if c.Servers < 1 {
		return fmt.Errorf("servers must be >= 1, got %d", c.Servers)
	}
	if c.Customers < 1 {
		return fmt.Errorf("customers must be >= 1, got %d", c.Customers)
	}
	if c.Servers == 1 && c.Lambda >= c.Mu {
		return errors.New("unstable: M/M/1 requires lambda < mu")
	}
	return nil
}

// RunMM1 运行 M/M/1 仿真，返回统计量。
func RunMM1(c Config) (Stats, error) {
	return run(c, 1)
}

// RunMMc 运行 M/M/c 仿真（servers 个服务台），返回统计量。
func RunMMc(c Config, servers int) (Stats, error) {
	if servers < 1 {
		return Stats{}, fmt.Errorf("servers must be >= 1, got %d", servers)
	}
	return run(c, servers)
}

// run 是共享的事件驱动仿真引擎：
// 到达时若有空闲服务台则立即服务并调度 Departure，否则排队；
// 服务完成时释放服务台并从队首取下一位顾客。
// FEL 中途为空则提前正常结束。
func run(c Config, servers int) (Stats, error) {
	if err := Validate(c); err != nil {
		return Stats{}, err
	}

	interArrival := NewExpRand(c.Seed, c.Lambda)
	service := NewExpRand(c.Seed+1, c.Mu)

	fel := &event.Queue{}
	fel.Push(event.Event{Time: interArrival.Next(), Kind: event.Arrival, Customer: 0})

	arrived := 0
	served := 0
	idle := servers
	n := 0         // 当前系统内人数（排队 + 服务中）
	var line []int // 等待队列中的顾客编号
	arrivalAt := make(map[int]float64)

	var now, lastT float64
	var busyArea, sysArea float64 // 忙服务台数 / 系统人数的时间积分
	var waitSum, systemSum float64

	for {
		e, ok := fel.Pop()
		if !ok {
			// FEL 提前为空，正常结束。
			break
		}
		now = e.Time
		dt := now - lastT
		busyArea += float64(servers-idle) * dt
		sysArea += float64(n) * dt
		lastT = now

		switch e.Kind {
		case event.Arrival:
			arrivalAt[e.Customer] = now
			n++
			if idle > 0 {
				idle--
				fel.Push(event.Event{Time: now + service.Next(), Kind: event.Departure, Customer: e.Customer})
			} else {
				line = append(line, e.Customer)
			}
			arrived++
			if arrived < c.Customers {
				fel.Push(event.Event{Time: now + interArrival.Next(), Kind: event.Arrival, Customer: arrived})
			}

		case event.Departure:
			served++
			idle++
			n--
			systemSum += now - arrivalAt[e.Customer]
			delete(arrivalAt, e.Customer)
			if len(line) > 0 {
				next := line[0]
				line = line[1:]
				idle--
				waitSum += now - arrivalAt[next]
				fel.Push(event.Event{Time: now + service.Next(), Kind: event.Departure, Customer: next})
			}
		}

		if served >= c.Customers {
			break
		}
	}

	if served == 0 || now <= 0 {
		return Stats{}, errors.New("no customers were served")
	}

	fs := float64(served)
	return Stats{
		Rho:        busyArea / (float64(servers) * now),
		L:          sysArea / now,
		W:          systemSum / fs,
		Wq:         waitSum / fs,
		Throughput: fs / now,
		Served:     served,
	}, nil
}
