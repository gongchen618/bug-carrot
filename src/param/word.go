package param

type WordSplit struct {
	Type string
	Word string
}

type WordsMap struct {
	Map map[WordSplit]bool
}

func (wm WordsMap) ExistWord(t string, ws []string) bool {
	for _, w := range ws {
		ans, flag := wm.Map[WordSplit{
			Type: t,
			Word: w,
		}]
		if flag && ans == true {
			return true
		}
	}
	return false
}
