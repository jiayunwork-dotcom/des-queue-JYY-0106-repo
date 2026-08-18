// discipline.go 提供更多排队纪律：LIFO（后进先出）、随机服务顺序。
package queue

import "math/rand"

// LIFO 后进先出队列（栈行为）。
type LIFO struct {
	buf []Customer
	cap int
}

// NewLIFO 创建容量为 cap 的 LIFO 队列。
func NewLIFO(cap int) *LIFO {
	return &LIFO{cap: cap}
}

// Enqueue 入队（push to top）。
func (l *LIFO) Enqueue(c Customer) error {
	if l.cap > 0 && len(l.buf) >= l.cap {
		return ErrFull
	}
	l.buf = append(l.buf, c)
	return nil
}

// Dequeue 出队（pop from top, LIFO）。
func (l *LIFO) Dequeue() (Customer, error) {
	if len(l.buf) == 0 {
		return Customer{}, ErrEmpty
	}
	n := len(l.buf) - 1
	c := l.buf[n]
	l.buf = l.buf[:n]
	return c, nil
}

// Len 当前排队人数。
func (l *LIFO) Len() int { return len(l.buf) }

// IsFull 是否已满。
func (l *LIFO) IsFull() bool { return l.cap > 0 && len(l.buf) >= l.cap }

// IsEmpty 是否为空。
func (l *LIFO) IsEmpty() bool { return len(l.buf) == 0 }

// RandomOrder 随机服务顺序队列。
type RandomOrder struct {
	buf  []Customer
	cap  int
	rng  *rand.Rand
}

// NewRandomOrder 创建随机服务顺序队列。
func NewRandomOrder(cap int, seed int64) *RandomOrder {
	return &RandomOrder{cap: cap, rng: rand.New(rand.NewSource(seed))}
}

// Enqueue 入队。
func (ro *RandomOrder) Enqueue(c Customer) error {
	if ro.cap > 0 && len(ro.buf) >= ro.cap {
		return ErrFull
	}
	ro.buf = append(ro.buf, c)
	return nil
}

// Dequeue 随机选择一个出队。
func (ro *RandomOrder) Dequeue() (Customer, error) {
	if len(ro.buf) == 0 {
		return Customer{}, ErrEmpty
	}
	idx := ro.rng.Intn(len(ro.buf))
	c := ro.buf[idx]
	ro.buf[idx] = ro.buf[len(ro.buf)-1]
	ro.buf = ro.buf[:len(ro.buf)-1]
	return c, nil
}

// Len 当前排队人数。
func (ro *RandomOrder) Len() int { return len(ro.buf) }

// IsFull 是否已满。
func (ro *RandomOrder) IsFull() bool { return ro.cap > 0 && len(ro.buf) >= ro.cap }

// IsEmpty 是否为空。
func (ro *RandomOrder) IsEmpty() bool { return len(ro.buf) == 0 }

// RoundRobin 轮询分配器：将顾客分配到 n 个子队列。
type RoundRobin struct {
	queues []Strategy
	next   int
}

// NewRoundRobin 创建包含 n 个 FIFO 子队列的轮询分配器。
func NewRoundRobin(n, perQueueCap int) *RoundRobin {
	queues := make([]Strategy, n)
	for i := range queues {
		queues[i] = NewFIFO(perQueueCap)
	}
	return &RoundRobin{queues: queues}
}

// Assign 将顾客分配到下一个队列。
func (rr *RoundRobin) Assign(c Customer) error {
	err := rr.queues[rr.next].Enqueue(c)
	rr.next = (rr.next + 1) % len(rr.queues)
	return err
}

// TotalLen 所有队列的总长度。
func (rr *RoundRobin) TotalLen() int {
	total := 0
	for _, q := range rr.queues {
		total += q.Len()
	}
	return total
}

// ShortestQueue 返回最短队列的索引。
func (rr *RoundRobin) ShortestQueue() int {
	min := 0
	for i := 1; i < len(rr.queues); i++ {
		if rr.queues[i].Len() < rr.queues[min].Len() {
			min = i
		}
	}
	return min
}
