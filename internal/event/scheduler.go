// scheduler.go 扩展 event 包，提供事件调度和时间管理工具。
package event

// BatchScheduler 批量事件调度器：一次性调度多个事件。
type BatchScheduler struct {
	queue *Queue
}

// NewBatchScheduler 创建批量调度器。
func NewBatchScheduler() *BatchScheduler {
	return &BatchScheduler{queue: &Queue{}}
}

// Schedule 添加一个事件。
func (bs *BatchScheduler) Schedule(e Event) {
	bs.queue.Push(e)
}

// ScheduleMany 批量添加事件。
func (bs *BatchScheduler) ScheduleMany(events []Event) {
	for _, e := range events {
		bs.queue.Push(e)
	}
}

// Next 取出下一个事件。
func (bs *BatchScheduler) Next() (Event, bool) {
	return bs.queue.Pop()
}

// Pending 返回待处理事件数。
func (bs *BatchScheduler) Pending() int {
	return bs.queue.Len()
}

// PeekTime 查看下一个事件的时间但不弹出。
func (bs *BatchScheduler) PeekTime() (float64, bool) {
	if bs.queue.Len() == 0 {
		return 0, false
	}
	// Pop and push back to peek.
	e, ok := bs.queue.Pop()
	if !ok {
		return 0, false
	}
	bs.queue.Push(e)
	return e.Time, true
}

// Clear 清空所有待处理事件。
func (bs *BatchScheduler) Clear() {
	bs.queue = &Queue{}
}

// Clock 是一个虚拟时钟，跟踪仿真时间。
type Clock struct {
	Now float64
}

// NewClock 创建初始时间为 0 的时钟。
func NewClock() *Clock {
	return &Clock{}
}

// Advance 将时钟推进到新时间。
func (c *Clock) Advance(t float64) {
	if t > c.Now {
		c.Now = t
	}
}

// Elapsed 返回自仿真开始经过的时间。
func (c *Clock) Elapsed() float64 {
	return c.Now
}

// EventCounter 事件计数器，按类型统计事件数量。
type EventCounter struct {
	counts map[int]int
}

// NewEventCounter 创建事件计数器。
func NewEventCounter() *EventCounter {
	return &EventCounter{counts: make(map[int]int)}
}

// Count 记录一个事件。
func (ec *EventCounter) Count(kind int) {
	ec.counts[kind]++
}

// Get 获取某种类型的事件数量。
func (ec *EventCounter) Get(kind int) int {
	return ec.counts[kind]
}

// Total 返回所有事件总数。
func (ec *EventCounter) Total() int {
	total := 0
	for _, c := range ec.counts {
		total += c
	}
	return total
}

// Reset 重置计数器。
func (ec *EventCounter) Reset() {
	ec.counts = make(map[int]int)
}
