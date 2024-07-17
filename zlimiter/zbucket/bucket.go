package zbucket

import (
	"context"
	"github.com/redis/go-redis/v9"
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
 @Time    : 2024/7/16 -- 11:44
 @Author  : bishop ❤️ MONEY
 @Description: bucket
*/

type Bucket interface {
	Refill(ctx context.Context)
	Available(ctx context.Context) bool
}

// TokenBucket 令牌桶结构
type TokenBucket struct {
	capacity   float64       // 桶的大小
	rate       float64       // 添加的令牌
	rateWindow time.Duration // 每一段时间添加的令牌
	tokens     float64       // 当前桶中的令牌容量
	lastTime   time.Time     // 上次增加令牌容量的时间
	mu         sync.Mutex    // lock
}

// LeakyBucket 令牌桶结构
type LeakyBucket struct {
	capacity      float64       // 桶的大小
	leakRate      time.Duration // 每一段时间添加的容量
	highWaterMark float64       // 当前桶中的水位线
	lastTime      time.Time     // 上次增加容量的时间
	mu            sync.Mutex    // lock
}

type SpeedGovernor struct {
	mu       sync.Mutex
	queue    []uint8 // 请求序列，0 表示不允许当前请求；1 表示允许当前请求。queue 中的 0:1 = req-allow:allow
	tag      string  // 标识
	req      int32   // 每 req 个请求
	allow    int32   // 允许通过 allow 个
	redisCli *redis.Client
}
