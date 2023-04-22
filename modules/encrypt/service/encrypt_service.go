package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"playground-go-api/domain"
)

type encryptService struct {
	encryptKey string
	aes256IV   string
}

func NewEncryptService() domain.EncryptService {
	return &encryptService{
		encryptKey: "a567ba1bd93e861554a2bb2ae3931243c128ddf818d4e1246823a3450f81687e",
		aes256IV:   "b2fe646ba775db180098d5266849d258",
	}
}

func (e *encryptService) NewCBCDecrypter(value string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover success.")
		}
	}()
	key, _ := hex.DecodeString(e.encryptKey)
	ciphertext, _ := hex.DecodeString(value)
	aes256IV, _ := hex.DecodeString(e.aes256IV)

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	// crypto/cipher: ciphertext too short
	if len(ciphertext) < block.BlockSize() {
		return ""
	}
	// crypto/cipher: input not full blocks
	if len(ciphertext)%block.BlockSize() != 0 {
		return ""
	}

	mode := cipher.NewCBCDecrypter(block, aes256IV)
	decrypted := make([]byte, len(ciphertext))

	// crypto/cipher: output smaller than input
	if len(decrypted) < len(ciphertext) {
		return ""
	}

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(decrypted, ciphertext)

	result := e.PKCS5UnPadding(decrypted)
	return string(result)
}

func (e *encryptService) NewCBCEncrypter(value string) string {
	key, _ := hex.DecodeString(e.encryptKey)
	plaintext := []byte(value)
	iv, _ := hex.DecodeString(e.aes256IV)

	block, _ := aes.NewCipher(key)

	content := e.PKCS5Padding(plaintext, block.BlockSize())
	encrypted := make([]byte, len(content))

	// crypto/cipher: ciphertext too short
	if len(content) < block.BlockSize() {
		return ""
	}
	// crypto/cipher: input not full blocks
	if len(content)%block.BlockSize() != 0 {
		return ""
	}
	// crypto/cipher: output smaller than input
	if len(encrypted) < len(content) {
		return ""
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(encrypted, content)
	return hex.EncodeToString(encrypted)
}

/* func (e *encryptService) NewCFBDecrypter(value string) string {
	key, _ := hex.DecodeString(e.encryptKey)
	ciphertext, _ := hex.DecodeString(value)

	block, _ := aes.NewCipher(key)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < block.BlockSize() {
		return ""
	}
	iv := ciphertext[:block.BlockSize()]
	ciphertext = ciphertext[block.BlockSize():]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func (e *encryptService) NewCFBEncrypter(value string) string {
	key, _ := hex.DecodeString(e.encryptKey)
	plaintext := []byte(value)

	block, _ := aes.NewCipher(key)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, block.BlockSize()+len(plaintext))
	iv := ciphertext[:block.BlockSize()]

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[block.BlockSize():], plaintext)
	return hex.EncodeToString(ciphertext)
} */

/* ==== */

// PKCS5UnPadding 解包装
func (e *encryptService) PKCS5UnPadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

// PKCS5Padding PKCS5包装
func (e *encryptService) PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (e *encryptService) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (e *encryptService) PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
