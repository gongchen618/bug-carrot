package plugin

import (
	"bug-carrot/src/constant"
	"bug-carrot/src/controller"
	"bug-carrot/src/model"
	"bug-carrot/src/param"
	"bug-carrot/src/util"
	"fmt"
	"strings"
)

type keyWord struct {
	Index               param.PluginIndex
	KeyWordAddPrefix    string
	KeyWordDeletePrefix string
	KeyWordQueryPrefix  string
	KeyWordCheckPrefix  string
	DividingString      string
}

func (p *keyWord) GetPluginName() string {
	return p.Index.PluginName
}
func (p *keyWord) GetPluginAuthor() string {
	return p.Index.PluginAuthor
}
func (p *keyWord) CanTime() bool {
	return p.Index.FlagCanTime
}
func (p *keyWord) CanMatchedGroup() bool {
	return p.Index.FlagCanMatchedGroup
}
func (p *keyWord) CanMatchedPrivate() bool {
	return p.Index.FlagCanMatchedPrivate
}
func (p *keyWord) CanListen() bool {
	return p.Index.FlagCanListen
}
func (p *keyWord) NeedDatabase() bool {
	return p.Index.FlagUseDatabase
}
func (p *keyWord) DoIgnoreRiskControl() bool {
	return p.Index.FlagIgnoreRiskControl
}

func (p *keyWord) IsTime() bool {
	return false
}
func (p *keyWord) DoTime() error {
	return nil
}

func (p *keyWord) IsMatchedGroup(msg param.GroupMessage) bool {
	if msg.Anonymous.Id != 0 { // 禁止匿名
		return false
	}
	if strings.HasPrefix(msg.RawMessage, p.KeyWordAddPrefix) ||
		strings.HasPrefix(msg.RawMessage, p.KeyWordDeletePrefix) ||
		strings.HasPrefix(msg.RawMessage, p.KeyWordQueryPrefix) ||
		strings.HasPrefix(msg.RawMessage, p.KeyWordCheckPrefix) {
		return true
	}
	return false
}
func (p *keyWord) DoMatchedGroup(msg param.GroupMessage) error {
	switch {
	case strings.HasPrefix(msg.RawMessage, p.KeyWordAddPrefix):
		info := msg.RawMessage[len(p.KeyWordAddPrefix):]
		if strings.Count(info, p.DividingString) == 0 {
			util.QQGroupSendAtSomeone(msg.GroupId, msg.UserId, constant.CarrotFoodStrangeInput)
			return nil
		}
		str := strings.Split(info, p.DividingString)
		keyWordAdd(msg.UserId, msg.GroupId, str[0], strings.Join(str[1:], p.DividingString))

	case strings.HasPrefix(msg.RawMessage, p.KeyWordDeletePrefix):
		name := msg.RawMessage[len(p.KeyWordDeletePrefix):]
		keyWordDelete(msg.GroupId, name)

	case strings.HasPrefix(msg.RawMessage, p.KeyWordQueryPrefix):
		kw := msg.RawMessage[len(p.KeyWordQueryPrefix):]
		keyWordQuery(msg.GroupId, kw)

	case strings.HasPrefix(msg.RawMessage, p.KeyWordCheckPrefix):
		kw := msg.RawMessage[len(p.KeyWordCheckPrefix):]
		keyWordCheck(msg.GroupId, kw)
	}
	return nil
}

func (p *keyWord) IsMatchedPrivate(msg param.PrivateMessage) bool {
	return false
}
func (p *keyWord) DoMatchedPrivate(msg param.PrivateMessage) error {
	return nil
}

func (p *keyWord) Listen(msg param.GroupMessage) {

}

func (p *keyWord) Close() {
}

func KeyWordPluginRegister() {
	p := &keyWord{
		Index: param.PluginIndex{
			PluginName:            "keyword",
			PluginAuthor:          "gongchen618",
			FlagCanTime:           false,
			FlagCanMatchedGroup:   true,
			FlagCanMatchedPrivate: false,
			FlagCanListen:         false,
			FlagUseDatabase:       true,
			FlagIgnoreRiskControl: false,
		},
		KeyWordAddPrefix:    "立契",
		KeyWordDeletePrefix: "悔契",
		KeyWordQueryPrefix:  "召唤",
		KeyWordCheckPrefix:  "调查",
		DividingString:      "#",
	}
	controller.PluginRegister(p)
}

func keyWordAdd(id int64, group int64, keyWord string, content string) {
	m := model.GetModel()
	defer m.Close()

	fd := param.KeyWord{
		KeyWord: keyWord,
		Content: content,
		Author:  id,
	}

	if err := m.AddKeyWord(fd); err != nil {
		util.QQGroupSend(group, "契约签订失败!")
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQGroupSend(group, fmt.Sprintf("契约「%s」已签订", keyWord))
}

func keyWordDelete(group int64, keyWord string) {
	m := model.GetModel()
	defer m.Close()

	if err := m.DeleteKeyWord(keyWord); err != nil {
		util.QQGroupSend(group, "契约销毁失败!")
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQGroupSend(group, fmt.Sprintf("契约「%s」已销毁", keyWord))
}

func keyWordQuery(group int64, keyWord string) {
	m := model.GetModel()
	defer m.Close()

	kw, err := m.GetKeyWord(keyWord)
	if err != nil {
		util.QQGroupSend(group, "不存在的契约!")
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQGroupSend(group, fmt.Sprintf("%s", kw.Content))
}

func keyWordCheck(group int64, keyWord string) {
	m := model.GetModel()
	defer m.Close()

	kw, err := m.GetKeyWord(keyWord)
	if err != nil {
		util.QQGroupSend(group, "不存在的契约!")
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQGroupSend(group, fmt.Sprintf("契约「%s」的签订者是 %d", keyWord, kw.Author))
}
