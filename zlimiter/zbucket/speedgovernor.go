package zbucket

import (
	"context"
	"github.com/Bishoptylaor/go-toolbox/zslice"
	"github.com/redis/go-redis/v9"
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
 @Time    : 2024/7/17 -- 10:05
 @Author  : bishop ❤️ MONEY
 @Description: speedgovernor.go
*/

func NewSpeedGovernor(tag string, req, allow int32, redisCli *redis.Client) (Bucket, bool) {
	// 最多100% allow 应 <= req
	if tag == "" || req == 0 {
		return nil, false
	}
	if allow > req {
		allow = req
	}
	return &SpeedGovernor{
		redisCli: redisCli,
		req:      req,
		allow:    allow,
		queue:    make([]uint8, 0, req),
		tag:      tag,
	}, true
}

func (pb *SpeedGovernor) Refill(ctx context.Context) {
	if err := pb.redisCli.Ping(ctx).Err(); err != nil {
		pb.refillLocal(ctx)
		return
	}
	pb.refillRedis(ctx)
	return
}

func (pb *SpeedGovernor) refillLocal(ctx context.Context) {

}

func (pb *SpeedGovernor) refillRedis(ctx context.Context) {

}

func (pb *SpeedGovernor) Available(ctx context.Context) bool {
	if err := pb.redisCli.Ping(ctx).Err(); err != nil {
		return pb.availableLocal(ctx)
	}
	return pb.availableRedis(ctx)
}

func (pb *SpeedGovernor) availableRedis(ctx context.Context) bool {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if len(pb.queue) == 0 {
		pb.refillRedis(ctx)
	}
	current := zslice.RPop(pb.queue)
	return current == 1
}

func (pb *SpeedGovernor) availableLocal(ctx context.Context) bool {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if len(pb.queue) == 0 {
		pb.refillLocal(ctx)
	}
	current := zslice.RPop(pb.queue)
	return current == 1
}
