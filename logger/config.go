package logger

import (
	"log"
	"os"
	"time"

	"github.com/jackma8ge8/pine/application/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// LogType 当前日志类型
var LogType string

// SetLogMode 设置log模式
func SetLogMode(logType string) {
	LogType = logType

	logrus.AddHook(&errorHook{})
	logrus.SetReportCaller(true)

	logLevel := logrus.DebugLevel // Default
	switch config.GetServerConfig().LogLevel {
	case LogLevelEnum.Debug: // Debug
		logLevel = logrus.DebugLevel
	case LogLevelEnum.Info: // Info
		logLevel = logrus.InfoLevel
	case LogLevelEnum.Warn: // Warn
		logLevel = logrus.WarnLevel
	case LogLevelEnum.Error: // Error
		logLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logLevel)
	if LogType == LogTypeEnum.File {
		logrus.SetFormatter(&logrus.JSONFormatter{})

		path, _ := os.Getwd()

		writer, err := rotatelogs.New(
			path+"/log/"+config.GetServerConfig().ID+"-%m_%d-%H_%M.log",
			rotatelogs.WithMaxAge(time.Hour*24*30),    // 保留时间
			rotatelogs.WithRotationTime(24*time.Hour), // 分割间隔
		)

		if err != nil {
			log.Fatalf("create file log.txt failed: %v", err)
		}

		logrus.SetOutput(writer)
	} else {
		logrus.SetFormatter(ErrorFormatter{})
	}
}
