// Package event 提供离散事件仿真所需的事件类型与未来事件表（FEL）。
package event

import "container/heap"

// 事件种类。
const (
	Arrival = iota
	Departure
)

// Event 表示一个离散事件：在 Time 时刻发生，与 Customer 编号关联。
type Event struct {
	Time     float64
	Kind     int
	Customer int
}

// felHeap 是实现 heap.Interface 的内部最小堆，
// 按 Time 升序；Time 相同则 Arrival 先于 Departure。
type felHeap []Event

func (h felHeap) Len() int { return len(h) }

func (h felHeap) Less(i, j int) bool {
	if h[i].Time != h[j].Time {
		return h[i].Time < h[j].Time
	}
	return h[i].Kind < h[j].Kind
}

func (h felHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *felHeap) Push(x any) {
	*h = append(*h, x.(Event))
}

func (h *felHeap) Pop() any {
	old := *h
	n := len(old)
	e := old[n-1]
	*h = old[:n-1]
	return e
}

// Queue 是基于最小堆的未来事件表。
type Queue struct {
	h felHeap
}

// Push 将事件加入 FEL。
func (q *Queue) Push(e Event) {
	heap.Push(&q.h, e)
}

// Pop 取出并返回时间最早的事件；空表返回 (零值, false)，不 panic。
func (q *Queue) Pop() (Event, bool) {
	if q.h.Len() == 0 {
		return Event{}, false
	}
	return heap.Pop(&q.h).(Event), true
}

// Len 返回 FEL 中事件数量。
func (q *Queue) Len() int {
	return q.h.Len()
}
