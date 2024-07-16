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
 @Time    : 2024/7/13 -- 14:11
 @Author  : bishop ❤️ MONEY
 @Description: 时间处理相关 const 定义集合
*/

const (
	TimeForever = 0 // 永久
	TimeAbs     = 1 // 绝对时间，单位为s
	TimeDay     = 2 // 自然日
	TimeWeek    = 3 // 自然周
	TimeMonth   = 4 // 自然月
	TimeYear    = 5 // 自然年
	TimeBackAbs = 6 // 绝对时间，单位为s, 从当前时间往后
)
