package zclipper

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"strings"
	"sync"
	"sync/atomic"
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
 @Time    : 2024/7/19 -- 11:44
 @Author  : bishop ❤️ MONEY
 @Description: z clipper
*/

type ZClipper interface {
	Add(ctx context.Context, api string, reloadFunc ReloadInt64Func, fallbackFunc FallbackInt64Func, intervalFunc IntervalFunc) error
	Allow(ctx context.Context, api string) bool
	Print()
}

var (
	ErrNotEnoughArgs           = errors.New("not enough args")
	ErrBadArgs                 = errors.New("bad args")
	ErrFallbackFail            = errors.New("fallback fail")
	ErrInitFailNoReloadControl = errors.New("init fail no reload control")
	ErrNotValidConf            = errors.New("not valid config")
	ErrWrongInterval           = errors.New("wrong interval")
	ErrReload2EmptyKey         = errors.New("reload to empty key")
)

const keySep string = "-"

type ReloadInt64Func func(ctx context.Context, key string) (int64, error)
type FallbackInt64Func func(ctx context.Context, key string) (int64, error)
type ReloadInt64PairFunc func(ctx context.Context, key string) (int64, int64, error)
type FallbackInt64PairFunc func(ctx context.Context, key string) (int64, int64, error)
type IntervalFunc func(ctx context.Context, key string) (time.Duration, error)

func DefaultReloadPairFunc(limit, bucket int64) ReloadInt64PairFunc {
	return func(ctx context.Context, key string) (int64, int64, error) {
		// todo assure max limit or max bucket ?
		return limit, bucket, nil
	}
}

func DefaultReloadFunc(arg int64) ReloadInt64Func {
	return func(ctx context.Context, key string) (int64, error) {
		return arg, nil
	}
}

func DefaultFallbackPairFunc(limit, bucket int64) FallbackInt64PairFunc {
	return func(ctx context.Context, key string) (int64, int64, error) {
		return limit, bucket, nil
	}
}

func DefaultFallbackFunc(arg int64) FallbackInt64Func {
	return func(ctx context.Context, key string) (int64, error) {
		return arg, nil
	}
}

func defaultIntervalFunc(interval time.Duration) IntervalFunc {
	return func(ctx context.Context, key string) (time.Duration, error) {
		return interval, nil
	}
}

func DefaultIntervalFuncWrapper(interval time.Duration) IntervalFunc {
	if interval < time.Second {
		return defaultIntervalFunc(time.Second)
	}
	return defaultIntervalFunc(interval)
}

func RedisReloadPairFunc(redisCli *redis.Client) ReloadInt64PairFunc {
	return func(ctx context.Context, key string) (int64, int64, error) {
		if rs, err := redisCli.Get(ctx, key).Result(); err != nil {
			return 0, 0, err
		} else {
			ss := strings.Split(rs, ":")
			if len(ss) != 2 {
				return 0, 0, fmt.Errorf("invalid Limite limit: %s", rs)
			}
			var limit, bucket int64
			if val, err := cast.ToInt64E(ss[0]); err != nil {
				return 0, 0, err
			} else {
				limit = val
			}
			if val, err := cast.ToInt64E(ss[1]); err != nil {
				return 0, 0, err
			} else {
				bucket = val
			}
			return limit, bucket, nil
		}
	}
}

func RedisReloadFunc(redisCli *redis.Client) ReloadInt64Func {
	return func(ctx context.Context, key string) (int64, error) {
		if rs, err := redisCli.Get(ctx, key).Result(); err != nil {
			return 0, err
		} else {
			if val, err := cast.ToInt64E(rs); err != nil {
				return 0, err
			} else {
				return val, nil
			}
		}
	}
}

type Mark struct {
	done atomic.Uint32
	sync.Mutex
}

func (m *Mark) Doing() bool {
	return m.done.Load() == 1
}

func (m *Mark) Done() {
	m.Lock()
	defer m.Unlock()
	m.done.CompareAndSwap(1, 0)
}

func (m *Mark) Do() {
	m.Lock()
	defer m.Unlock()
	m.done.CompareAndSwap(0, 1)
}

func buildKey(group, namespace, api string) string {
	return strings.Join([]string{
		group,
		namespace,
		api,
	}, keySep)
}
