package controller

import (
	"bug-carrot/controller/param"
	"bug-carrot/util"
)

func dealGroupHomeworkMessage(groupId int64, userId int64, words []param.WordSplit) {
	for _, word := range words {
		switch word.Type {
		case "n":
			switch word.Word {
			case "大物", "大学物理", "物理", "大雾":
				util.QQGroupSendAtSomeone(groupId, userId, util.GetHomeworkStringSubject("大物"))
				return
			case "微积分":
				util.QQGroupSendAtSomeone(groupId, userId, util.GetHomeworkStringSubject("微积分"))
				return
			case "离散", "离散数学":
				util.QQGroupSendAtSomeone(groupId, userId, util.GetHomeworkStringSubject("离散"))
				return
			}
		}
	}

	util.QQGroupSendAtSomeone(groupId, userId, util.GetHomeworkString())
}

func dealGroupWeatherMessage(groupId int64, userId int64, words []param.WordSplit) {
	location := "武汉"
	for _, word := range words {
		if word.Type == "ns" {
			location = word.Word
		}
	}
	util.QQGroupSendAtSomeone(groupId, userId, util.GetWeatherInfoString(location))
}

func delGroupCodeforcesMessage(groupId int64, userId int64, words []param.WordSplit) {

}
