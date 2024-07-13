package ztime

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
 @Time    : 2024/7/13 -- 14:09
 @Author  : bishop ❤️ MONEY
 @Software: GoLand
 @Description: 时间函数封装
 @TODO: 结构整理 & 增加时区逻辑
*/

import (
	"context"
	"fmt"
	"time"
)

// 本年初
func EndOfYear(now time.Time) time.Time {
	now = now.AddDate(1, 0, 0)
	t := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Add(-1)
	return t
}

// 本月末
func StartOfYear(now time.Time) time.Time {
	t := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	return t
}

// 本月末
func EndOfMonth(now time.Time) time.Time {
	now = now.AddDate(0, 1, 0)
	t := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Add(-1)
	return t
}

// 本月初
func StartOfMonth(now time.Time) time.Time {
	t := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return t
}

// 当天末
func EndOfDay(now time.Time) time.Time {
	now = now.AddDate(0, 0, 1)
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(-1)
	return t
}

// 当天初
func StartOfDay(now time.Time) time.Time {
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return t
}

// StartOfDayStampFromStr 获取指定天的时间范围
// 天格式 2006-01-02
// 为空时候返回当天的
func StartOfDayStampFromStr(day string) (int64, error) {
	nowt := time.Now()

	var begin int64
	if len(day) > 0 {
		tm, err := time.ParseInLocation("2006-01-02", day, nowt.Location())
		if err != nil {
			return 0, err
		}

		begin = tm.Unix()

	} else {
		begin = StartOfDay(nowt).Unix()

	}

	return begin, nil
}

// 本周初
func StartOfWeek(now time.Time) time.Time {
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	w := t.Weekday() - 1
	for w < 0 {
		w += 7
	}
	return t.AddDate(0, 0, int(w*-1))
}

// 本周末
func EndOfWeek(now time.Time) time.Time {
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	w := 8 - t.Weekday()
	if t.Weekday() == time.Sunday {
		w = 1
	}
	return t.AddDate(0, 0, int(w)).Add(-1)
}

// 时间转换 返回请求对应的起止时间，如果是永久, 则返回
// timeType 时间类型
// timeValue
//
func GetTimeSE(timeType int64, timeValue int64, now time.Time) (ST *time.Time, ET *time.Time, err error) {
	if timeType == TimeForever {
		return nil, nil, nil
	}
	a, b, err := getTimeSE(timeType, timeValue, now)
	if err != nil {
		return nil, nil, err
	}
	ST, ET = &a, &b
	return
}

func GetTimeSE2(timeType int64, timeValue int64, now time.Time) (ST int64, ET int64, err error) {
	a, b, err := GetTimeSE(timeType, timeValue, now)
	if err != nil {
		return 0, 0, err
	}
	if a != nil {
		ST = a.Unix()
	}
	if b != nil {
		ET = b.Unix()
	}
	return
}

func getTimeSE(timeType int64, timeValue int64, now time.Time) (ST time.Time, ET time.Time, err error) {
	switch timeType {
	case TimeAbs:
		ST = now.Add(time.Duration(timeValue*-1) * time.Second)
		ET = now
	case TimeDay:
		ST = StartOfDay(now)
		ET = EndOfDay(now.AddDate(0, 0, int(timeValue)))
	case TimeWeek:
		ST = StartOfWeek(now)
		ET = EndOfDay(now.AddDate(0, 0, int(timeValue*7)))
	case TimeMonth:
		ST = StartOfMonth(now)
		ET = EndOfMonth(now.AddDate(0, int(timeValue), 0))
	case TimeYear:
		ST = StartOfMonth(now)
		ET = EndOfMonth(now.AddDate(int(timeValue), 0, 0))
	case TimeBackAbs:
		ST = now
		ET = now.Add(time.Duration(timeValue) * time.Second)
	default:
		err = fmt.Errorf("not support timetype %d", timeType)
	}
	return
}

func GetLocalTimeFromString(ctx context.Context, format string, timestring string) (localtime time.Time) {
	loc, _ := time.LoadLocation("Local")
	localtime, _ = time.ParseInLocation(format, timestring, loc)
	return
}

func GetTimeFormStamp(strTime int64, format string) string {
	datetime := time.Unix(strTime, 0).Format(format)
	return datetime
}

// build format with string type and build a struct for this format func
const (
	TimeFormatYear   = "2006-01-02"
	TimeFormatDate   = "01-02"
	TimeFormatTime24 = "15:04:05"
	TimeFormatTime12 = "03:04:05"
	TimeFormatAll    = "2006-01-02 15:04:05"
)
