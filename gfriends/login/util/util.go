package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	fllog "github.com/forlifeproj/msf/log"
)

func IntArray2Bytes(intArray []int) []byte {
	var byteArray []byte

	for _, n := range intArray {
		// 创建一个 bytes.Buffer，用于存储转换后的字节
		buf := new(bytes.Buffer)

		// 将整数 n 转换为字节并写入 buf，使用大端字节序
		err := binary.Write(buf, binary.BigEndian, int32(n))
		if err != nil {
			fllog.Log().Error(fmt.Sprintf("binary.Write failed:", err))
		}

		// 获取字节切片
		intBytes := buf.Bytes()

		// 输出结果
		// fmt.Printf("整数 %d 转换为字节切片: %v\n", n, intBytes)

		// 将 intBytes 追加到 byteArray 中
		byteArray = append(byteArray, intBytes...)
	}
	fllog.Log().Debug(fmt.Sprintf("intArray:%+v byteArray:%s", intArray, string(byteArray)))
	return byteArray
}

func Bytes2IntArray(byteArray []byte) []int {
	intArray := make([]int, len(byteArray)/4)
	buf := bytes.NewReader(byteArray)

	for i := 0; i < len(intArray); i++ {
		var n int32
		err := binary.Read(buf, binary.BigEndian, &n)
		if err != nil {
			fllog.Log().Error(fmt.Sprintf("binary.Read failed:", err))
		}
		intArray[i] = int(n)
	}

	return intArray
}

func GetCrc32(src string) uint32 {
	return crc32.ChecksumIEEE([]byte(src))
}

func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
