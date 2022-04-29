package util

import (
	"bug-carrot/param"
	"github.com/yanyiwu/gojieba"
	"strings"
	"sync"
)

var jiebaMutex sync.Mutex

func GetWordsFromMessage(message string) []param.WordSplit {
	jiebaMutex.Lock()
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

	jiebaMutex.Unlock()
	return wordsResponse
}

func GetWordsMapFromMessage(message string) map[param.WordSplit]bool {
	jiebaMutex.Lock()
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

	jiebaMutex.Unlock()
	return wordsMap
}
