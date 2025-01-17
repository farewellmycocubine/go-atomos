package go_atomos

import (
	"bytes"
	"fmt"
	"io"
)

type Logging interface {
	io.Writer

	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	coreInfo(format string, args ...interface{})

	Error(format string, args ...interface{})
	coreError(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	coreFatal(format string, args ...interface{})
}

// Atomos日志管理器
// Atomos Logs Manager

type atomosLogging struct {
	id      *IDInfo
	level   LogLevel
	logging *loggingAtomos
}

// 初始化atomosLogsManager的内容。
// 没有构造和释构函数，因为atomosLogsManager是AtomCore内部使用的。
//
// Initialization of atomosLogging.
// No New and Delete function because atomosLogging is struct inner AtomCore.
func initAtomosLog(l *atomosLogging, id *IDInfo, lv LogLevel, logging *loggingAtomos) {
	l.id = id
	l.level = lv
	l.logging = logging
}

// 把Log以邮件的方式发送到Cosmos的Log实例处理。
// write Logs as Mails to Cosmos Log instance.
func (l *atomosLogging) pushAtomosLog(id *IDInfo, level LogLevel, msg string) {
	l.logging.PushLogging(id, level, msg)
}

func (l *atomosLogging) Write(p []byte) (n int, err error) {
	// Remove return character, if any.
	p = bytes.TrimRight(p, "\n")
	l.pushAtomosLog(l.id, LogLevel_CoreInfo, string(p))
	return len(p), nil
}

// 各种级别的日志函数。
// Log functions in difference levels.

func (l *atomosLogging) Debug(format string, args ...interface{}) {
	if l.level > LogLevel_Debug {
		return
	}
	l.pushAtomosLog(l.id, LogLevel_Debug, fmt.Sprintf(format, args...))
}

func (l *atomosLogging) Info(format string, args ...interface{}) {
	if l.level > LogLevel_Info {
		return
	}
	l.pushAtomosLog(l.id, LogLevel_Info, fmt.Sprintf(format, args...))
}

func (l *atomosLogging) coreInfo(format string, args ...interface{}) {
	if l.level > LogLevel_CoreInfo {
		return
	}
	l.pushAtomosLog(l.id, LogLevel_CoreInfo, fmt.Sprintf(format, args...))
}

func (l *atomosLogging) Warn(format string, args ...interface{}) {
	if l.level > LogLevel_Warn {
		return
	}
	l.pushAtomosLog(l.id, LogLevel_Warn, fmt.Sprintf(format, args...))
}

func (l *atomosLogging) Error(format string, args ...interface{}) {
	if l.level > LogLevel_Err {
		return
	}
	l.pushAtomosLog(l.id, LogLevel_Err, fmt.Sprintf(format, args...))
}

func (l *atomosLogging) coreError(format string, args ...interface{}) {
	if l.level > LogLevel_CoreErr {
		return
	}
	l.pushAtomosLog(l.id, LogLevel_CoreErr, fmt.Sprintf(format, args...))
}

func (l *atomosLogging) Fatal(format string, args ...interface{}) {
	if l.level > LogLevel_Fatal {
		return
	}
	l.pushAtomosLog(l.id, LogLevel_Fatal, fmt.Sprintf(format, args...))
}

func (l *atomosLogging) coreFatal(format string, args ...interface{}) {
	if l.level > LogLevel_CoreFatal {
		return
	}
	l.pushAtomosLog(l.id, LogLevel_CoreErr, fmt.Sprintf(format, args...))
}
