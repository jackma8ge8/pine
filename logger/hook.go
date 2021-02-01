package logger

import (
	"github.com/facebookgo/stack"
	"github.com/sirupsen/logrus"
)

// error hook
type errorHook struct{}

func (h *errorHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel}
}

func (h *errorHook) Fire(entry *logrus.Entry) error {

	var frames stack.Stack
	if len(entry.Data) == 0 {
		frames = stack.Callers(8)
	} else {
		frames = stack.Callers(6)
	}

	if LogType == LogTypeEnum.Console {
		for index, frame := range frames {
			frames[index].File = "    " + frame.File
		}
	}

	entry.Data["Stack"] = frames

	return nil
}
