package go_atomos

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultLogMailID = 0

// Cosmos的Log接口。
// Interface of Cosmos Log.

type LoggingAtomos struct {
	logBox    *mailBox
	accessLog LoggingFn
	errorLog  LoggingFn
}

var sharedLogging LoggingAtomos

type LoggingFn func(string)

func SharedLogging() *LoggingAtomos {
	return &sharedLogging
}

func (c *LoggingAtomos) PushLogging(id *IDInfo, level LogLevel, msg string) {
	lm := &LogMail{
		Id:      id,
		Time:    timestamppb.Now(),
		Level:   level,
		Message: msg,
	}
	m := &mail{
		next:   nil,
		id:     defaultLogMailID,
		action: MailActionRun,
		mail:   nil,
		log:    lm,
	}
	if ok := c.logBox.pushTail(m); !ok {
		c.errorLog(fmt.Sprintf("LoggingAtomos: Add log mail failed. id=(%+v),level=(%v),msg=(%s)", id, level, msg))
	}
}

func initSharedLoggingAtomos(accessLog, errLog LoggingFn) {
	sharedLogging = LoggingAtomos{
		logBox: newMailBox("sharedLogging", MailBoxHandler{
			OnReceive: sharedLogging.onLogMessage,
			OnStop:    sharedLogging.onLogStop,
		}),
		accessLog: accessLog,
		errorLog:  errLog,
	}
	sharedLogging.logBox.start()
}

func (c *LoggingAtomos) pushFrameworkErrorLog(format string, args ...interface{}) {
	c.PushLogging(&IDInfo{
		Type:    IDType_Process,
		Cosmos:  "",
		Element: "",
		Atomos:  "",
	}, LogLevel_Fatal, fmt.Sprintf(format, args...))
}

func (c *LoggingAtomos) pushProcessLog(level LogLevel, format string, args ...interface{}) {
	c.PushLogging(&IDInfo{
		Type:    IDType_Process,
		Cosmos:  "",
		Element: "",
		Atomos:  "",
	}, level, fmt.Sprintf(format, args...))
}

// Logging Atomos的实现。
// Implementation of Logging Atomos.

func (c *LoggingAtomos) onLogMessage(mail *mail) {
	c.logging(mail.log)
}

func (c *LoggingAtomos) onLogStop(killMail, remainMails *mail, num uint32) {
	for curMail := remainMails; curMail != nil; curMail = curMail.next {
		c.onLogMessage(curMail)
	}
}

func (c *LoggingAtomos) logging(lm *LogMail) {
	var msg string
	if id := lm.Id; id != nil {
		switch id.Type {
		case IDType_Atomos:
			msg = fmt.Sprintf("%s::%s::%s => %s", id.Cosmos, id.Element, id.Atomos, lm.Message)
		case IDType_Element:
			msg = fmt.Sprintf("%s::%s => %s", id.Cosmos, id.Element, lm.Message)
		case IDType_Cosmos:
			msg = fmt.Sprintf("%s => %s", id.Cosmos, lm.Message)
		case IDType_Main:
			msg = fmt.Sprintf("Main => %s", lm.Message)
		default:
			msg = fmt.Sprintf("Unknown => %s", lm.Message)
		}
	} else {
		msg = fmt.Sprintf("%s", lm.Message)
	}
	switch lm.Level {
	case LogLevel_Debug:
		c.accessLog(fmt.Sprintf("%s [DEBUG] %s\n", lm.Time.AsTime().Format(logTimeFmt), msg))
	case LogLevel_Info:
		c.accessLog(fmt.Sprintf("%s [INFO]  %s\n", lm.Time.AsTime().Format(logTimeFmt), msg))
	case LogLevel_Warn:
		c.accessLog(fmt.Sprintf("%s [WARN]  %s\n", lm.Time.AsTime().Format(logTimeFmt), msg))
	case LogLevel_Err:
		c.errorLog(fmt.Sprintf("%s [ERROR] %s\n", lm.Time.AsTime().Format(logTimeFmt), msg))
	case LogLevel_Fatal:
		fallthrough
	default:
		c.errorLog(fmt.Sprintf("%s [FATAL] %s\n", lm.Time.AsTime().Format(logTimeFmt), msg))
	}
}
