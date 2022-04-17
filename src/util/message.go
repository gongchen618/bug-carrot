package util

import (
	"bug-carrot/config"
	"bug-carrot/controller/param"
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

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

func GetQQGroupUserId(msg param.GroupMessage) int64 {
	userId := msg.UserId
	if msg.SubType == "anonymous" {
		userId = msg.Anonymous.Id
	}
	return userId
}

func packageMessage(message string) string {
	data := rand.Int63n(222)
	return fmt.Sprintf("%s[CQ:face,id=%d]", message, data)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
