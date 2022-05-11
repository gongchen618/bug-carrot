package plugin

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller"
	"bug-carrot/param"
	"bug-carrot/util"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/togatoga/goforces"
)

type codeforces struct {
	Index param.PluginIndex
}

func (p *codeforces) GetPluginName() string {
	return p.Index.PluginName
}
func (p *codeforces) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *codeforces) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *codeforces) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *codeforces) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *codeforces) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *codeforces) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *codeforces) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

// 以上六个函数在注册时被唯一调用，并以此为依据加入相应的 queue
// 无需修改

// IsTime : 是你需要的时间吗？
func (p *codeforces) IsTime() bool {
	return false
}

// DoTime : 当到了你需要的时间，要做什么呢？
func (p *codeforces) DoTime() error {
	return nil
}

// IsMatchedGroup : 是你想收到的群 @ 消息吗？
func (p *codeforces) IsMatchedGroup(msg param.GroupMessage) bool {
	return true
}

// DoMatchedGroup : 收到了想收到的群 @ 消息，要做什么呢？
func (p *codeforces) DoMatchedGroup(msg param.GroupMessage) error {
	if !config.C.RiskControl {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), constant.CarrotGroupPuzzled)
	} else {
		util.QQSend(config.C.Plugin.Default.Admin, constant.CarrotRiskControlAngry)
	}
	return nil
}

// IsMatchedPrivate : 是你想收到的私聊消息吗？
func (p *codeforces) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return msg.WordsMap.ExistWord("eng", []string{"cf", "codeforces"})
}

// DoMatchedPrivate : 收到了想收到的私聊消息，要做什么呢？
func (p *codeforces) DoMatchedPrivate(msg param.PrivateMessage) error {
	if msg.SubType == "friend" {
		if config.C.RiskControl {
			util.QQSend(msg.UserId, constant.CarrotRiskControlAngry)
		} else {
			util.QQSend(msg.UserId, getCodeforcesContestList())
		}
	}
	return nil
}

// Listen : 监听到非 @ 的群消息，要做什么呢？
func (p *codeforces) Listen(msg param.GroupMessage) {
	if !msg.WordsMap.ExistWord("eng", []string{"cf", "codeforces"}) {
		return
	}
	util.QQGroupSend(msg.GroupId, getCodeforcesContestList())
}

// Close : 项目要关闭了，要做什么呢？
func (p *codeforces) Close() {
}

func CodeforcesPluginRegister() {
	p := &codeforces{
		Index: param.PluginIndex{
			PluginName:            "codeforces",
			PluginAuthor:          "ligen131",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   false,
			FlagCanMatchedPrivate: true,
			FlagCanListen:         true,
			FlagUseDatabase:       false,
			FlagIgnoreRiskControl: false,
		},
	}
	controller.PluginRegister(p)
}

func getCodeforcesContestList() string {
	// { DurationSeconds: 7200 Frozen: false ID: 1681 Name:Educational Codeforces Round 129 (Rated for Div. 2) Phase:BEFORE RelativeTimeSeconds: -1046232 StartTimeSeconds: 1653316500 Type:ICPC }
	ctx := context.Background()
	logger := log.New(os.Stderr, "[Goforces] ", log.LstdFlags)
	api, err := goforces.NewClient(logger)
	if err != nil {
		util.ErrorPrint(err, nil, "[Goforces] Failed to initial goforces.")
		return "Something error. Please try again later."
	}
	contestList, err := api.GetContestList(ctx, nil)
	if err != nil {
		util.ErrorPrint(err, nil, "[Goforces] Failed to get contests list.")
		return "Failed to get contests list. Please try again later."
	}
	ListLen := len(contestList)
	text := "Codeforces Upcoming Contests:"
	tot := 0
	for i := 0; i < ListLen; i++ {
		if contestList[i].Before() {
			tot++
			dur, _ := time.ParseDuration(fmt.Sprintf("%ds", contestList[i].DurationSeconds))
			st := time.Unix(contestList[i].StartTimeSeconds, 0).Format("2006-01-02 15:04:05") // Do not change this time
			text += fmt.Sprintf("\n\n%d. %v\n%v, Duration: %v\n%v", tot, contestList[i].Name, st, dur, contestList[i].ContestURL())
		}
	}
	if tot == 0 {
		text += "\nNone"
	}
	return text
}
