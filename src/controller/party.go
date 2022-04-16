package controller

import (
	"bug-carrot/constant"
	"bug-carrot/controller/param"
	"bug-carrot/util"
	"fmt"
	"time"
)

var (
	goodMorningDay     = time.Now().Day() - 1
	goodMorningUserDay = make(map[int64]int)
	goodNightDay       = time.Now().Day() - 1
	goodNightUserDay   = make(map[int64]int)
)

func dealGroupPartyMessage(groupId int64, userId int64, words []param.WordSplit) {
	for _, word := range words {
		switch word.Type {
		case "n":
			switch word.Word {
			case "早安": // 早安
				now := time.Now()
				fmt.Println(now.Hour())
				if now.Hour() > 12 || now.Hour() < 5 {
					util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodMorningCheat)
					return
				}

				userDay, flag := goodMorningUserDay[userId]
				if flag && userDay == now.Day() {
					util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodMorningRepeat)
					return
				}
				goodMorningUserDay[userId] = now.Day()
				if goodMorningDay != now.Day() && now.Hour() > 5 {
					goodMorningDay = now.Day()
					util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodMorningFirst)
				} else {
					util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodMorning)
				}
				return

			case "晚安": // 晚安
				now := time.Now()
				if now.Hour() > 5 && now.Hour() < 20 {
					util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodNightCheat)
					return
				}

				userDay, flag := goodNightUserDay[userId]
				if flag && userDay == now.Day() {
					util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodNightRepeat)
					return
				}
				goodNightUserDay[userId] = now.Day()
				if goodNightDay != now.Day() && now.Hour() > 5 {
					goodNightDay = now.Day()
					util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodNightFirst)
				} else {
					util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodNight)
				}
				return

			case "萝卜": // 萝卜
				util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupCarrot)
				return
			}
		}
	}
	util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupPuzzled)
}

var (
	messageRecorded   [10]string
	messageRepeatLast string
	goodNightCheck    = make(map[int64]int64)
)

func recordGroupMessage(groupId int64, userId int64, message string) {
	flagSame := true
	for i := 0; i < 3; i++ {
		if messageRecorded[i] != message {
			flagSame = false
		}
	}

	if flagSame && messageRepeatLast != message {
		util.QQGroupSend(groupId, message)
		messageRepeatLast = message
	}

	for i := len(messageRecorded) - 1; i > 0; i-- {
		messageRecorded[i] = messageRecorded[i-1]
	}
	messageRecorded[0] = message

	//check good night
	//userDay, flag := goodNightUserDay[userId]
	//if flag && userDay == time.Now().Day() {
	//	userTim, flag := goodNightCheck[userId]
	//	if !flag {
	//		userTim = 0
	//	}
	//	userTim = userTim + 1
	//	if userTim%7 != 0 {
	//		if userTim%3 == 0 {
	//			util.QQGroupSend(groupId, fmt.Sprintf("[CQ:poke,qq=%d]", userId))
	//		}
	//	} else {
	//		util.QQGroupBan(groupId, userId, userTim/7)
	//		util.QQGroupSendAtSomeone(groupId, userId, constant.CarrotGroupGoodNightButChat)
	//	}
	//	goodNightCheck[userId] = userTim
	//}
}
