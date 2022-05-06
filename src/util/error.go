package util

import "github.com/sirupsen/logrus"

// ErrorPrint 在本地输出错误信息的日志
func ErrorPrint(err error, data interface{}, info string) {
	logrus.WithFields(logrus.Fields{"err": err.Error(), "data": data}).Info(info)
}
