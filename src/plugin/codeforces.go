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
	"strings"
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

func (p *codeforces) IsTime() bool {
	return false
}

func (p *codeforces) DoTime() error {
	return nil
}

func (p *codeforces) IsMatchedGroup(msg param.GroupMessage) bool {
	return msg.WordsMap.ExistWord("eng", []string{"cf", "codeforces"})
}

func (p *codeforces) DoMatchedGroup(msg param.GroupMessage) error {
	if !config.C.RiskControl {
		util.QQGroupSendAtSomeone(msg.GroupId, util.GetQQGroupUserId(msg), getCodeforcesContestList(msg.RawMessage))
	} else {
		util.QQSend(config.C.Plugin.Default.Admin, constant.CarrotRiskControlAngry)
	}
	return nil
}

func (p *codeforces) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return msg.WordsMap.ExistWord("eng", []string{"cf", "codeforces"})
}

func (p *codeforces) DoMatchedPrivate(msg param.PrivateMessage) error {
	if msg.SubType == "friend" {
		if config.C.RiskControl {
			util.QQSend(msg.UserId, constant.CarrotRiskControlAngry)
		} else {
			util.QQSend(msg.UserId, getCodeforcesContestList(msg.RawMessage))
		}
	}
	return nil
}

func (p *codeforces) Listen(msg param.GroupMessage) {
}

func (p *codeforces) Close() {
}

func CodeforcesPluginRegister() {
	p := &codeforces{
		Index: param.PluginIndex{
			PluginName:            "codeforces",
			PluginAuthor:          "ligen131",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: true,
			FlagCanListen:         false,
			FlagUseDatabase:       false,
			FlagIgnoreRiskControl: false,
		},
	}
	controller.PluginRegister(p)
}

func getCodeforcesContestList(msg string) string {
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
	for i := ListLen - 1; i >= 0; i-- {
		if contestList[i].Before() {
			tot++
			dur, _ := time.ParseDuration(fmt.Sprintf("%ds", contestList[i].DurationSeconds))
			st := time.Unix(contestList[i].StartTimeSeconds, 0).Format("2006-01-02 15:04:05") // Do not change this time
			text += fmt.Sprintf("\n%d. %v, %v, %v", tot, parseCodeforcesContestName(contestList[i].Name), st, dur)
		}
	}
	if tot == 0 {
		text += "\nNone"
	}
	return text
}
func parseCodeforcesContestName(contest string) string {
	ans := contest[strings.Index(contest, "("):]
	num := 0
	ind := 0
	len := len(contest)
	if strings.Contains(contest, "Educational") {
		ans = "Edu "
		ind = strings.Index(contest, "Round")
		if ind != -1 {
			ind += 6
			for ind < len && contest[ind] >= '0' && contest[ind] <= '9' {
				num = num*10 + (int)(contest[ind]) - '0'
				ind++
			}
			ans += fmt.Sprintf("#%d ", num)
		}
	} else if ind = strings.Index(contest, "#"); ind != -1 {
		ind++
		for ind < len && contest[ind] >= '0' && contest[ind] <= '9' {
			num = num*10 + (int)(contest[ind]) - '0'
			ind++
		}
		ans = fmt.Sprintf("#%d ", num)
	}
	if strings.Contains(contest, "Div. 1 + Div. 2") {
		ans += "(Div. 1 + Div. 2)"
	} else if ind = strings.Index(contest, "Div. "); ind != -1 {
		ind += 5
		num = 0
		for ind < len && contest[ind] >= '0' && contest[ind] <= '9' {
			num = num*10 + (int)(contest[ind]) - '0'
			ind++
		}
		ans += fmt.Sprintf("(Div. %d)", num)
	}
	return ans
}
