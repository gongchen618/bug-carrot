package util

import (
	"bug-carrot/config"
	"bug-carrot/param"
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// QQApproveFriendAddRequest 通过好友申请
func QQApproveFriendAddRequest(flag string) {
	url := fmt.Sprintf("%s%s", config.C.QQBot.Host, "/set_friend_add_request")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		ErrorPrint(err, flag, "friend add")
		return
	}

	q := req.URL.Query()
	q.Add("flag", flag)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		ErrorPrint(err, flag, "friend add")
	}
}

// QQSend 接受 userId 和 message，并私聊发送
// 如果不是好友，会发送失败
func QQSend(userId int64, message string) {
	sendMsgUrl := fmt.Sprintf("%s%s", config.C.QQBot.Host, "/send_private_msg")
	req, err := http.NewRequest("POST", sendMsgUrl, bytes.NewBuffer(nil))
	if err != nil {
		ErrorPrint(err, userId, "message send")
		return
	}

	q := req.URL.Query()
	q.Add("user_id", strconv.FormatInt(userId, 10))
	q.Add("message", packageMessage(message))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		ErrorPrint(err, userId, "message send")
	}
}

// QQGroupSend 接受 groupId 和 message，并在对应群聊中发送
// 如果不在群聊中，会发送失败
func QQGroupSend(groupId int64, message string) {
	sendMsgUrl := fmt.Sprintf("%s%s", config.C.QQBot.Host, "/send_group_msg")
	req, err := http.NewRequest("POST", sendMsgUrl, bytes.NewBuffer(nil))
	if err != nil {
		ErrorPrint(err, groupId, "group message send")
		return
	}

	q := req.URL.Query()
	q.Add("group_id", strconv.FormatInt(groupId, 10))
	q.Add("message", packageMessage(message))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		ErrorPrint(err, groupId, "group message send")
	}
}

// QQGroupSendAtSomeone 接受 groupId 和 userId 和 message，并在对应群聊中 @ 对应成员并发送
// 如果不在群聊中，会发送失败
func QQGroupSendAtSomeone(groupId int64, userId int64, message string) {
	sendMsgUrl := fmt.Sprintf("%s%s", config.C.QQBot.Host, "/send_group_msg")
	req, err := http.NewRequest("POST", sendMsgUrl, bytes.NewBuffer(nil))
	if err != nil {
		ErrorPrint(err, groupId, "group message send")
		return
	}

	q := req.URL.Query()
	q.Add("group_id", strconv.FormatInt(groupId, 10))
	q.Add("message", fmt.Sprintf("[CQ:at,qq=%d] %s", userId, packageMessage(message)))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		ErrorPrint(err, groupId, "group message send")
	}
}

// QQGroupBan 接受 groupId 和 userId，并在对应群聊中禁言用户，时长 cnt (单位分钟）
// 如果不是管理员，会禁言失败
func QQGroupBan(groupId int64, userId int64, cnt int64) {
	sendMsgUrl := fmt.Sprintf("%s%s", config.C.QQBot.Host, "/set_group_ban")
	req, err := http.NewRequest("POST", sendMsgUrl, bytes.NewBuffer(nil))
	if err != nil {
		ErrorPrint(err, groupId, "group ban set")
		return
	}

	q := req.URL.Query()
	q.Add("group_id", strconv.FormatInt(groupId, 10))
	q.Add("user_id", strconv.FormatInt(userId, 10))
	q.Add("duration", strconv.FormatInt((cnt*2-1)*60, 10))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		ErrorPrint(err, groupId, "group ban set")
	}
}

// GetQQGroupUserId 获取 msg 中的 userId (处理匿名用户的问题）
func GetQQGroupUserId(msg param.GroupMessage) int64 {
	userId := msg.UserId
	if msg.SubType == "anonymous" {
		userId = msg.Anonymous.Id
	}
	return userId
}

// packageMessage 在消息后面增加一个随机表情
func packageMessage(message string) string {
	data := rand.Int63n(222)
	return fmt.Sprintf("%s[CQ:face,id=%d]", message, data)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
