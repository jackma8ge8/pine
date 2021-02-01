package logger

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/jackma8ge8/pine/application/config"
	"github.com/sirupsen/logrus"
)

// 字体 背景  颜色
// 31   41   红色
// 32   42   绿色
// 33   43   黄色
// 34   44   蓝色
// 35   45   洋红
// 36   46   青色
// 37   47   白色
const (
	Error = 31
	Warn  = 33
	Debug = 34
	Fatal = 35
	Info  = 36
	Trace = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.ErrorLevel, logrus.PanicLevel:
		return Error
	case logrus.FatalLevel:
		return Fatal
	case logrus.WarnLevel:
		return Warn
	case logrus.DebugLevel:
		return Debug
	case logrus.TraceLevel:
		return Trace
	case logrus.InfoLevel:
		return Info
	default:
		return Info
	}
}

// ErrorFormatter 错误格式化器
type ErrorFormatter struct{}

// Format Format函数
func (f ErrorFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	// output buffer
	b := &bytes.Buffer{}

	levelColor := getColorByLevel(entry.Level)
	fmt.Fprintf(b, "\x1b[%dm", levelColor)

	timestampFormat := time.RFC3339

	// 时间
	b.WriteString(entry.Time.Format(timestampFormat))

	// 日志等级
	b.WriteString(" [")
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString("]")

	// message
	if entry.Message != "" {
		b.WriteString(" [" + config.GetServerConfig().ID + ": " + strings.TrimSpace(entry.Message) + "]")
	} else {
		b.WriteString(" [" + config.GetServerConfig().ID + "]")
	}

	stack, hasStack := entry.Data["Stack"]
	if hasStack {
		delete(entry.Data, "Stack")
	}

	// fields
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprint(" [", key, ":", value, "] "))
	}

	if hasStack {
		// stack
		b.WriteString(fmt.Sprint("\nCall stack:\n", stack))
		b.WriteString("\x1b[0m\n")
	} else {
		// caller
		b.WriteString("\x1b[0m")
		b.WriteString(fmt.Sprint(" ", entry.Caller.File, ":", entry.Caller.Line, "\n"))
	}

	return b.Bytes(), nil

}
