package util

import (
	"bug-carrot/controller/param"
	"github.com/yanyiwu/gojieba"
	"strings"
)

func GetWordsFromMessage(message string) []param.WordSplit {
	x := gojieba.NewJieba()
	defer x.Free()

	words := x.Tag(message)
	var wordsResponse []param.WordSplit
	for _, word := range words {
		wordsSplit := strings.Split(word, "/")
		wordsResponse = append(wordsResponse, param.WordSplit{
			Type: wordsSplit[1],
			Word: wordsSplit[0],
		})
	}

	return wordsResponse
}

func GetWordsMapFromMessage(message string) map[param.WordSplit]bool {
	x := gojieba.NewJieba()
	defer x.Free()

	words := x.Tag(message)
	wordsMap := make(map[param.WordSplit]bool)
	for _, word := range words {
		wordsSplit := strings.Split(word, "/")
		wordsMap[param.WordSplit{
			Type: wordsSplit[1],
			Word: wordsSplit[0],
		}] = true
	}

	return wordsMap
}

func IsWordInMessage(t string, ws []string, msg param.GroupMessage) bool {
	for _, w := range ws {
		ans, flag := msg.WordsMap[param.WordSplit{
			Type: t,
			Word: w,
		}]
		if flag && ans == true {
			return true
		}
	}
	return false
}
