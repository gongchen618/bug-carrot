package util

import (
	"fmt"
	"testing"
	"time"
)

var messages = []string{
	//"今天武汉的天气怎么样",
	//"这周的CF",
	//"今晚的Codeforces",
	//"这周的作业",
	//"微积分作业是什么",
	//"reset微积分",
	//"submit微积分",
	"近期考试",
	"卡洛晚安",
	"召唤一个你",
	"立契",
	"悔契",
	"卡洛喜欢吃萝卜吗",
	"卡洛说不好",
	"卡洛收拾卷怪可是超厉害的哦",
	"大学物理作业是什么呢",
	"离散数学是坏的",
	"卡洛知道明天武汉的天气吗?",
	"卡洛这两天的天气是什么",
	"任务清单",
	"TODO清单",
	"约定",
	"这是一份约定",
	"div1榜单",
	"榜单div2",
}

func TestGetWordsFromString(t *testing.T) {
	for _, message := range messages {
		words := GetWordsFromMessage(message)
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

func TestTimePattern(t *testing.T) {
	timePattern := "2006年01月02日15:04"
	dateStr := "2022年02月03日03:14"
	date, err := time.ParseInLocation(timePattern, dateStr, time.Local)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(date)
}
