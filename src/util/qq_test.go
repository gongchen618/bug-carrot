package util

import (
	"testing"
)

func TestQQSend(t *testing.T) {
	QQSend(2039799616, "hello")
}

func TestEmoji(t *testing.T) {
	//for i := 301; i <= 400; i++ {
	//	QQGroupSend(config.C.Plugin.Homework.Group, fmt.Sprintf("%d[CQ:face,id=%d]", i, i))
	//}
}
