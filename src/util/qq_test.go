package util

import (
	"bug-carrot/config"
	"fmt"
	"testing"
)

func TestQQSend(t *testing.T) {
	QQSend(2039799616, "hello")
}

func TestEmoji(t *testing.T) {
	for i := 0; i <= 222; i++ {
		QQGroupSend(config.C.Plugin.Homework.Group, fmt.Sprintf("%d[CQ:face,id=%d]", i, i))
	}
}
