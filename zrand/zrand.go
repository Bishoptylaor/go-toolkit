package zrand

import (
	"math/rand"
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
 @Time    : 2024/7/13 -- 14:53
 @Author  : bishop ❤️ MONEY
 @Description: zrand.go
*/

var RanD random

type random struct {
	rd *rand.Rand
}

var rd *rand.Rand

const (
	PROBABILITY int32 = 10000
)

func (r *random) InitProb() (err error) {
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
	return nil
}

func (r *random) GetRandomNumber(n int32) int32 {
	return r.getRandomNumber(n)
}

func (r *random) getRandomNumber(n int32) int32 {
	if n == 0 {
		n = PROBABILITY
	}
	return rd.Int31n(n)
}

func (r *random) CheckProbabilityJackpot(prob int32) bool {
	number := r.getRandomNumber(PROBABILITY)
	if number < prob {
		return true
	}
	return false
}

func (r *random) calRandomNumberBetween(min, max int32) int32 {
	// 认为只要调用数量至少为1
	if max == 0 {
		return 1
	}
	number := r.getRandomNumber(max - min)
	cal := min + number
	if cal != 0 {
		return cal
	}
	return 1
}
