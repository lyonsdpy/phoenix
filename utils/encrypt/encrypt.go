// @Author:DaiPengyuan
// @Date:2018/9/20
// @Desc: Des秘钥加解密

package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
)

const k = "lyonsdpy"

var key = []byte(k)

func pKCS5Padding(src []byte, blocksize int) []byte {
	padding := blocksize - len(src)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func DesEncrypt(origStr string) (string, error) {
	origData := []byte(origStr)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	origData = pKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return hex.EncodeToString(crypted), nil
}

func DesDeCrypt(cryptedStr string) (r string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
			return
		}
	}()
	if cryptedStr != "" {
		crypted, err := hex.DecodeString(cryptedStr)
		if err != nil {
			return "", err
		}
		block, err := des.NewCipher(key)
		if err != nil {
			return "", err
		}
		blockMode := cipher.NewCBCDecrypter(block, key)
		origData := make([]byte, len(crypted))
		blockMode.CryptBlocks(origData, crypted)
		origData = pKCS5UnPadding(origData)
		return string(origData), nil
	} else {
		return "", errors.New("input cryptedStr is nil")
	}
}

// 快速解密,如果解密成功则返回解密后的值,解密失败返回原始值
func FastDecrypt(s string) string {
	ds, err := DesDeCrypt(s)
	if err != nil {
		return s
	}
	return ds
}
