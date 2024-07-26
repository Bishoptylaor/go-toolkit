package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	. "github.com/Bishoptylaor/go-toolbox/zcrypto"
)

func main() {
	key := []byte("1234567890abcdefghijklmnopqrstuv") // 32字节的密钥
	plaintext := []byte("Hello, AES encryption!")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// 填充原始数据以满足块大小要求
	paddedText := PKCS7.Padding(plaintext, aes.BlockSize)

	// 加密
	cbcEncrypter := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
	cbcEncrypter.CryptBlocks(paddedText, paddedText)
	fmt.Printf("Encrypted: %x\n", paddedText)

	url := "https://a/b/c/thisisaverylongurlwaitingtobeshorter"
	b62 := Base62.SEncode(url)
	fmt.Println("b62: ", b62)
	fmt.Println(Base62.SDecode(b62))
}
