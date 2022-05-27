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
	QueryLen := -1
	if ind := strings.Index(msg, "-l"); ind != -1 {
		QueryLen = getNumber(msg, ind+2)
		if QueryLen > 100 {
			text = "Query list is too long.\n\n"
			QueryLen = 100
		} else {
			text = ""
		}
		text += fmt.Sprintf("Codeforces Recent %d Contests:", QueryLen)
	}
	for i := ListLen - 1; i >= 0; i-- {
		if (contestList[i].Before() && QueryLen < 0) || (QueryLen > 0 && i+1 <= QueryLen) {
			tot++
			dur := ParseTime(contestList[i].DurationSeconds)
			st := time.Unix(contestList[i].StartTimeSeconds, 0).Format("2006-01-02 15:04:05")[2:16] // Do not change this time
			text += fmt.Sprintf("\n%d. %v, %v, %v", tot, parseCodeforcesContestName(contestList[i].Name), st, dur)
		}
	}
	if tot == 0 {
		text += "\nNone"
	}
	fmt.Println(text)
	return text
}
func parseCodeforcesContestName(contest string) string {
	ind := 0
	ans := ""
	if ind = strings.Index(contest, "("); ind != -1 {
		ans = contest[:ind]
	} else {
		ans = contest
	}
	if strings.Contains(contest, "Educational") {
		ans = "Edu "
		ind = strings.Index(contest, "Round")
		if ind != -1 {
			ans += fmt.Sprintf("#%d ", getNumber(contest, ind+6))
		}
	} else if ind = strings.Index(contest, "#"); ind != -1 {
		ans = fmt.Sprintf("#%d ", getNumber(contest, ind+1))
	}
	if strings.Contains(contest, "Div. 1 + Div. 2") {
		ans += "(Div. 1 + Div. 2)"
	} else if ind = strings.Index(contest, "Div. "); ind != -1 {
		ans += fmt.Sprintf("(Div. %d)", getNumber(contest, ind+5))
	}
	return ans
}

func getNumber(s string, st int) int {
	ans := 0
	len := len(s)
	for st < len && (s[st] < '0' || s[st] > '9') {
		st++
	}
	for st < len && (s[st] >= '0' && s[st] <= '9') {
		ans = ans*10 + int(s[st]) - '0'
		st++
	}
	return ans
}

func ParseTime(second int64) string {
	if second <= 0 {
		return "0s"
	}
	ans := ""
	d := second / (60 * 60 * 24)
	second -= d * 60 * 60 * 24
	h := second / (60 * 60)
	second -= h * 60 * 60
	m := second / 60
	second -= m * 60
	s := second
	if d > 0 {
		ans += fmt.Sprintf("%dd", d)
	}
	if h > 0 {
		ans += fmt.Sprintf("%dh", h)
	}
	if m > 0 {
		ans += fmt.Sprintf("%dmin", m)
	}
	if s > 0 {
		ans += fmt.Sprintf("%ds", s)
	}
	return ans
}
