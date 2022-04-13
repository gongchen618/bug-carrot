package util

import (
	"fmt"
	"testing"
)

var messages = []string{
	//"今天武汉的天气怎么样",
	//"这周的CF",
	//"今晚的Codeforces",
	//"这周的作业",
	//"微积分作业是什么",
	//"reset微积分",
	//"submit微积分",
	"卡洛早安",
	"卡洛晚安",
	"卡洛喜欢吃萝卜吗",
	"卡洛说不好",
	"卡洛收拾卷怪可是超厉害的哦",
	"大学物理作业是什么呢",
	"离散数学是坏的",
	"卡洛知道明天武汉的天气吗?",
	"卡洛知道这几天的天气吗",
	"卡洛这三天天气",
	"卡洛这两天的天气是什么",
}

func TestGetWordsFromString(t *testing.T) {
	for _, message := range messages {
		words := GetWordsFromString(message)
		fmt.Println(words)
	}
}

/*
- 天气：[t 今天] [ns 武汉] [n 天气]
- 作业：[n 作业]
- cf: [eng codeforces/cf] [t 今晚]
[{x 卡洛} {r 这} {m 三天} {n 天气}]
[{x 卡洛} {r 这} {m 两天} {uj 的} {n 天气} {v 是} {r 什么}]
*/
