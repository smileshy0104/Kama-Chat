package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// encryptAES 使用AES算法加密数据。
// 该函数需要三个参数：待加密的数据(data)、加密密钥(key)和初始化向量(iv)。
// 所有参数都以字节切片的形式提供。
// 函数返回加密后的数据，以Base64编码的字符串形式输出，以及一个错误对象。
func encryptAES(data, key, iv []byte) (string, error) {
	// 创建AES密钥对应的加密器。
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建一个足够大的字节切片来存储加密后的数据，包括初始化向量的长度。
	ciphertext := make([]byte, aes.BlockSize+len(data))
	// 复制初始化向量到一个新的切片，以避免修改原始的初始化向量。
	ivCopy := make([]byte, aes.BlockSize)
	copy(ivCopy, iv)
	// 创建一个基于CFB模式的加密器。
	stream := cipher.NewCFBEncrypter(block, ivCopy)
	// 使用加密器对数据进行加密。
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	// 将初始化向量和加密后的数据合并，并以Base64编码，然后返回编码后的字符串。
	return base64.StdEncoding.EncodeToString(append(ivCopy, ciphertext[aes.BlockSize:]...)), nil
}
