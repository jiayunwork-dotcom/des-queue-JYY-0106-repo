// advanced.go 提供高级仿真模式：带有限缓冲的 M/M/c/K 系统、
// 带不耐烦顾客的 M/M/c/K+M 系统、以及多优先级服务。
package sim

import (
	"errors"
	"math"
	"math/rand"

	"des-queue/internal/event"
)

// MMcKConfig M/M/c/K 有限缓冲排队系统配置。
type MMcKConfig struct {
	Lambda    float64
	Mu        float64
	Servers   int
	Capacity  int // 系统最大容量（含服务中），超过则拒绝
	Customers int
	Seed      int64
}

// MMcKStats 结果包含拒绝率。
type MMcKStats struct {
	Stats
	Rejected int     // 被拒绝的顾客数
	RejectRate float64 // 拒绝率
}

// RunMMcK 运行 M/M/c/K 仿真。
func RunMMcK(c MMcKConfig) (MMcKStats, error) {
	if c.Lambda <= 0 || c.Mu <= 0 || c.Servers < 1 || c.Customers < 1 {
		return MMcKStats{}, errors.New("sim: invalid MMcK config")
	}
	if c.Capacity < c.Servers {
		return MMcKStats{}, errors.New("sim: capacity must be >= servers")
	}
	r := rand.New(rand.NewSource(c.Seed))
	svcR := rand.New(rand.NewSource(c.Seed + 1))

	fel := &event.Queue{}
	fel.Push(event.Event{Time: expRand(r, c.Lambda), Kind: event.Arrival, Customer: 0})

	arrived, served, rejected := 0, 0, 0
	idle := c.Servers
	n := 0 // 当前系统内人数
	var line []int
	arrivalAt := make(map[int]float64)
	var now, lastT, busyArea, sysArea, waitSum, systemSum float64

	for {
		e, ok := fel.Pop()
		if !ok {
			break
		}
		now = e.Time
		dt := now - lastT
		busyArea += float64(c.Servers-idle) * dt
		sysArea += float64(n) * dt
		lastT = now

		switch e.Kind {
		case event.Arrival:
			arrived++
			if n >= c.Capacity {
				rejected++
			} else {
				arrivalAt[e.Customer] = now
				n++
				if idle > 0 {
					idle--
					fel.Push(event.Event{Time: now + expRand(svcR, c.Mu), Kind: event.Departure, Customer: e.Customer})
				} else {
					line = append(line, e.Customer)
				}
			}
			if arrived < c.Customers {
				fel.Push(event.Event{Time: now + expRand(r, c.Lambda), Kind: event.Arrival, Customer: arrived})
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
				fel.Push(event.Event{Time: now + expRand(svcR, c.Mu), Kind: event.Departure, Customer: next})
			}
		}
		if served+rejected >= c.Customers {
			break
		}
	}
	if served == 0 || now <= 0 {
		return MMcKStats{}, errors.New("sim: no customers served in MMcK")
	}
	fs := float64(served)
	return MMcKStats{
		Stats: Stats{
			Rho:        busyArea / (float64(c.Servers) * now),
			L:          sysArea / now,
			W:          systemSum / fs,
			Wq:         waitSum / fs,
			Throughput: fs / now,
			Served:     served,
		},
		Rejected:   rejected,
		RejectRate: float64(rejected) / float64(arrived),
	}, nil
}

// ImpatienceConfig 带不耐烦顾客的配置。
type ImpatienceConfig struct {
	Lambda     float64
	Mu         float64
	Servers    int
	Patience   float64 // 平均耐心时间（指数分布）
	Customers  int
	Seed       int64
}

// ImpatienceStats 含放弃率的结果。
type ImpatienceStats struct {
	Stats
	Abandoned    int
	AbandonRate  float64
}

// RunWithImpatience 运行带不耐烦顾客的仿真。
func RunWithImpatience(c ImpatienceConfig) (ImpatienceStats, error) {
	if c.Lambda <= 0 || c.Mu <= 0 || c.Servers < 1 || c.Customers < 1 || c.Patience <= 0 {
		return ImpatienceStats{}, errors.New("sim: invalid impatience config")
	}
	r := rand.New(rand.NewSource(c.Seed))
	svcR := rand.New(rand.NewSource(c.Seed + 1))
	patR := rand.New(rand.NewSource(c.Seed + 2))

	arrived, served, abandoned := 0, 0, 0
	idle := c.Servers
	n := 0
	var now, lastT, busyArea, waitSum, systemSum float64

	type waiter struct {
		id       int
		deadline float64
	}
	var line []waiter
	arrivalAt := make(map[int]float64)

	fel := &event.Queue{}
	fel.Push(event.Event{Time: expRand(r, c.Lambda), Kind: event.Arrival, Customer: 0})

	for {
		e, ok := fel.Pop()
		if !ok {
			break
		}
		now = e.Time
		dt := now - lastT
		busyArea += float64(c.Servers-idle) * dt
		lastT = now

		switch e.Kind {
		case event.Arrival:
			arrivalAt[e.Customer] = now
			n++
			if idle > 0 {
				idle--
				fel.Push(event.Event{Time: now + expRand(svcR, c.Mu), Kind: event.Departure, Customer: e.Customer})
			} else {
				deadline := now + expRand(patR, 1.0/c.Patience)
				line = append(line, waiter{id: e.Customer, deadline: deadline})
			}
			arrived++
			if arrived < c.Customers {
				fel.Push(event.Event{Time: now + expRand(r, c.Lambda), Kind: event.Arrival, Customer: arrived})
			}
		case event.Departure:
			served++
			idle++
			n--
			systemSum += now - arrivalAt[e.Customer]
			delete(arrivalAt, e.Customer)
			// 从队列中取下一个未放弃的顾客。
			for len(line) > 0 {
				w := line[0]
				line = line[1:]
				if now > w.deadline {
					abandoned++
					n--
					delete(arrivalAt, w.id)
					continue
				}
				idle--
				waitSum += now - arrivalAt[w.id]
				fel.Push(event.Event{Time: now + expRand(svcR, c.Mu), Kind: event.Departure, Customer: w.id})
				break
			}
		}
		if served+abandoned >= c.Customers {
			break
		}
	}
	if served == 0 {
		return ImpatienceStats{}, errors.New("sim: no customers served")
	}
	fs := float64(served)
	return ImpatienceStats{
		Stats: Stats{
			Rho:        busyArea / (float64(c.Servers) * now),
			W:          systemSum / fs,
			Wq:         waitSum / fs,
			Throughput: fs / now,
			Served:     served,
		},
		Abandoned:   abandoned,
		AbandonRate: float64(abandoned) / float64(arrived),
	}, nil
}

func expRand(r *rand.Rand, rate float64) float64 {
	u := r.Float64()
	for u == 0 {
		u = r.Float64()
	}
	return -math.Log(u) / rate
}
