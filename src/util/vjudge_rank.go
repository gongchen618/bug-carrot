package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type probInfo struct {
	accepted bool
	last     int
	times    int
}

func (info *probInfo) getPenalty() int {
	return info.last + info.times*20
}

type responseRankList struct {
	Title        string `json:title`
	Participants map[int][]string
	Submissions  [][]int
}

type rankUser struct {
	Name    string
	ACList  []int
	Penalty int
}

type RankUserList []*rankUser

func (a RankUserList) cmp(i, j int) bool {
	if len(a[i].ACList) != len(a[j].ACList) {
		return len(a[i].ACList) > len(a[j].ACList)
	} else {
		return a[i].Penalty < a[j].Penalty
	}
}

func calcTime(val int) string {
	sec := val % 60
	val /= 60
	minu := val % 60
	val /= 60
	hou := val % 24
	return fmt.Sprintf("%.2d:%.2d:%.2d", hou, minu, sec)
}

func GetRankString(rankListID string) string {
	url := "https://vjudge.net/contest/rank/single/" + rankListID
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	client := &http.Client{}
	resp, _ := client.Do(req)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	body, _ := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	var responseRankList responseRankList
	json.Unmarshal(body, &responseRankList)

	counter := make(map[int]map[int]*probInfo)
	for i := 0; i < len(responseRankList.Submissions); i++ {
		var arr = responseRankList.Submissions[i]
		if counter[arr[0]] == nil {
			counter[arr[0]] = make(map[int]*probInfo)
		}
		if counter[arr[0]][arr[1]] == nil {
			counter[arr[0]][arr[1]] = &probInfo{}
		}
		if counter[arr[0]][arr[1]].accepted {
			continue
		}
		if arr[2] == 0 {
			counter[arr[0]][arr[1]].times++
			if counter[arr[0]][arr[1]].last < arr[3] {
				counter[arr[0]][arr[1]].last = arr[3]
			}
		} else {
			counter[arr[0]][arr[1]].accepted = true
			if counter[arr[0]][arr[1]].last < arr[3] {
				counter[arr[0]][arr[1]].last = arr[3]
			}
		}
	}

	var resList RankUserList
	for key, val := range counter {
		nowUser := rankUser{
			Name: responseRankList.Participants[key][1],
		}
		for id, prob := range val {
			if prob.accepted {
				nowUser.ACList = append(nowUser.ACList, id)
				nowUser.Penalty += prob.getPenalty()
			}
		}
		resList = append(resList, &nowUser)
	}
	sort.Slice(resList, resList.cmp)

	var res string
	res = ""
	for rank, val := range resList {
		var prob string
		for _, id := range val.ACList {
			prob = prob + " " + strconv.Itoa(id+1)
		}
		res = res + "第 " + strconv.Itoa(rank+1) + " 名: " + val.Name + " 通过了题目: " + prob + ", 罚时: " + calcTime(val.Penalty) + "\n"
	}

	return res

}
