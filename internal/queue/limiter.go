// limiter.go 提供令牌桶和漏桶速率限制器，用于控制系统入口流量。
package queue

import (
	"errors"
	"math"
)

// ErrNoToken 没有可用令牌。
var ErrNoToken = errors.New("limiter: no token available")

// TokenBucket 令牌桶限流器。
type TokenBucket struct {
	capacity   float64 // 桶容量
	tokens     float64 // 当前令牌数
	rate       float64 // 令牌生成速率（个/单位时间）
	lastUpdate float64 // 上次更新时间
}

// NewTokenBucket 创建令牌桶。
func NewTokenBucket(capacity, rate float64) *TokenBucket {
	return &TokenBucket{capacity: capacity, tokens: capacity, rate: rate}
}

// TryConsume 在时间 now 尝试消耗一个令牌。
func (tb *TokenBucket) TryConsume(now float64) bool {
	tb.refill(now)
	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// Available 返回当前可用令牌数。
func (tb *TokenBucket) Available(now float64) float64 {
	tb.refill(now)
	return tb.tokens
}

func (tb *TokenBucket) refill(now float64) {
	elapsed := now - tb.lastUpdate
	if elapsed > 0 {
		tb.tokens += elapsed * tb.rate
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastUpdate = now
	}
}

// LeakyBucket 漏桶限流器：以固定速率输出，溢出则丢弃。
type LeakyBucket struct {
	capacity float64 // 桶容量（可缓冲的最大量）
	level    float64 // 当前水位
	leakRate float64 // 漏出速率
	lastLeak float64 // 上次漏出时间
}

// NewLeakyBucket 创建漏桶。
func NewLeakyBucket(capacity, leakRate float64) *LeakyBucket {
	return &LeakyBucket{capacity: capacity, leakRate: leakRate}
}

// TryAdd 在时间 now 尝试向桶中加入 amount 的水。
func (lb *LeakyBucket) TryAdd(now, amount float64) bool {
	lb.leak(now)
	if lb.level+amount > lb.capacity {
		return false
	}
	lb.level += amount
	return true
}

// Level 返回当前水位。
func (lb *LeakyBucket) Level(now float64) float64 {
	lb.leak(now)
	return lb.level
}

func (lb *LeakyBucket) leak(now float64) {
	elapsed := now - lb.lastLeak
	if elapsed > 0 {
		leaked := elapsed * lb.leakRate
		lb.level -= leaked
		if lb.level < 0 {
			lb.level = 0
		}
		lb.lastLeak = now
	}
}

// SlidingWindowCounter 滑动窗口计数器限流。
type SlidingWindowCounter struct {
	window    float64 // 窗口长度
	limit     int     // 窗口内最大请求数
	timestamps []float64
}

// NewSlidingWindow 创建滑动窗口限流器。
func NewSlidingWindow(window float64, limit int) *SlidingWindowCounter {
	return &SlidingWindowCounter{window: window, limit: limit}
}

// Allow 检查时间 now 是否允许通过。
func (sw *SlidingWindowCounter) Allow(now float64) bool {
	// 清理过期时间戳。
	cutoff := now - sw.window
	valid := 0
	for _, ts := range sw.timestamps {
		if ts > cutoff {
			sw.timestamps[valid] = ts
			valid++
		}
	}
	sw.timestamps = sw.timestamps[:valid]
	if len(sw.timestamps) >= sw.limit {
		return false
	}
	sw.timestamps = append(sw.timestamps, now)
	return true
}

// Count 当前窗口内的请求数。
func (sw *SlidingWindowCounter) Count(now float64) int {
	cutoff := now - sw.window
	n := 0
	for _, ts := range sw.timestamps {
		if ts > cutoff {
			n++
		}
	}
	return n
}

// WaitTime 估算下一个请求需要等待多久。
func (sw *SlidingWindowCounter) WaitTime(now float64) float64 {
	if sw.Allow(now) {
		// 已经通过了，撤回。
		sw.timestamps = sw.timestamps[:len(sw.timestamps)-1]
		return 0
	}
	if len(sw.timestamps) == 0 {
		return 0
	}
	// 最早的时间戳过期后即可通过。
	cutoff := now - sw.window
	earliest := math.MaxFloat64
	for _, ts := range sw.timestamps {
		if ts > cutoff && ts < earliest {
			earliest = ts
		}
	}
	return (earliest + sw.window) - now
}
