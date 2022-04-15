package controller

import (
	"bug-carrot/config"
	"bug-carrot/constant"
	"bug-carrot/controller/param"
	"bug-carrot/model"
	"bug-carrot/util"
	"fmt"
)

func solveAdminHomeworkDeleteMessage(id int64, subject string, context string) {
	m := model.GetModel()
	defer m.Close()

	homework := param.Homework{
		Subject: subject,
		Context: context,
	}

	if err := m.DeleteHomework(homework); err != nil {
		util.QQSend(id, constant.CarrotHomeworkDeleteFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotHomeworkDeleteSuccess)
	util.QQGroupSend(config.C.QQ.Group, fmt.Sprintf("删除作业 %s %s", subject, context))
}

func solveAdminHomeworkShowMessage(id int64) {
	util.QQSend(id, util.GetHomeworkString())
}

func solveAdminHomeworkAddMessage(id int64, subject string, context string) {
	m := model.GetModel()
	defer m.Close()

	homework := param.Homework{
		Subject: subject,
		Context: context,
	}

	if err := m.AddHomework(homework); err != nil {
		util.QQSend(id, constant.CarrotHomeworkAddFailed)
		util.ErrorPrint(err, nil, "mongo")
		return
	}

	util.QQSend(id, constant.CarrotHomeworkAddSuccess)
	util.QQGroupSend(config.C.QQ.Group, fmt.Sprintf("新增作业 %s %s", subject, context))
}
