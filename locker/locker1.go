package locker

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
 @Time    : 2024/7/13 -- 14:05
 @Author  : bishop ❤️ MONEY
 @Software: GoLand
 @Description: locker1.go
*/

import (
	"context"
	"errors"
	"fmt"
	"gitlab.pri.ibanyu.com/middleware/seaweed/xcache/redisext"
	"gitlab.pri.ibanyu.com/middleware/seaweed/xlog"
	"math/rand"
	"time"
)

var NotNeedLockError = errors.New("need not lock")

type Locker interface {
	Lock(name string) error
	Unlock(name string) error
}

// 二重判定锁
func DoubleLockUtil(ctx context.Context, key string, locker Locker, doCheck, operation func() error) error {
	if err := doCheck(); err != nil {
		return err
	}
	if err := locker.Lock(key); err != nil {
		return err
	}
	defer locker.Unlock(key)
	if err := doCheck(); err != nil {
		return err
	}
	return operation()
}

// 一重判定锁
func SingleLockUtil(ctx context.Context, key string, locker Locker, doCheck, operation func() error) error {
	if err := locker.Lock(key); err != nil {
		return err
	}
	defer locker.Unlock(key)
	if err := doCheck(); err != nil {
		return err
	}
	return operation()
}

type RedisLocker struct {
	Redis *redisext.RedisExt
	ttl   time.Duration
}

func (r *RedisLocker) Lock(key string) error {
	b, err := r.Redis.SetNX(context.TODO(), key, 1, r.ttl)
	if err != nil {
		return err
	}
	if !b {
		return errors.New("try lock failed")
	}
	return nil
}

func (r *RedisLocker) Unlock(key string) error {
	_, err := r.Redis.Del(context.TODO(), key)
	return err
}

func NewRedisLocker(redis *redisext.RedisExt, ttl time.Duration) Locker {
	return &RedisLocker{
		Redis: redis,
		ttl:   ttl,
	}
}

type Lockers struct {
	Redis *redisext.RedisExt
	Rd    *rand.Rand
}

func NewLockers(redis *redisext.RedisExt, rd *rand.Rand) *Lockers {
	return &Lockers{
		Redis: redis,
		Rd:    rd,
	}
}

func (l *Lockers) UnLock(ctx context.Context, lockKey string, value int32, fun string) {
	r, err := l.Redis.Unlock(ctx, lockKey, value)
	if err != nil {
		xlog.Warnf(ctx, "%s key:%s try unlock fail:%s,%s", fun, lockKey, r, err)
	}
}
func (l *Lockers) Lock(ctx context.Context, fun, lockKey string, lockTime time.Duration) (int32, error) {
	fun += "Lock -->"
	value := l.Rd.Int31n(100)
	lock, err := l.Redis.Lock(ctx, lockKey, value, lockTime)
	if err != nil {
		xlog.Errorf(ctx, "%s key:%s try lock fail:%s", fun, lockKey, err)
		return value, err
	}
	if !lock {
		return value, fmt.Errorf("multiclick")
	}
	return value, nil
}

// 循环等待获取锁
func (l *Lockers) LockCyclic(ctx context.Context, fun, lockKey string, lockTime time.Duration) (int32, error) {
	fun += "LockCyclic-->"
	value := l.Rd.Int31n(100)
	lock, err := l.Redis.Lock(ctx, lockKey, value, lockTime)
	if err != nil {
		xlog.Errorf(ctx, "%s key:%s try lock fail:%s", fun, lockKey, err)
		return 0, err
	}
	// 加锁失败重新加锁
	if !lock {
		for i := 0; i < 500; i++ {
			time.Sleep(10 * time.Millisecond)
			lock, err = l.Redis.Lock(ctx, lockKey, value, lockTime)
			if err != nil || !lock {
				xlog.Infof(ctx, "%v 获取锁失败，重新获取。lockKey：%v err:%s", fun, lockKey, err)
			} else if lock {
				return value, nil
			}
		}
		return value, fmt.Errorf("multiclick")
	}
	return value, nil
}
