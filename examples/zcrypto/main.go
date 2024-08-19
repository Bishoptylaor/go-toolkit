package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	. "github.com/Bishoptylaor/go-toolkit/zcrypto"
)

// online aes check https://www.lddgo.net/en/encrypt/aes

func main() {
	TestBase()
	TestAESCBCEncrypt()
	TestAESCBCDecrypt()
	TestAESCFBEncrypt()
	TestAESCFBDecrypt()
	TestAESGCMDecryptWithNonce()
}

func TestBase() {
	key := []byte("1234567890abcdefghijklmnopqrstuv") // 32字节的密钥
	origData := []byte("Hello, AES encryption!")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// 填充原始数据以满足块大小要求
	paddedText := PKCS7.Padding(origData, aes.BlockSize)

	// 加密
	cbcEncrypter := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
	cbcEncrypter.CryptBlocks(paddedText, paddedText)
	fmt.Printf("Encrypted: %x\n", paddedText)

	url := "https://a/b/c/thisisaverylongurlwaitingtobeshorter"
	b62 := Base62.SEncode(url)
	fmt.Println("b62: ", b62)
	fmt.Println(Base62.SDecode(b62))
}

func TestAESCBCEncrypt() {
	var testTbl = []struct {
		origData  []byte
		key       []byte
		iv        []byte
		encrypted string
	}{
		{
			origData:  []byte("test data"),
			key:       []byte("test-key-aes-128"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "ecddfa122db3975e5534b3809a9d8def",
		},
		{
			origData:  []byte("test data"),
			key:       []byte("test-key-aes-192-0000000"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "a94071f5ad281cf4d296abfd7d280bea",
		},
		{
			origData:  []byte("test data"),
			key:       []byte("test-key-aes-192-000000000000000"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "042959da3c38bf27a0183b62705c55ea",
		},
	}

	var pad = PKCS5

	for _, test := range testTbl {
		var encrypted, err = AESCBCEncrypt(test.origData, test.key, test.iv, pad)
		if err != nil {
			fmt.Println(err)
		}

		var r = hex.EncodeToString(encrypted)

		if r != test.encrypted {
			fmt.Printf("AES CBC 加密 %s 结果，期望: %s, 实际: %s \n", string(test.origData), test.encrypted, r)
		}
	}
}

func TestAESCBCDecrypt() {
	var testTbl = []struct {
		origData  string
		key       []byte
		iv        []byte
		encrypted string
	}{
		{
			origData:  "test data",
			key:       []byte("test-key-aes-128"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "ecddfa122db3975e5534b3809a9d8def",
		},
		{
			origData:  "test data",
			key:       []byte("test-key-aes-192-0000000"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "a94071f5ad281cf4d296abfd7d280bea",
		},
		{
			origData:  "test data",
			key:       []byte("test-key-aes-192-000000000000000"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "042959da3c38bf27a0183b62705c55ea",
		},
	}

	var padding = PKCS5

	for _, test := range testTbl {
		var encrypted, _ = hex.DecodeString(test.encrypted)

		var origData, err = AESCBCDecrypt(encrypted, test.key, test.iv, padding)
		if err != nil {
			fmt.Println(err)
		}

		var r = string(origData)

		if r != test.origData {
			fmt.Printf("AES CBC 解密 %s 结果，期望: %s, 实际: %s \n", test.encrypted, test.origData, r)
		}
	}
}

func TestAESCFBEncrypt() {
	var testTbl = []struct {
		origData  []byte
		key       []byte
		iv        []byte
		encrypted string
	}{
		{
			origData:  []byte("test data"),
			key:       []byte("test-key-aes-128"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "610a6e98e8cc15fc7a984e5945d76a7a",
		},
		{
			origData:  []byte("test data"),
			key:       []byte("test-key-aes-192-0000000"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "5f36e20df6d953e50acc24f45fc7cc53",
		},
		{
			origData:  []byte("test data"),
			key:       []byte("test-key-aes-192-000000000000000"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "2dc97774aab66260870b385251024abd",
		},
	}

	var padding = PKCS5

	for _, test := range testTbl {
		var encrypted, err = AESCFBEncrypt(test.origData, test.key, test.iv, padding)
		if err != nil {
			fmt.Println(err)
		}

		var r = hex.EncodeToString(encrypted)

		if r != test.encrypted {
			fmt.Printf("AES CFB 加密 %s 结果，期望: %s, 实际: %s \n", string(test.origData), test.encrypted, r)
		}
	}
}

func TestAESCFBDecrypt() {
	var testTbl = []struct {
		origData  string
		key       []byte
		iv        []byte
		encrypted string
	}{
		{
			origData:  "test data",
			key:       []byte("test-key-aes-128"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "610a6e98e8cc15fc7a984e5945d76a7a",
		},
		{
			origData:  "test data",
			key:       []byte("test-key-aes-192-0000000"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "5f36e20df6d953e50acc24f45fc7cc53",
		},
		{
			origData:  "test data",
			key:       []byte("test-key-aes-192-000000000000000"),
			iv:        []byte("1234567890abcdef"),
			encrypted: "2dc97774aab66260870b385251024abd",
		},
	}

	var padding = PKCS5

	for _, test := range testTbl {
		var encrypted, _ = hex.DecodeString(test.encrypted)

		var origData, err = AESCFBDecrypt(encrypted, test.key, test.iv, padding)
		if err != nil {
			fmt.Println(err)
		}

		var r = string(origData)

		if r != test.origData {
			fmt.Printf("AES CFB 解密 %s 结果，期望: %s, 实际: %s \n", test.encrypted, test.origData, r)
		}
	}
}

func TestAESGCMDecryptWithNonce() {
	var testTbl = []struct {
		encrypted string
		origData  string
		nonce     []byte
		key       []byte
	}{
		{
			encrypted: "01b7207ef05ec95d5650ddf34bfc45f936c49c35",
			origData:  "test",
			nonce:     []byte("1234567890ab"),
			key:       []byte("test-key-aes-128"),
		},
		{
			encrypted: "1db73f6671091ab1cfda13e51e544a1e89923709a9",
			origData:  "hello",
			nonce:     []byte("1234567890ab"),
			key:       []byte("test-key-aes-128"),
		},
	}

	for _, test := range testTbl {
		var encrypted, _ = hex.DecodeString(test.encrypted)

		var origData, err = AESGCMDecryptWithNonce(encrypted, test.key, test.nonce, nil)
		if err != nil {
			fmt.Println(err)
		}

		var r = string(origData)

		if r != test.origData {
			fmt.Printf("AES GCM 解密 %s 结果，期望: %s, 实际: %s \n", test.encrypted, test.origData, r)
		}
	}
}
