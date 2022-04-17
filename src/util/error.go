package util

import "github.com/sirupsen/logrus"

func ErrorPrint(err error, data interface{}, info string) {
	logrus.WithFields(logrus.Fields{"err": err.Error(), "data": data}).Info(info)
}
