package util

import "github.com/sirupsen/logrus"

func ErrorReturnAndPrint(err error, data interface{}, info string) error {
	logrus.WithFields(logrus.Fields{"err": err.Error(), "data": data}).Info(info)
	return err
}

func ErrorPrint(err error, data interface{}, info string) {
	logrus.WithFields(logrus.Fields{"err": err.Error(), "data": data}).Info(info)
}
