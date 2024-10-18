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
 @Time    : 2024/7/13 -- 14:14
 @Author  : bishop ❤️ MONEY
 @Description: 时间窗口校验，是否还在窗口内
*/

package xtime

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	TimeTypeHand      = 0
	TimeTypeAbs       = 1
	TimeTypeNatureDay = 2
)

func CheckNeedUpdate(updateType int64, extra string, tm int64) (bool, error) {
	return CheckNeedUpdate2(updateType, extra, tm, time.Now().Unix())
}

func CheckNeedUpdate2(updateType int64, extra string, start, end int64) (bool, error) {
	var t TimeChecker
	switch updateType {
	case TimeTypeHand:
		t = &HandTimeChecker{}
	case TimeTypeAbs:
		t = &AbsTimeChecker{}
	case TimeTypeNatureDay:
		t = &NatureDayTimeChecker{}
	default:
		return false, errors.New("not support updatetype")
	}
	if extra != "" {
		if err := json.Unmarshal([]byte(extra), t); err != nil {
			return false, err
		}
	}
	if err := t.Valid(); err != nil {
		return false, err
	}
	return t.Check(start, end), nil
}

func TimeCheckerFactory(updateType int64, extra string) (TimeChecker, error) {
	var t TimeChecker
	switch updateType {
	case TimeTypeHand:
		t = &HandTimeChecker{}
	case TimeTypeAbs:
		t = &AbsTimeChecker{}
	case TimeTypeNatureDay:
		t = &NatureDayTimeChecker{}
	default:
		return nil, errors.New("not support updatetype")
	}
	if extra != "" {
		if err := json.Unmarshal([]byte(extra), t); err != nil {
			return nil, err
		}
	}
	return t, nil
}

func TimeCheckerValid(updateType int64, extra string) error {
	_, err := TimeCheckerFactory(updateType, extra)
	return err
}

type TimeChecker interface {
	Check(start, end int64) bool
	Valid() error
}

// ==== 手动更新
type HandTimeChecker struct{}

func (HandTimeChecker) Check(start, end int64) bool {
	return false
}

func (HandTimeChecker) Valid() error {
	return nil
}

// ==== 绝对时间
type AbsTimeChecker struct {
	Value int64 `json:"value"`
}

func (a *AbsTimeChecker) Check(start, end int64) bool {
	if start+a.Value <= end {
		return true
	}
	return false
}

func (a *AbsTimeChecker) Valid() error {
	if a.Value <= 0 {
		return errors.New("lag value must > 0")
	}
	return nil
}

// ==== 自然日
type NatureDayTimeChecker struct {
	Value int `json:"value"`
}

func (n *NatureDayTimeChecker) Check(start, end int64) bool {
	_start := time.Unix(start, 0).AddDate(0, 0, n.Value).Unix()
	if DayBeginStamp(_start) <= DayBeginStamp(end) {
		return true
	}
	return false
}

func (n *NatureDayTimeChecker) Valid() error {
	if n.Value <= 0 {
		return errors.New("lag value must > 0")
	}
	return nil
}
