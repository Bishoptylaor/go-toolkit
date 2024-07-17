package zlimiter

import (
	"context"
	"fmt"
	"github.com/Bishoptylaor/go-toolbox/zlimiter/zbucket"
	"sync"
	"time"
)

/*
*  ┏┓      ┏┓
*┏━┛┻━━━━━━┛┻┓
*┃　　　━　　  ┃
*┃   ┳┛ ┗┳   ┃
*┃           ┃
*┃     ┻     ┃
*┗━━━┓     ┏━┛
*　　 ┃　　　┃神兽保佑
*　　 ┃　　　┃代码无BUG！
*　　 ┃　　　┗━━━┓
*　　 ┃         ┣┓
*　　 ┃         ┏┛
*　　 ┗━┓┓┏━━┳┓┏┛
*　　   ┃┫┫  ┃┫┫
*      ┗┻┛　 ┗┻┛
@Time    : 2024/7/16 -- 14:44
@Author  : bishop ❤️ MONEY
@Description: limiter.go
*/

type TokenBuckets struct {
	bs map[string]*zbucket.Bucket
	mu sync.Mutex
}

// NewTokenBuckets 创建一个新的令牌桶实例
func NewTokenBuckets(capacity int, rate float64) ZLimiter {
	return &TokenBuckets{
		bs: make(map[string]*zbucket.Bucket),
		mu: sync.Mutex{},
	}
}

func (tbs *TokenBuckets) Add(ctx context.Context, api string, cap, rate float64) {
	tbs.mu.Lock()
	defer tbs.mu.Unlock()
	tbs.bs[api] = &zbucket.Bucket{
		capacity:   cap,
		rate:       rate,
		rateWindow: time.Second,
		tokens:     cap,
		lastTime:   time.Time{},
		mu:         sync.Mutex{},
	}
}

func (tbs *TokenBuckets) Available(ctx context.Context, api string) bool {
	return false
}

// Take 从桶中取出一定数量的令牌
// 如果桶中令牌不足，等待直到有足够的令牌
func (tb *zbucket.TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	// 计算自上次以来应该添加的令牌数
	passed := now.Sub(tb.lastTime).Seconds()
	addTokens := tb.rate * passed

	// 更新桶中的令牌数，不超过桶的容量
	tb.tokens = tb.tokens + addTokens
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}

	tb.lastTime = now

	// 如果桶中的令牌足够，取出令牌，返回true
	if tb.tokens >= 1 {
		tb.tokens-- // 取出一个令牌
		return true
	}

	// 如果桶中的令牌不足，返回false
	return false
}

// RateLimiter 限流器结构
type RateLimiter struct {
	buckets map[string]*zbucket.TokenBucket // 接口名到令牌桶的映射
	mu      sync.Mutex                      // 互斥锁
}

// SetRate 为指定接口设置令牌桶的速率和容量
func (rl *RateLimiter) SetRate(api string, capacity int, rate float64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.buckets[api] = NewTokenBucket(capacity, rate)
}

// Allow 检查是否可以处理指定接口的请求
func (rl *RateLimiter) Allow(api string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[api]
	if !exists {
		// 如果没有为该接口设置令牌桶，可以认为请求是被允许的
		return true
	}
	return bucket.Take()
}

func main() {
	// 创建限流器实例
	limiter := NewRateLimiter()

	// 为不同的接口设置不同的速率和容量
	limiter.SetRate("/api1", 10, 5) // 桶大小为10，每秒添加5个令牌
	limiter.SetRate("/api2", 5, 2)  // 桶大小为5，每秒添加2个令牌

	// 模拟请求
	for i := 0; i < 20; i++ {
		api := "/api1" // 假设所有请求都是对/api1的
		if limiter.Allow(api) {
			fmt.Printf("Request to %s is allowed at %v\n", api, time.Now())
		} else {
			fmt.Printf("Request to %s is denied at %v\n", api, time.Now())
		}
	}
}
