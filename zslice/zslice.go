// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package cast provides easy and safe casting in Go.

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
 @Time    : 2024/7/15 -- 14:43
 @Author  : bishop ❤️ MONEY
 @Description: slice 常用函数封装
*/

// Distinct 对 src 去重
func Distinct[T comparable](src []T) (dst []T) {
	tempMap := map[T]struct{}{}
	for _, t := range src {
		// directly use hmap.count
		l := len(tempMap)
		tempMap[t] = struct{}{}

		// if len changes, means the item we add to tempMap is first seen in the progress.
		if len(tempMap) != l {
			dst = append(dst, t)
		}
	}
	return
}

// Contains 判断 tar 是否在 src 中
func Contains[T comparable](src []T, tar T) bool {
	for _, t := range src {
		if t == tar {
			return true
		}
	}
	return false
}

// Except 返回在 left 但是不在 right 中的 left 元素
func Except[T comparable](left []T, right []T) (dst []T) {
	tempMap := make(map[T]struct{}, len(right))
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

// Reverse 原地倒序
func Reverse[T any](src []T) {
	i := 0
	j := len(src) - 1

	for i < j {
		src[i], src[j] = src[j], src[i]
		i++
		j--
	}
}

// Cut 左开右闭
func Cut[T any](src []T, start, end int64) []T {

	return nil
}
