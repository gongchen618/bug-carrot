package util

import (
	"bug-carrot/param"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
	"unicode"
)

type hitokotoResponse struct {
	Hitokoto string `json:"hitokoto"`
	From     string `json:"from"`
}

var (
	defaultSentence = hitokotoResponse{
		Hitokoto: "嗨呀",
		From:     "Carrot卡洛塔",
	}
)

func getHitokotoSentence() hitokotoResponse {
	url := "https://v1.hitokoto.cn"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return defaultSentence
	}

	q := req.URL.Query()
	q.Add("c", "a")
	q.Add("c", "b")
	q.Add("c", "c")
	q.Add("encode", "json")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return defaultSentence
	}

	body, _ := ioutil.ReadAll(resp.Body)
	hitokotoResp := hitokotoResponse{}
	if err = json.Unmarshal(body, &hitokotoResp); err != nil {
		return defaultSentence
	}
	return hitokotoResp
}

func GetHitokotoWarpedMessage(message string) string {
	hitokotoResp := getHitokotoSentence()
	return fmt.Sprintf("「%s」\n%s\nfrom.%s", hitokotoResp.Hitokoto, message, hitokotoResp.From)
}

func getMessageChaosVersion(message string) string {
	messageRune := []rune(message)
	messageLen := len(messageRune)
	for i := 0; i <= messageLen-1 && i <= messageLen/9; i++ {
		a, err := rand.Int(rand.Reader, big.NewInt(int64(messageLen-1)))
		wordi, wordii := messageRune[a.Int64()], messageRune[a.Int64()+1]
		if err != nil ||
			unicode.IsNumber(wordi) || unicode.IsNumber(wordii) ||
			unicode.IsLetter(wordi) || unicode.IsLetter(wordii) {
			continue
		}
		messageRune[a.Int64()], messageRune[a.Int64()+1] = wordii, wordi
	}
	return string(messageRune)
}

func getMessageLinkMixedVersion(message string) string {
	messageRune := []rune(message)
	messageLen := len(messageRune)
	messageNew := ""
	for i := 0; i < messageLen; i++ {
		wordi := messageRune[i]
		if wordi == '.' || wordi == ':' {
			messageNew = fmt.Sprintf("%s%s", messageNew, GetRandomEmojiCQString())
			continue
		}
		messageNew = fmt.Sprintf("%s%s", messageNew, string(messageRune[i]))
	}
	return messageNew
}

// SendSameMessageToManyFriends : 批量发送同一条消息，混淆汉字顺序和添加无关内容后不均匀延迟发送
func SendSameMessageToManyFriends(message string, people []param.PersonWithQQ) []param.PersonWithQQ {
	var failed []param.PersonWithQQ
	timer := time.NewTimer(time.Duration(20))
	for _, person := range people {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(20)))
		if err != nil {
			num = big.NewInt(20)
		}
		timer.Reset(time.Second * time.Duration(num.Int64()))
		select {
		case <-timer.C:
			message = getMessageLinkMixedVersion(message)
			message = getMessageChaosVersion(message)
			message = GetHitokotoWarpedMessage(message)

			status := QQSendAndFindWhetherSuccess(person.QQ, message)
			if status == false {
				failed = append(failed, person)
			}
		}
	}
	return failed
}
