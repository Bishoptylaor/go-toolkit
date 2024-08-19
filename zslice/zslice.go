package zslice

import "slices"

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
 @Time    : 2024/7/15 -- 14:43
 @Author  : bishop ❤️ MONEY
 @Description: 对 go/src/slice 的补充 https://github.com/golang/go/blob/master/src/slices/slices.go
*/

func Distinct[S ~[]E, E comparable](src S) (dst S) {
	tempMap := make(map[E]struct{})
	for _, v := range src {
		l := len(tempMap)
		tempMap[v] = struct{}{}
		if len(tempMap) != l {
			dst = append(dst, v)
		}
	}
	return
}

func DistinctFunc[S ~[]E, E any](src S, eq func(a, b E) bool) (dst S) {
	for _, v := range src {
		if slices.ContainsFunc(dst, func(e E) bool {
			return eq(v, e)
		}) {
			dst = append(dst, v)
		}
	}
	return
}

// Except 返回在 left 但是不在 right 中的 left 元素
func Except[S ~[]E, E comparable](left S, right S) (dst S) {
	tempMap := make(map[E]struct{}, len(right))
	for _, r := range right {
		tempMap[r] = struct{}{}
	}
	for _, l := range left {
		if _, ok := tempMap[l]; !ok {
			dst = append(dst, l)
		}
	}
	return
}

func ExceptFunc[S ~[]E, E any](left S, right S, eq func(a, b E) bool) (dst S) {
	// 在 left
	for _, v := range left {
		// 不在 right
		if !slices.ContainsFunc(right, func(e E) bool {
			return eq(v, e)
		}) {
			// 加入结果集
			dst = append(dst, v)
		}
	}
	return dst
}

func RPop[S ~[]E, E any](src S) (last E) {
	if len(src) == 0 {
		return
	}
	last = src[len(src)-1]
	src = src[:len(src)-1]
	return
}

func LPop[S ~[]E, E any](src S) (first E) {
	if len(src) == 0 {
		return
	}
	first = src[0]
	src = src[1:]
	return
}

func FrontAdd[S ~[]E, E any](src S, added S) S {
	return append(added, src...)
}

func BackendAdd[S ~[]E, E any](src S, added S) S {
	return append(src, added...)
}
