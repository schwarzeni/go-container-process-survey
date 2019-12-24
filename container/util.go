package container

import (
	"math/rand"
	"time"
)

// randStringBytes 生成指定长度的随机字符串
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
