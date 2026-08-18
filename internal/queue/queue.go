// Package queue 提供多种排队策略实现：FIFO、优先级队列、限长队列、SJF。
package queue

import (
	"container/heap"
	"errors"
)

// ErrFull 队列已满时返回。
var ErrFull = errors.New("queue: full")

// ErrEmpty 队列为空时返回。
var ErrEmpty = errors.New("queue: empty")

// Customer 表示一个等待中的顾客。
type Customer struct {
	ID        int
	Priority  int     // 越小越优先
	ArriveAt  float64
	ServiceDur float64 // 预期服务时长（SJF 用）
}

// Strategy 排队策略接口。
type Strategy interface {
	Enqueue(c Customer) error
	Dequeue() (Customer, error)
	Len() int
	IsFull() bool
	IsEmpty() bool
}

// FIFO 先进先出队列。
type FIFO struct {
	buf []Customer
	cap int
}

// NewFIFO 创建容量为 cap 的 FIFO 队列；cap<=0 表示无限。
func NewFIFO(cap int) *FIFO {
	return &FIFO{cap: cap}
}

// Enqueue 入队。
func (f *FIFO) Enqueue(c Customer) error {
	if f.cap > 0 && len(f.buf) >= f.cap {
		return ErrFull
	}
	f.buf = append(f.buf, c)
	return nil
}

// Dequeue 出队。
func (f *FIFO) Dequeue() (Customer, error) {
	if len(f.buf) == 0 {
		return Customer{}, ErrEmpty
	}
	c := f.buf[0]
	f.buf = f.buf[1:]
	return c, nil
}

// Len 当前排队人数。
func (f *FIFO) Len() int { return len(f.buf) }

// IsFull 是否已满。
func (f *FIFO) IsFull() bool { return f.cap > 0 && len(f.buf) >= f.cap }

// IsEmpty 是否为空。
func (f *FIFO) IsEmpty() bool { return len(f.buf) == 0 }

// Priority 优先级队列（最小堆按 Priority 字段）。
type Priority struct {
	h   priorityHeap
	cap int
}

// NewPriority 创建容量为 cap 的优先级队列。
func NewPriority(cap int) *Priority {
	return &Priority{cap: cap}
}

// Enqueue 入队。
func (p *Priority) Enqueue(c Customer) error {
	if p.cap > 0 && p.h.Len() >= p.cap {
		return ErrFull
	}
	heap.Push(&p.h, c)
	return nil
}

// Dequeue 出队（优先级最高者）。
func (p *Priority) Dequeue() (Customer, error) {
	if p.h.Len() == 0 {
		return Customer{}, ErrEmpty
	}
	return heap.Pop(&p.h).(Customer), nil
}

// Len 当前排队人数。
func (p *Priority) Len() int { return p.h.Len() }

// IsFull 是否已满。
func (p *Priority) IsFull() bool { return p.cap > 0 && p.h.Len() >= p.cap }

// IsEmpty 是否为空。
func (p *Priority) IsEmpty() bool { return p.h.Len() == 0 }

// priorityHeap 优先级堆。
type priorityHeap []Customer

func (h priorityHeap) Len() int           { return len(h) }
func (h priorityHeap) Less(i, j int) bool { return h[i].Priority < h[j].Priority }
func (h priorityHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *priorityHeap) Push(x any)        { *h = append(*h, x.(Customer)) }
func (h *priorityHeap) Pop() any {
	old := *h
	n := len(old)
	c := old[n-1]
	*h = old[:n-1]
	return c
}

// SJF 最短作业优先队列（按 ServiceDur 排序）。
type SJF struct {
	h   sjfHeap
	cap int
}

// NewSJF 创建 SJF 队列。
func NewSJF(cap int) *SJF { return &SJF{cap: cap} }

// Enqueue 入队。
func (s *SJF) Enqueue(c Customer) error {
	if s.cap > 0 && s.h.Len() >= s.cap {
		return ErrFull
	}
	heap.Push(&s.h, c)
	return nil
}

// Dequeue 出队。
func (s *SJF) Dequeue() (Customer, error) {
	if s.h.Len() == 0 {
		return Customer{}, ErrEmpty
	}
	return heap.Pop(&s.h).(Customer), nil
}

// Len 当前人数。
func (s *SJF) Len() int { return s.h.Len() }

// IsFull 是否已满。
func (s *SJF) IsFull() bool { return s.cap > 0 && s.h.Len() >= s.cap }

// IsEmpty 是否为空。
func (s *SJF) IsEmpty() bool { return s.h.Len() == 0 }

type sjfHeap []Customer

func (h sjfHeap) Len() int           { return len(h) }
func (h sjfHeap) Less(i, j int) bool { return h[i].ServiceDur < h[j].ServiceDur }
func (h sjfHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *sjfHeap) Push(x any)        { *h = append(*h, x.(Customer)) }
func (h *sjfHeap) Pop() any {
	old := *h
	n := len(old)
	c := old[n-1]
	*h = old[:n-1]
	return c
}
