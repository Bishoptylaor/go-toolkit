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
 @Time    : 2024/7/13 -- 13:58
 @Author  : bishop ❤️ MONEY
 @Description: 给运行的函数f封装，避免panic导致全局退出
*/

package safefun

import (
	"fmt"
	"runtime"
	"time"
)

// Retry 重试 func 最大次数，间隔
func Retry(callback func() error, maxRetries int, interval time.Duration) (err error) {
	for i := 1; i <= maxRetries; i++ {
		if err = callback(); err != nil {
			time.Sleep(interval)
			continue
		}
		return
	}
	return
}

// FunWrapperWithArgs 带参数函数 wrapper
func FunWrapperWithArgs(f func(args ...interface{}), args ...interface{}) (err error) {

	defer func() {
		if p := recover(); p != nil {
			err = DumpStack(p)
		}
	}()

	f(args...)
	return
}

// FunWrapper 无参数 wrapper
func FunWrapper(f func()) (err error) {

	defer func() {
		if p := recover(); p != nil {
			err = DumpStack(p)
		}
	}()

	f()
	return
}

func DumpStack(e interface{}) (err error) {
	if e == nil {
		return
	}
	err = fmt.Errorf("%+v", e)
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		err = fmt.Errorf("%s\t %s:%d", err.Error(), file, line)
	}
	return
}

func SafelyGo(function func(), handleError func(error)) {
	go func() {
		err := FunWrapper(function)
		if err != nil && handleError != nil {
			handleError(err)
		}
	}()
}
