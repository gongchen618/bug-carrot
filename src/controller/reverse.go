package controller

// qqbot 消息的接收层
// 注意一般情况下不要给 qqbot 返回 error

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller/param"
	"bug-carrot/util"
	"bug-carrot/util/context"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strings"
)

func QQReverseHTTPMiddleHandler(c echo.Context) error {
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	p := struct {
		PostType string `json:"post_type"`
	}{}
	if err := json.Unmarshal(bodyBytes, &p); err != nil {
		util.ErrorPrint(err, nil, "first unmarshal failed")
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	switch p.PostType {
	case "message":
		p2 := struct {
			MessageType string `json:"message_type"`
		}{}
		if err := json.Unmarshal(bodyBytes, &p2); err != nil {
			util.ErrorPrint(err, nil, "second unmarshal failed")
			return context.Error(c, http.StatusBadRequest, "bad request", err)
		}
		switch p2.MessageType {
		case "private":
			return privateMessageHandler(c)
		case "group":
			return groupMessageHandler(c)
		}
	case "request":
		p2 := struct {
			RequestType string `json:"request_type"`
		}{}
		if err := json.Unmarshal(bodyBytes, &p2); err != nil {
			util.ErrorPrint(err, nil, "second unmarshal failed")
			return context.Error(c, http.StatusBadRequest, "bad request", err)
		}
		switch p2.RequestType {
		//case "friend":
		//return friendAddRequestHandler(c)
		}
	}

	return context.Success(c, nil)
}

//func friendAddRequestHandler(c echo.Context) error {
//	req := param.RequestQQFriendAdd{}
//	if err := c.Bind(&req); err != nil {
//		logrus.WithFields(logrus.Fields{"err": err.Error()}).Info("bind failed")
//		return context.Error(c, http.StatusBadRequest, "bad request", err)
//	}
//
//	// approve friend add request
//	err := util.QQApproveFriendAddRequest(req.Flag)
//	if err != nil {
//		logrus.WithFields(logrus.Fields{"err": err.Error()}).Info("approve friend add request failed")
//		return context.Error(c, http.StatusInternalServerError, "can't approve", err)
//	}
//
//	// send hello message
//	err = util.QQSend(req.UserId, BotHelloString)
//	if err != nil {
//		logrus.WithFields(logrus.Fields{"err": err.Error()}).Info("send message failed")
//		return context.Error(c, http.StatusInternalServerError, "can't send message", err)
//	}
//	return context.Success(c, nil)
//}

func privateMessageHandler(c echo.Context) error {
	req := param.RequestPrivateMessage{}
	if err := c.Bind(&req); err != nil {
		util.ErrorPrint(err, nil, "bind failed")
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	if req.UserId != config.C.QQ.Admin {
		if req.SubType == "friend" {
			util.QQSend(req.UserId, constant.CarrotFriendNotAdmin)
		}
		return context.Success(c, nil)
	}

	str := strings.Split(req.RawMessage, " ")
	if len(str) >= 2 {
		switch str[1] {
		case "delete":
			if len(str) >= 4 {
				solveAdminHomeworkDeleteMessage(req.UserId, str[2], str[3])
				return context.Success(c, nil)
			}
		case "show":
			solveAdminHomeworkShowMessage(req.UserId)
			return context.Success(c, nil)
		case "add":
			if len(str) >= 4 {
				solveAdminHomeworkAddMessage(req.UserId, str[2], str[3])
				return context.Success(c, nil)
			}
		}
	}

	return context.Success(c, nil)
}

func groupMessageHandler(c echo.Context) error {
	req := param.RequestGroupMessage{}
	if err := c.Bind(&req); err != nil {
		util.ErrorPrint(err, nil, "bind failed")
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	userId := req.UserId
	if req.SubType == "anonymous" {
		userId = req.Anonymous.Id
	}

	recordGroupMessage(req.GroupId, userId, req.RawMessage) // 记录x

	prefixAt := fmt.Sprintf("[CQ:at,qq=%d]", config.C.QQBot.Bot)
	if !strings.HasPrefix(req.RawMessage, prefixAt) {
		return context.Success(c, http.StatusOK)
	}

	words := util.GetWordsFromString(req.RawMessage)
	for _, word := range words { // 强功能性关键词
		switch word.Type {
		case "n":
			switch word.Word {
			case "作业":
				dealGroupHomeworkMessage(req.GroupId, userId, words)
				return context.Success(c, nil)
			case "天气":
				dealGroupWeatherMessage(req.GroupId, userId, words)
				return context.Success(c, nil)
			}
		case "eng":
			switch word.Word {
			case "cf", "codeforces":
				delGroupCodeforcesMessage(req.GroupId, userId, words)
				return context.Success(c, nil)
			}
		}
	}

	dealGroupPartyMessage(req.GroupId, userId, words) // 是要和卡洛一起玩吗?
	return context.Success(c, nil)
}
