package safe

import (
	"fmt"
	"runtime"
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
 @Time    : 2024/7/13 -- 13:58
 @Author  : bishop ❤️ MONEY
 @Software: GoLand
 @Description: 给运行的函数f封装，避免panic导致全局退出
*/

// FunWrapperWithArgs 带参数函数 wrapper
func FunWrapperWithArgs(f func(args ...interface{}), args ...interface{}) (err error) {

	defer func() {
		if panic := recover(); panic != nil {
			err = DumpStack(panic)
		}
	}()

	f(args...)
	return
}

// FunWrapper 无参数 wrapper
func FunWrapper(f func()) (err error) {

	defer func() {
		if panic := recover(); panic != nil {
			err = DumpStack(panic)
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
