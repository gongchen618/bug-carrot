package util

import (
	"bug-carrot/config"
	"fmt"
	"math/rand"
	"time"
)

var (
	Token      string
	AdminGroup = int64(1028801782)
)

func init() {
	rand.Seed(time.Now().UnixNano())

	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < 4; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	Token = string(result)

	QQGroupSend(AdminGroup, fmt.Sprintf("卡洛塔已重启！生成了新的前端 token：%s", Token))
	QQSend(config.C.Plugin.Default.Admin, Token)
}
