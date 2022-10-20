package util

import (
	"bug-carrot/config"
	"bug-carrot/param"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestQQSend(t *testing.T) {
	QQSend(2039799616, "hello")
}

func TestEmoji(t *testing.T) {
	//message := "😁😂😃😄👿😉😊😌😍😏😒😓😔😖😘😚😜😝😞😠😡😢😣😥😨😪😭😰😱😲😳😷🙃😋😗😛🤑🤓😎🤗🙄🤔😩😤🤐🤒😴😀😆😅😇🙂😙😟😕🙁😫😶😐😑😯😦😧😮😵😬🤕😈👻\U0001F97A\U0001F974🤣\U0001F970🤩🤤🤫🤪🧐🤬🤧🤭🤠🤯🤥\U0001F973🤨🤢🤡🤮\U0001F975\U0001F976💩💀👽👾👺👹🤖😺😸😹😻😼😽🙀😿😾"
	message := "🙈🙉🙊💘💔💯💤"
	messageRune := []rune(message)
	test := ""
	for i := range messageRune {
		test = fmt.Sprintf("%s%d%s", test, i, string(messageRune[i]))
	}
	//for i := 128530; i <= 128563; i++ {
	//	message = fmt.Sprintf("%s%d[CQ:face,id=%d]", message, i, i)
	//	//QQGroupSend(1028801782, fmt.Sprintf("%d[CQ:face,id=%d]", i, i))
	//}
	QQGroupSend(1028801782, getMessageChaosVersion("作业来咯~你可接稳啦!\n【数电】(1)学习通\n【离散】(1)也是学习通\n【大物】(1)71-72\n【电路】(1)1-15,18,26,27,30\n【复变】(1)到练习十二"))
}

func TestSendSameMessageToManyFriends(t *testing.T) {
	mus := param.Muster{
		Title: "123",
	}
	mus.People = append(mus.People, param.PersonWithQQ{
		Name: "name",
		QQ:   config.C.Plugin.Homework.Admin,
	})

	mus.People = append(mus.People, param.PersonWithQQ{
		Name: "name",
		QQ:   1437342516,
	})

	SendSameMessageToManyFriends("作业来咯~你可接稳啦!\n【数电】(1)学习通\n【离散】(1)也是学习通\n【大物】(1)71-72\n【电路】(1)1-15,18,26,27,30\n【复变】(1)到练习十二",
		mus.People)
}

func TestRandomEmoji(t *testing.T) {
	QQGroupSend(1028801782, GetRandomEmojiCQString())
}

func TestParseStruct(t *testing.T) {
	member := param.FamilyMember{
		StudentID: "U00000000",
		Name:      "amy",
		QQ:        123456,
		Phone:     "13054059182",
		Mail:      "amy@qq.com",
		Birthday:  time.Now(),
	}
	v := reflect.ValueOf(member)
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("bson")
		fmt.Println(name, v.Field(i), v.Field(i).Type(), v.Field(i).Kind())

		if v.Field(i).Kind() == reflect.String {
			fmt.Println("!")
		}
	}
}

func TestGetRankString(t *testing.T) {
	ans := GetRankString("518528")
	fmt.Println(ans)
}
