package go_atomos

// CHECKED!

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
)

// Cosmos的Log接口。
// Interface of Cosmos Log.

func (c *CosmosSelf) Debug(format string, args ...interface{}) {
	if c.config.LogLevel > LogLevel_Debug {
		return
	}
	c.pushCosmosLog(LogLevel_Debug, fmt.Sprintf(format, args...))
}

func (c *CosmosSelf) logInfo(format string, args ...interface{}) {
	if c.config.LogLevel > LogLevel_Info {
		return
	}
	c.pushCosmosLog(LogLevel_Info, fmt.Sprintf(format, args...))
}

func (c *CosmosSelf) logWarn(format string, args ...interface{}) {
	if c.config.LogLevel > LogLevel_Warn {
		return
	}
	c.pushCosmosLog(LogLevel_Warn, fmt.Sprintf(format, args...))
}

func (c *CosmosSelf) logError(format string, args ...interface{}) {
	if c.config.LogLevel > LogLevel_Error {
		return
	}
	c.pushCosmosLog(LogLevel_Error, fmt.Sprintf(format, args...))
}

func (c *CosmosSelf) logFatal(format string, args ...interface{}) {
	if c.config.LogLevel > LogLevel_Fatal {
		return
	}
	c.pushCosmosLog(LogLevel_Fatal, fmt.Sprintf(format, args...))
}

// Cosmos的Log实现。
// Implementation of Cosmos Log.

func (c *CosmosSelf) onLogMessage(mail *Mail) {
	lm := mail.Content.(*LogMail)
	c.logging(lm)
	logMailsPool.Put(lm)
	DelMail(mail)
}

func (c *CosmosSelf) onLogPanic(mail *Mail, trace string) {
	lm := mail.Content.(*LogMail)
	c.logging(&LogMail{
		Id:      lm.Id,
		Time:    lm.Time,
		Level:   LogLevel_Fatal,
		Message: trace,
	})
}

func (c *CosmosSelf) onLogStop(killMail, remainMails *Mail, num uint32) {
	for curMail := remainMails; curMail != nil; curMail = curMail.next {
		c.onLogMessage(curMail)
	}
}

const logTimeFmt = "2006-01-02 15:04:05.000000"

func (c *CosmosSelf) logging(lm *LogMail) {
	var msg string
	if lm.Id != nil {
		msg = fmt.Sprintf("%s::%s::%s => %s", lm.Id.Node, lm.Id.Element, lm.Id.Name, lm.Message)
	} else {
		msg = fmt.Sprintf("%s", lm.Message)
	}
	switch lm.Level {
	case LogLevel_Debug:
		logWrite(fmt.Sprintf("%s [DEBUG] %s\n", lm.Time.AsTime().Format(logTimeFmt), msg), false)
	case LogLevel_Info:
		logWrite(fmt.Sprintf("%s [INFO]  %s\n", lm.Time.AsTime().Format(logTimeFmt), msg), false)
	case LogLevel_Warn:
		logWrite(fmt.Sprintf("%s [WARN]  %s\n", lm.Time.AsTime().Format(logTimeFmt), msg), false)
	case LogLevel_Error:
		logWrite(fmt.Sprintf("%s [ERROR] %s\n", lm.Time.AsTime().Format(logTimeFmt), msg), true)
	case LogLevel_Fatal:
		logWrite(fmt.Sprintf("%s [FATAL] %s\n", lm.Time.AsTime().Format(logTimeFmt), msg), true)
	default:
		logWrite(fmt.Sprintf("%s [WARN]  %s\n", lm.Time.AsTime().Format(logTimeFmt), msg), true)
	}
}

func (c *CosmosSelf) pushCosmosLog(level LogLevel, msg string) {
	lm := logMailsPool.Get().(*LogMail)
	lm.Id = nil
	lm.Time = timestamppb.Now()
	lm.Level = level
	lm.Message = msg
	m := NewMail(defaultLogMailId, lm)
	if ok := c.log.PushTail(m); !ok {
		log.Println("Cosmos Log Mail failed", level, msg)
	}
}

// Concrete log to file logic.

func logWrite(msg string, err bool) {
	if err {
		os.Stderr.WriteString(msg)
	}
	os.Stdout.WriteString(msg)
}
