// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 21:17
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : sha1.go

package Sha1

import (
	"bytes"
	"crypto/des"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
)

// SignKey DES加解密密钥
const SignKey = "wT2uRiKm"

// Sha1 哈希加密
func Sha1(s string) string {
	// 产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})
	h := sha1.New()
	// 写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组
	h.Write([]byte(s))
	// 这个用来得到最终地散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// Md5 MD5加密
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// DESPlus DES加密
func DESPlus(s string) string {
	src := []byte(s)
	block, err := des.NewCipher([]byte(SignKey))
	if err != nil {
		zap.L().Error("[DES加密]" + err.Error())
		return ""
	}
	bs := block.BlockSize()
	src = ZeroPadding(src, bs)
	if len(src)%bs != 0 {
		zap.L().Error("[DES加密]Need a multiple of the blocksize")
		return ""
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out)
}

// DESReduce DES解密
func DESReduce(s string) string {
	src, err := hex.DecodeString(s)
	if err != nil {
		zap.L().Error("[DES解密]" + err.Error())
		return ""
	}
	block, err := des.NewCipher([]byte(SignKey))
	if err != nil {
		zap.L().Error("[DES解密]" + err.Error())
		return ""
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		zap.L().Error("[DES加密]input not full blocks")
		return ""
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = ZeroUnPadding(out)
	return string(out)
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
