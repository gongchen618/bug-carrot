package controller

// qqbot 消息的接收层
// 注意一般情况下不要给 qqbot 返回 error

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/param"
	"bug-carrot/util"
	"bug-carrot/util/context"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// QQReverseHTTPMiddleHandler 接受 cqhttp 的消息
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
		case "friend":
			return friendAddRequestHandler(c)
		}
	}

	return context.Success(c, nil)
}

func friendAddRequestHandler(c echo.Context) error {
	req := param.RequestQQFriendAdd{}
	if err := c.Bind(&req); err != nil {
		util.ErrorPrint(err, nil, "bind failed")
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	// approve friend add request
	util.QQApproveFriendAddRequest(req.Flag)
	// send hello message
	time.Sleep(2 * time.Second)
	util.QQSend(req.UserId, constant.CarrotFriendAddHello)
	return context.Success(c, nil)
}

func privateMessageHandler(c echo.Context) error {
	req := param.RequestPrivateMessage{}
	if err := c.Bind(&req); err != nil {
		util.ErrorPrint(err, nil, "bind failed")
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	r := param.PrivateMessage{
		RequestPrivateMessage: req,
		WordsMap: param.WordsMap{
			Map: util.GetWordsMapFromMessage(req.RawMessage),
		},
	}
	go func(msg param.PrivateMessage) {
		WorkPrivateMessagePlugins(msg)
	}(r)

	return context.Success(c, nil)
}

func groupMessageHandler(c echo.Context) error {
	req := param.RequestGroupMessage{}
	if err := c.Bind(&req); err != nil {
		util.ErrorPrint(err, nil, "bind failed")
		return context.Error(c, http.StatusBadRequest, "bad request", err)
	}

	prefA := fmt.Sprintf("[CQ:at,qq=%d]", config.C.QQBot.QQ)
	prefT := fmt.Sprintf("@%s", config.C.QQBot.Name)
	if strings.HasPrefix(req.RawMessage, prefA) || strings.HasPrefix(req.RawMessage, prefT) {
		var message string
		if strings.HasPrefix(req.RawMessage, prefA) {
			message = req.RawMessage[len(prefA):]
		} else {
			message = req.RawMessage[len(prefT):]
		}
		message = strings.TrimLeft(message, " ")
		r := param.GroupMessage{
			RequestGroupMessage: param.RequestGroupMessage{
				SubType:    req.SubType,
				RawMessage: message,
				UserId:     req.UserId,
				GroupId:    req.GroupId,
				Anonymous:  req.Anonymous,
				Sender:     req.Sender,
			},
			WordsMap: param.WordsMap{
				Map: util.GetWordsMapFromMessage(req.RawMessage),
			},
		}
		go func(msg param.GroupMessage) {
			WorkGroupMessagePlugins(msg)
		}(r)
	} else {
		r := param.GroupMessage{
			RequestGroupMessage: req,
			WordsMap: param.WordsMap{
				Map: util.GetWordsMapFromMessage(req.RawMessage),
			},
		}
		go func(msg param.GroupMessage) {
			WorkListenPlugins(msg)
		}(r)
	}
	return context.Success(c, nil)
}
