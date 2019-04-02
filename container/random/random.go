package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	src         = rand.NewSource(time.Now().UnixNano())
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func RandStringBytes(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

// 使用余数
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

// 掩码
func RandStringBytesMask(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}

	return string(b)
}

// 掩码加强版
func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)

	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// math source
// 那就是提高随机数的产生
// 直接使用rand.Source，而不是全局或者共享的随机源
// 全局的(默认的)随机源是线程安全，里面用到了锁，所以没有我们直接rand.Source更好
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func main() {
	var randomStr string
	n := 16

	randomStr = RandStringRunes(n)
	fmt.Println("RandStringRunes : ", randomStr)

	randomStr = RandStringBytes(n)
	fmt.Println("RandStringBytes : ", randomStr)

	randomStr = RandStringBytesRmndr(n)
	fmt.Println("RandStringBytesRmndr : ", randomStr)

	randomStr = RandStringBytesMask(n)
	fmt.Println("RandStringBytesMask : ", randomStr)

	randomStr = RandStringBytesMaskImpr(n)
	fmt.Println("RandStringBytesMaskImpr : ", randomStr)

	randomStr = RandStringBytesMaskImprSrc(n)
	fmt.Println("RandStringBytesMaskImprSrc : ", randomStr)

}
