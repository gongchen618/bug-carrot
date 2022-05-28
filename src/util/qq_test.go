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
	message := ""
	for i := 0; i <= 222; i++ {
		message = fmt.Sprintf("%s[CQ:face,id=%d]", message, i)
	}
	QQGroupSend(config.C.Plugin.Schedule.Group, message)
}
