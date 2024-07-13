package zstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
 @Time    : 2024/7/13 -- 14:54
 @Author  : bishop ❤️ MONEY
 @Software: GoLand
 @Description: zredis.go
*/

var RedisCounter rCounter

type rCounter struct{}

func (r *rCounter) PackAddOne(ctx context.Context, counterKey string, fun string, dur time.Duration) (int32, error) {
	current := r.Get(ctx, counterKey, fun)
	if current > 0 {
		r.AddOne(ctx, counterKey, fun)
	} else {
		_, err := r.Set(ctx, counterKey, fun, dur)
		if err != nil {
			return 0, err
		}
		r.AddOne(ctx, counterKey, fun)
	}
	return 0, nil
}
func (r *rCounter) AddOneWithoutSet(ctx context.Context, counterKey string, fun string) {
	current := r.Get(ctx, counterKey, fun)
	if current > 0 {
		r.AddOne(ctx, counterKey, fun)
	} else {
		return
	}
}
func (rCounter) AddOne(ctx context.Context, counterKey string, fun string) {
	r, err := config.RedisClient.Incr(ctx, counterKey)
	if err != nil {
		xlog.Warnf(ctx, "%s key:%s try incr fail:%s,%s", fun, counterKey, r, err)
	}
}
func (rCounter) Get(ctx context.Context, counterKey string, fun string) (val int32) {
	fun += "Get -->"
	r, err := config.RedisClient.Get(ctx, counterKey)
	if err != nil {
		xlog.Infof(ctx, "%s key:%s try Get fail:%s,%s", fun, counterKey, r, err)
		return
	}
	v, err := strconv.ParseInt(r, 10, 32)
	if err != nil {
		xlog.Warnf(ctx, "%s key:%s try convert str2int32 fail:%d,%s", fun, counterKey, v, err)
		return
	}
	val = int32(v)
	return
}
func (rCounter) Set(ctx context.Context, fun, lockKey string, lockTime time.Duration) (int32, error) {
	fun += "Set -->"
	_, err := config.RedisClient.Set(ctx, lockKey, 0, lockTime)
	if err != nil {
		xlog.Errorf(ctx, "%s key:%s try Set fail:%s", fun, lockKey, err)
		return 0, err
	}
	return 0, nil
}
func (rCounter) SetVal(ctx context.Context, fun, lockKey string, val interface{}, lockTime time.Duration) (int32, error) {
	fun += "SetVal -->"
	_, err := config.RedisClient.Set(ctx, lockKey, val, lockTime)
	if err != nil {
		xlog.Errorf(ctx, "%s key:%s try Set fail:%s", fun, lockKey, err)
		return 0, err
	}
	return 0, nil
}

var RedisCache rc

type rc struct{}

