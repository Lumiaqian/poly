package log

import (
	"context"

	"github.com/Lumiaqian/go-sdk-core/log"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/lib/logger"
)

type LogAdapter struct {
	log *logger.CustomLogger
}

func NewLogAdapter(log *logger.CustomLogger) *LogAdapter {
	logger.GlobalLogger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	return &LogAdapter{log: log}
}

func (a *LogAdapter) Log(ctx context.Context, level log.Level, keyvals ...interface{}) {
	var (
		fields logger.Fields = make(map[string]interface{})
		msg    string
	)
	if len(keyvals) == 0 {
		return
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		if key == logrus.FieldKeyMsg {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		fields[key] = keyvals[i+1]
	}
	switch level {
	case log.DEBUG:
		a.log.DebugFields(msg, fields)
	case log.INFO:
		a.log.InfoFields(msg, fields)
	case log.WARN:
		a.log.WarnFields(msg, fields)
	case log.ERROR:
		a.log.ErrorFields(msg, fields)
	case log.FATAL:
		a.log.FatalFields(msg, fields)
	default:
		a.log.DebugFields(msg, fields)
	}

}
