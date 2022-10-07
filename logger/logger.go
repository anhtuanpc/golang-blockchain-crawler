package logger

import (
	"cbridgewrapper/config"
	"io/ioutil"

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true
	Logger.SetFormatter(customFormatter)
	Logger.SetNoLock()
	Logger.SetOutput(&lumberjack.Logger{
		Filename:   config.GetConfig("log_path"),
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     2, //days
	})
	logLevel := config.GetConfig("log.level")

	if logLevel == "" {
		Logger.Out = ioutil.Discard
	} else {
		if logLevel == "DEBUG" {
			Logger.SetLevel(logrus.DebugLevel)
		} else if logLevel == "INFO" {
			Logger.SetLevel(logrus.InfoLevel)
		} else if logLevel == "WARN" {
			Logger.SetLevel(logrus.WarnLevel)
		} else if logLevel == "ERROR" {
			Logger.SetLevel(logrus.ErrorLevel)
		}
	}
}
