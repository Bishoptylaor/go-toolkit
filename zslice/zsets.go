package zslice

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
 @Time    : 2024/7/13 -- 14:43
 @Author  : bishop ❤️ MONEY
 @Description: 处理集合的各种关系操作。这里只考虑两个集合间的关系。
https://baike.baidu.com/item/%E9%9B%86%E5%90%88/2908117?fr=ge_ala
*/

import (
	"github.com/google/go-cmp/cmp"
)

/*
ZSets
Normally, the slices in ZSets requite to be distinct.
*/
type ZSets[T comparable] struct {
	A []T
	B []T
}

func NewZSets[T comparable](a []T, b []T) ZSets[T] {
	zs := ZSets[T]{
		A: a,
		B: b,
	}
	return zs
}

// Equal 比较是否相等，无论 A B 是否有重复元素
func (z ZSets[T]) Equal() bool {
	return cmp.Equal(z.A, z.B)
}

// Intersection 交集
func (z ZSets[T]) Intersection() []T {
	var is []T
	tempMap := make(map[T]struct{})
	for _, t := range z.B {
		tempMap[t] = struct{}{}
	}
	for _, t := range z.A {
		if _, ok := tempMap[t]; ok {
			is = append(is, t)
		}
	}
	return is
}

// Union ...
func (z ZSets[T]) Union() []T {
	u := make([]T, 0, len(z.A)+len(z.B))
	for _, t := range z.A {
		u = append(u, t)
	}
	for _, t := range z.B {
		u = append(u, t)
	}
	return u
}

// DistinctUnion 去重union
func (z ZSets[T]) DistinctUnion() []T {
	dst := make([]T, 0, len(z.A)+len(z.B))
	tempMap := make(map[T]struct{}, len(z.A)+len(z.B))
	for _, t := range z.A {
		l := len(tempMap)
		tempMap[t] = struct{}{}

		if len(tempMap) != l {
			dst = append(dst, t)
		}
	}
	for _, t := range z.B {
		l := len(tempMap)
		tempMap[t] = struct{}{}

		if len(tempMap) != l {
			dst = append(dst, t)
		}
	}

	return dst
}

// AExceptB 属于A而不属于B的元素组成的集合，称为B关于A的相对补集，记作A-B或A\B，即A-B={x|x∈A，且x∉B}。
func (z ZSets[T]) AExceptB() []T {
	return Except(z.A, z.B)
}

// BExceptA 属于B而不属于A的元素组成的集合，称为A关于B的相对补集，记作B-A或B\A，即B-A={x|x∈B，且x∉A}。
func (z ZSets[T]) BExceptA() []T {
	return Except(z.B, z.A)
}
