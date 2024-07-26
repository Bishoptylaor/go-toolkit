package zcrypto

import (
	"bytes"
	"fmt"
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
 @Time    : 2024/7/26 -- 11:38
 @Author  : bishop ❤️ MONEY
 @Description: padding
*/

type Pad interface {
	// Padding 填充明文
	Padding(plaintext []byte, blockSize int) []byte
	// UnPadding 去除填充数据
	UnPadding(src []byte) ([]byte, error)
}

var PKCS5 pkcs5
var PKCS7 pkcs5
var Zero zero
var AnsiX923 ansiX923
var Iso7816 iso7816

// PKCS5Padding（也称为 PKCS7Padding）是一种基于 PKCS#5 标准的填充方案，用于加密数据以确保其长度符合加密算法的块大小要求。
// 这种填充方法主要用于块加密模式，如 CBC（Cipher Block Chaining，密码块链接模式）。
type pkcs5 struct{}

func GetPKCS5() Pad {
	return pkcs5{}
}

func (pkcs5) Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padText...)
}

func (pkcs5) UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)], nil
}

// 也称为 Null Padding，简单地在数据末尾添加零直到满足块大小。
type zero struct{}

func (zero) Padding(plaintext []byte, blockSize int) []byte {
	paddingLength := blockSize - len(plaintext)%blockSize
	padText := bytes.Repeat([]byte{0x00}, paddingLength)
	return append(plaintext, padText...)
}

func (zero) UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 || src[length-1] != 0x00 {
		return nil, fmt.Errorf("invalid zero padding")
	}
	unPadding := int(src[length-1])
	return src[:length-unPadding], nil
}

// ANSI X.923 Padding：
// 与 PKCS#5 填充类似，但填充值是随机的，而不是固定模式。
type ansiX923 struct{}

func (ansiX923) Padding(plaintext []byte, blockSize int) []byte {
	paddingLength := blockSize - len(plaintext)%blockSize
	padText := make([]byte, paddingLength)
	for i := range padText {
		padText[i] = byte(paddingLength)
	}
	return append(plaintext, padText...)
}

func (ansiX923) UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("invalid X.923 padding")
	}
	unPadding := int(src[length-1])
	if unPadding > length {
		return nil, fmt.Errorf("invalid X.923 padding")
	}
	return src[:length-unPadding], nil
}

// 用于智能卡和金融服务，填充模式为 0x80 后跟一系列 0x00。
type iso7816 struct{}

func (iso7816) Padding(plaintext []byte, blockSize int) []byte {
	paddingLength := blockSize - len(plaintext)%blockSize
	padText := []byte{0x80}
	if paddingLength > 1 {
		padText = append(padText, bytes.Repeat([]byte{0x00}, paddingLength-1)...)
	}
	return append(plaintext, padText...)
}

func (iso7816) UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 || src[0] != 0x80 {
		return nil, fmt.Errorf("invalid ISO/IEC 7816-4 padding")
	}
	for i := 1; i < length; i++ {
		if src[i] != 0x00 {
			return src[:i], nil
		}
	}
	return nil, fmt.Errorf("invalid ISO/IEC 7816-4 padding")
}
