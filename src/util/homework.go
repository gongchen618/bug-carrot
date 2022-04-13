package util

import (
	"bug-carrot/constant"
	"bug-carrot/model"
	"fmt"
	"time"
)

func GetHomeworkString() string {
	m := model.GetModel()
	defer m.Close()

	timeR := time.Now()
	d, err := time.ParseDuration("-300h")
	if err != nil {
		ErrorPrint(err, timeR, "time init")
		return constant.CarrotHomeworkShowFailed
	}
	timeL := timeR.Add(d)
	homeworks, err := m.GetHomeworkByTimeRange(timeL, timeR)
	if err != nil {
		ErrorPrint(err, timeR, "mongo")
		return constant.CarrotHomeworkShowFailed
	}

	if len(homeworks) == 0 {
		return constant.CarrotHomeworkShowEmpty
	}

	message := constant.CarrotHomeworkShowStart
	subjectMap := make(map[string]int)
	subjectInfoMap := make(map[string]string)
	for _, homework := range homeworks {
		cnt, exist := subjectMap[homework.Subject]
		info, exist := subjectInfoMap[homework.Subject]
		if exist == false {
			subjectMap[homework.Subject] = 1
			subjectInfoMap[homework.Subject] = fmt.Sprintf("(%d)%s", cnt+1, homework.Context)
		} else {
			subjectMap[homework.Subject] = cnt + 1
			subjectInfoMap[homework.Subject] = fmt.Sprintf("%s (%d)%s", info, cnt+1, homework.Context)
		}
	}
	for subject, info := range subjectInfoMap {
		message = fmt.Sprintf("%s\n【%s】%s", message, subject, info)
	}

	return message
}

func GetHomeworkStringSubject(subject string) string {
	return GetHomeworkString()
}
