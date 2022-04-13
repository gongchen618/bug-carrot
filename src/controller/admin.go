package controller

import (
	"bug-carrot/constant"
	"bug-carrot/controller/param"
	"bug-carrot/model"
	"bug-carrot/util"
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
}
