package zcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/Bishoptylaor/go-toolbox/zutils"
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
 @Description: 加解密函数
*/

// AesEncrypt 加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	origData = PKCS5.Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 初始向量的长度必须等于块block的长度16字节
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// AesEncryptWithString 加密 结果返回base64编码后的string
func AesEncryptWithString(plainText string, key string) (string, error) {
	encrypted, err := AesEncrypt([]byte(plainText), []byte(key))
	if err != nil {
		return "", err
	}
	return zutils.Bytes2Str(encrypted), nil
}

// AesDecrypt 解密
func AesDecrypt(encrypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData, _ = PKCS5.UnPadding(origData)
	return origData, nil
}

// AesDecryptWithString 解密, 输入的加密串需是已经用base64编码过的
func AesDecryptWithString(secretText string, key string) (string, error) {
	encrypted, err := Base64.Decode(secretText)
	if err != nil {
		return "", err
	}
	res, err := AesDecrypt(encrypted, []byte(key))
	if err != nil {
		return "", err
	}
	return string(res), nil
}