func (r *rc) Set(ctx context.Context, key string, payload interface{}, cacheTime time.Duration) error {
	// fun := "RedisCache.Set -->"
	if key == "" {
		return fmt.Errorf("missing Key")
	}
	s, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = config.RedisClient.Set(ctx, key, string(s), cacheTime)
	if err != nil {
		return fmt.Errorf("set err:%s", err)
	}
	return nil
}
func (r *rc) Get(ctx context.Context, key string) (payload string, err error) {
	// fun := "RedisCache.Get -->"
	if key == "" {
		return "", fmt.Errorf("missing Key")
	}
	payload, err = config.RedisClient.Get(ctx, key)
	return payload, err
}
func (r *rc) Del(ctx context.Context, key string) (err error) {
	// fun := "RedisCache.Get -->"
	if key == "" {
		return fmt.Errorf("missing Key")
	}
	_, err = config.RedisClient.Del(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (r *rc) Zadd(ctx context.Context, key, mainstay string, score float64) error {
	fun := "RedisCache.Zadd -->"
	if key == "" {
		return fmt.Errorf("missing Key")
	}
	member := redisext.Z{Score: score, Member: mainstay}
	var members []redisext.Z
	members = append(members, member)
	_, err := config.RedisClient.ZAdd(ctx, key, members)
	if err != nil {
		return fmt.Errorf("%s Zadd value key %s value %s score %f err %s", fun, key, mainstay, score, err.Error())
	}
	return nil
}

func (r *rc) Zscore(ctx context.Context, key, mainstay string) (error, float64) {
	fun := "RedisCache.Zscore -->"
	if key == "" || mainstay == "" {
		return fmt.Errorf("missing Key"), 0
	}
	score, err := config.RedisClient.ZScore(ctx, key, mainstay)
	if err != nil {
		return fmt.Errorf("%s ZScore value key %s member %s err %s", fun, key, mainstay, err.Error()), 0
	}
	return nil, score
}

// func (r *rc) Zrangewithscores(ctx context.Context, key string, start, end int64) ([]*LeaderBoardData, error) {
// 	fun := "RedisCache.Zrangewithscores -->"
// 	if key == "" {
// 		return nil, fmt.Errorf("missing Key")
// 	}
// 	var rs []*LeaderBoardData
// 	members, err := config.RedisClient.ZRangeWithScores(ctx, key, start, end)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s ZRangeWithScores key %s err %s", fun, key, err.Error())
// 	}
// 	for i := 0; i < len(members); i++ {
// 		ms := members[i].Member.(string)
// 		if ms == "" {
// 			continue
// 		}
// 		data := &LeaderBoardData{
// 			Key:   ms,
// 			Score: members[i].Score,
// 		}
// 		rs = append(rs, data)
// 	}
// 	return rs, nil
// }

func (r *rc) Zrem(ctx context.Context, key, mainstay string) error {
	fun := "RedisCache.Zrem -->"
	if key == "" || mainstay == "" {
		return fmt.Errorf("missing Key")
	}
	member := make([]interface{}, 0)
	member = append(member, mainstay)
	_, err := config.RedisClient.ZRem(ctx, key, member)
	if err != nil {
		return fmt.Errorf("%s ZScore value key %s member %s err %s", fun, key, mainstay, err.Error())
	}
	return nil
}

func (r *rc) Zrank(ctx context.Context, key, mainstay string) (error, int64) {
	fun := "RedisCache.Zrank -->"
	if key == "" || mainstay == "" {
		return fmt.Errorf("missing Key"), -1
	}
	n, err := config.RedisClient.ZRank(ctx, key, mainstay)
	if err != nil {
		return fmt.Errorf("%s ZScore value key %s member %s err %s", fun, key, mainstay, err.Error()), -1
	}
	return nil, n
}

func (r *rc) Zcount(ctx context.Context, key, min, max string) (error, int64) {
	fun := "RedisCache.Zcount -->"
	if key == "" {
		return fmt.Errorf("missing Key"), -1
	}
	n, err := config.RedisClient.ZCount(ctx, key, min, max)
	if err != nil {
		return fmt.Errorf("%s ZCount value key %s err %s", fun, key, err.Error()), -1
	}
	return nil, n
}

func (r *rc) Hset(ctx context.Context, key, innerkey string, value int64) error {
	fun := "RedisCache.Hset -->"
	if key == "" {
		return fmt.Errorf("missing Key")
	}
	_, err := config.RedisClient.HSet(ctx, key, innerkey, value)
	if err != nil {
		return fmt.Errorf("%s Hset value key %s innerkey %s value %d  err %s", fun, key, innerkey, value, err.Error())
	}
	return nil
}

func (r *rc) Hget(ctx context.Context, key, innerkey string) (string, error) {
	fun := "RedisCache.Hget -->"
	if key == "" {
		return "", fmt.Errorf("missing Key")
	}
	res, err := config.RedisClient.HGet(ctx, key, innerkey)
	if err != nil {
		return "", fmt.Errorf("%s Hget value key %s innerkey %s err %s", fun, key, innerkey, err.Error())
	}
	return res, nil
}

func (r *rc) Hincby(ctx context.Context, key, innerkey string, value int64) error {
	fun := "RedisCache.Hincby -->"
	if key == "" {
		return fmt.Errorf("missing Key")
	}
	_, err := config.RedisClient.HIncrBy(ctx, key, innerkey, value)
	if err != nil {
		return fmt.Errorf("%s Hincby value key %s innerkey %s value %d  err %s", fun, key, innerkey, value, err.Error())
	}
	return nil
}

func (r *rc) Zcard(ctx context.Context, key string) (error, int64) {
	fun := "RedisCache.Zcard -->"
	if key == "" {
		return fmt.Errorf("missing Key"), -1
	}
	n, err := config.RedisClient.ZCard(ctx, key)
	if err != nil {
		return fmt.Errorf("%s ZCount value key %s err %s", fun, key, err.Error()), -1
	}
	return nil, n
}
