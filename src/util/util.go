package util

import (
	"bug-carrot/controller/param"
	"github.com/yanyiwu/gojieba"
	"strings"
)

func GetWordsFromString(message string) []param.WordSplit {
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
