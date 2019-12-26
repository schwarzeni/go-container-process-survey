package container

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
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

// getEnvsByPid 获取进程的环境变量
func getEnvsByPid(pid string) (envs []string, err error) {
	var (
		path         = fmt.Sprintf("/proc/%s/environ", pid)
		contentBytes []byte
	)

	if contentBytes, err = ioutil.ReadFile(path); err != nil {
		return nil, fmt.Errorf("read file %s failed, %v", path, err)
	}
	envs = strings.Split(string(contentBytes), "\u0000")
	return
}
