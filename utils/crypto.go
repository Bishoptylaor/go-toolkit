package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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
 @Time    : 2024/7/13 -- 14:00
 @Author  : bishop ❤️ MONEY
 @Software: GoLand
 @Description: 加解密函数
*/

// 填充明文
func _PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// 去除填充数据
func _PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	origData = _PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AES加密 结果返回base64编码后的string
func AesEncryptWithString(plainText string, key string) (string, error) {
	crypted, err := AesEncrypt([]byte(plainText), []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(crypted), nil
}

// AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = _PKCS5UnPadding(origData)
	return origData, nil
}

// AES解密, 输入的加密串需是已经用base64编码过的
func AesDecryptWithString(secretText string, key string) (string, error) {
	crypted, err := base64.StdEncoding.DecodeString(secretText)
	if err != nil {
		return "", err
	}
	res, err := AesDecrypt(crypted, []byte(key))
	if err != nil {
		return "", err
	}
	return string(res), nil
}