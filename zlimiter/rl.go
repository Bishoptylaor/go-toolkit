package zlimiter

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type RLimiter interface {
	Fill(ctx context.Context)
	Tick(ctx context.Context)
	TickDone()
	Available(ctx context.Context, redisCli *redis.Client) bool
}

type redisLimiter struct {
	mu         sync.Mutex
	redisCli   *redis.Client
	capacity   uint64             // capacity
	every      time.Duration      // reset to cap every duration
	lastFill   time.Time          // last refill time
	key        string             // key value to call redis
	cancelFunc context.CancelFunc // a way to stop ticker
}

func (rl *redisLimiter) Tick(ctx context.Context) {
	c, cancel := context.WithCancel(ctx)
	rl.cancelFunc = cancel
	go rl.Fill(c)
}

func (rl *redisLimiter) TickDone() {
	rl.cancelFunc()
}

func (rl *redisLimiter) Fill(ctx context.Context) {
	rl.fill(ctx)
	ticker := time.NewTicker(rl.every)
	for {
		select {
		case <-ticker.C:
			fmt.Println("ticking")
			rl.fill(ctx)
		case <-ctx.Done():
			fmt.Println("done")
			return
		}
	}
}

func (rl *redisLimiter) fill(ctx context.Context) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	val, err := rl.redisCli.Set(ctx, rl.key, rl.capacity, 24*time.Hour).Result()
	fmt.Println("filling", rl.key, val, err, "at", time.Now().Unix(), "last fill", rl.lastFill.Unix())
	if err != nil {
		//log.Ctx(ctx).Error().Msg(fmt.Sprintf("trying to refill key:%s with:%d fail, err:%s", rl.key, rl.capacity, err))
		// retry
		return
	}
	rl.lastFill = time.Now()
	return
}

func (rl *redisLimiter) Available(ctx context.Context, redisCli *redis.Client) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	val, err := redisCli.Decr(ctx, rl.key).Result()
	fmt.Println(val, err)
	switch {
	case err != nil:
		// something wrong.
		//log.Ctx(ctx).Error().Msg(fmt.Sprintf("trying to incr key:%s err:%s", rl.key, err))
		return true
	case val == -1:
		// cannot find key. val means the result after incr
		redisCli.Expire(ctx, rl.key, 24*time.Hour)
		return true
	case val < 0:
		// out of capacity
		return false
	default:
		// maybe something wrong. but we do not care
		return true
	}
}
