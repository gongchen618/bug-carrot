package util

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
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

func GetMessageChaosVersion(message string) string {
	hitokotoResp := getHitokotoSentence()
	messageRune := []rune(message)
	messageLen := len(messageRune)
	for i := 0; i <= messageLen-1 && i <= messageLen/9; i++ {
		a, err := rand.Int(rand.Reader, big.NewInt(int64(messageLen-1)))
		if err != nil || unicode.IsNumber(messageRune[a.Int64()]) || unicode.IsNumber(messageRune[a.Int64()+1]) {
			a = big.NewInt(0)
		}
		messageRune[a.Int64()], messageRune[a.Int64()+1] = messageRune[a.Int64()+1], messageRune[a.Int64()]
	}
	return fmt.Sprintf("「%s」\n%s\nfrom.%s", hitokotoResp.Hitokoto, string(messageRune), hitokotoResp.From)
}
