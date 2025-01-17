package go_atomos

import (
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

//
// AtomosMail
//
// 关于Atomos的并发，我不采用Go语言常用的CSP模型，因为CSP模型有些问题，是比较难解决的：
// #1 Go的Channel在队列内容超出了容量之后，插入Channel的内容顺序不能确定，而且会阻塞发送方。
// #2 Go的Channel无法对内容进行"插队"或"取消"。
// #3 Go的Channel是单向的，没有系统的办法去处理回调问题。
//

const DefaultMailID = 0

// 邮件类型

type MailType int

const (
	// MailHalt
	// 终止邮件，用于停止Atomos的运行。
	// Stopping Mail, for stopping an atomos from running.
	MailHalt MailType = 0

	// MailMessage
	// 信息邮件，用于外部给运行中的Atomos传递信息。
	// Message Mail, for messaging to a running atomos from outer.
	MailMessage MailType = 1

	// MailAsyncMessage
	// 异步信息邮件，用于外部给运行中的Atomos传递信息。
	// Async Message Mail, for messaging to a running atomos from outer.
	MailAsyncMessage MailType = 2

	// MailAsyncMessageCallback
	// 异步信息回调邮件。
	MailAsyncMessageCallback MailType = 3

	// MailTask
	// 任务邮件，用于内部给运行中的Atomos新增任务。
	// Task Mail, for adding task to a running atomos from inner.
	MailTask MailType = 4

	// MailWormhole
	// 虫洞邮件，用于传递不属于"Atomos宇宙"概念的对象。
	// Wormhole Mail, for transporting non-"Atomos Cosmos" object.
	MailWormhole MailType = 5

	// MailScale
	// Scale邮件。
	MailScale MailType = 6
)

// Atomos邮件
// Atomos Mail

type atomosMail struct {
	// 具体的Mail实例
	// Concrete Mail instance.
	*mail

	// Atomos邮件类型
	// Atomos mail type.
	//
	// Stopping, Message, Task, Reload
	mailType MailType

	// 从哪个ID发来的邮件。
	// Mail send from which ID.
	from SelfID

	// 发送到哪个ID的邮件。
	// Mail send to which ID.
	to SelfID

	// Message和Task邮件会使用到的，调用的目标对象的名称。
	// Mail target name, used by Message mail and Task mail.
	name          string
	fromCallChain []string

	// 超时时间
	// Timeout
	timeout time.Duration

	// Message和Task邮件的参数。
	// Argument that pass to target, used by Message mail and Task mail.
	arg proto.Message
	err *Error

	tracker *IDTracker

	wormhole AtomosWormhole

	taskClosure func(uint64)

	asyncMessageCallbackClosure func(proto.Message, *Error)

	// 用于发邮件时阻塞调用go程，以及返回结果用的channel。
	// A channel used to block messaging goroutine, and return the result.
	mailReply mailReply
	waitCh    chan *mailReply

	executeStop bool

	mutex sync.Mutex
}

// Construct and destruct of Mail may be in different part of code.

func allocAtomosMail() *atomosMail {
	am := &atomosMail{}
	am.mail = &mail{mail: am}
	return am
}

func deallocAtomosMail(_ *atomosMail) {
}

// 消息邮件
// Message Mail
func initMessageMail(am *atomosMail, from SelfID, fromCallChain []string, name string, arg proto.Message) {
	am.mail.id = DefaultMailID
	am.mail.action = MailActionRun
	am.mailType = MailMessage
	am.from = from
	am.fromCallChain = fromCallChain
	am.name = name
	// I think it has to be cloned, because argument is passing between atomos.
	if arg != nil {
		if ShouldArgumentClone {
			am.arg = proto.Clone(arg)
		} else {
			am.arg = arg
		}
	} else {
		am.arg = nil
	}
	am.waitCh = make(chan *mailReply, 1)
}

// 异步消息邮件
func initAsyncMessageMail(am *atomosMail, from, to SelfID, name string, timeout time.Duration, callback func(proto.Message, *Error), arg proto.Message) {
	am.mail.id = DefaultMailID
	am.mail.action = MailActionRun
	am.mailType = MailAsyncMessage
	am.from = from
	am.to = to
	am.name = name
	am.timeout = timeout
	am.asyncMessageCallbackClosure = callback
	// I think it has to be cloned, because argument is passing between atomos.
	if arg != nil {
		if ShouldArgumentClone {
			am.arg = proto.Clone(arg)
		} else {
			am.arg = arg
		}
	} else {
		am.arg = nil
	}
}

// AsyncMessageCallback邮件
// Async Message Callback Mail
func initAsyncMessageCallbackMail(am *atomosMail, from SelfID, name string, callback func(proto.Message, *Error), arg proto.Message, err *Error) {
	am.mail.id = DefaultMailID
	am.mail.action = MailActionRun
	am.mailType = MailAsyncMessageCallback
	am.from = from
	am.fromCallChain = nil
	am.name = name
	am.arg = arg
	am.err = err
	am.asyncMessageCallbackClosure = callback
	am.waitCh = make(chan *mailReply, 1)
}

// Scale邮件
// Scale Mail
func initScaleMail(am *atomosMail, from SelfID, fromCallChain []string, name string, arg proto.Message) {
	am.mail.id = DefaultMailID
	am.mail.action = MailActionRun
	am.mailType = MailScale
	am.from = from
	am.fromCallChain = fromCallChain
	am.name = name
	// I think it has to be cloned, because argument is passing between atomos.
	if arg != nil {
		if ShouldArgumentClone {
			am.arg = proto.Clone(arg)
		} else {
			am.arg = arg
		}
	} else {
		am.arg = nil
	}
	am.tracker = nil
	am.wormhole = nil
	am.mailReply = mailReply{}
	am.executeStop = false
	am.waitCh = make(chan *mailReply, 1)
}

// 任务闭包邮件
// Task Closure Mail
// name中记录调用的闭包代码定位信息。
func initTaskClosureMail(am *atomosMail, name string, taskID uint64, closure func(uint64)) {
	am.mail.id = taskID
	am.mail.action = MailActionRun
	am.mailType = MailTask

	am.name = name
	am.taskClosure = closure
	am.waitCh = make(chan *mailReply, 1)
}

// 虫洞邮件
// Reload Mail
func initWormholeMail(am *atomosMail, from SelfID, fromCallChain []string, wormhole AtomosWormhole) {
	am.mail.id = DefaultMailID
	am.mail.action = MailActionRun
	am.mailType = MailWormhole
	am.from = from
	am.fromCallChain = fromCallChain
	am.name = ""
	am.arg = nil
	am.tracker = nil
	am.wormhole = wormhole
	am.mailReply = mailReply{}
	am.executeStop = false
	am.waitCh = make(chan *mailReply, 1)
}

// 终止邮件
// Stopping Mail
func initKillMail(am *atomosMail, from SelfID, fromCallChain []string) {
	am.mail.id = DefaultMailID
	am.mail.action = MailActionExit
	am.mailType = MailHalt
	am.from = from
	am.fromCallChain = fromCallChain
	am.name = ""
	am.tracker = nil
	am.wormhole = nil
	am.mailReply = mailReply{}
	am.executeStop = true
	am.waitCh = make(chan *mailReply, 1)
}

// Mail返回
// Mail Reply
type mailReply struct {
	resp proto.Message
	id   ID
	err  *Error
}

// Method sendReply() will only be called in for-loop of MailBox, it's safe to do so, because while an atomos is
// waiting for replying, the atomos must still be running. Or if the atomos is not waiting for replying, after mailReply
// has been sent to waitCh, there will have no reference to the waitCh, waitCh will be collected.
func (m *atomosMail) sendReply(resp proto.Message, err *Error) {
	m.mutex.Lock()
	waitCh := m.waitCh
	//m.waitCh = nil
	m.mutex.Unlock()
	if waitCh == nil {
		return
	}

	m.mailReply.resp = resp
	m.mailReply.err = err
	select {
	case waitCh <- &m.mailReply:
	default:
	}
	//waitCh = nil
}

func (m *atomosMail) sendReplyID(id ID, err *Error) {
	m.mutex.Lock()
	waitCh := m.waitCh
	//m.waitCh = nil
	m.mutex.Unlock()
	if waitCh == nil {
		return
	}

	m.mailReply.id = id
	m.mailReply.err = err
	waitCh <- &m.mailReply
}

func (m *atomosMail) asyncReply(resp proto.Message, err *Error) {
	if m.asyncMessageCallbackClosure == nil {
		// No callback, no reply.
		return
	}
	m.from.asyncCallback(m.from, m.name, resp, err, m.asyncMessageCallbackClosure)
}

// TODO: Think about waitReply() is still waiting when cosmos runnable is exiting.
func (m *atomosMail) waitReply(a *BaseAtomos, timeout time.Duration) (resp proto.Message, err *Error) {
	m.mutex.Lock()
	waitCh := m.waitCh
	m.mutex.Unlock()
	// An empty channel here means the receiver has received. It must be framework problem otherwise it won't happen.
	if waitCh == nil {
		return nil, NewErrorf(ErrFrameworkRecoverFromPanic, "Atomos Message wait invalid.").AddStack(nil)
	}
	var reply *mailReply
	if timeout == 0 {
		timeout = messageDeadlineDefault
	}
	select {
	case reply = <-waitCh:
	case <-time.After(timeout):
		if a.mailbox.removeMail(m.mail) {
			return nil, NewErrorf(ErrAtomosPushTimeoutReject, "Atomos: Message is timeout and rejected. id=(%s),name=(%s),timeout=(%v)", a.id.Info(), m.name, timeout).AddStack(nil)
		} else {
			return nil, NewErrorf(ErrAtomosPushTimeoutHandling, "Atomos: Message is handling timeout. id=(%s),name=(%s),timeout=(%v),current=(%s)", a.id.Info(), m.name, timeout, a.mt.current).AddStack(nil)
		}
	}
	// Wait channel must be empty before delete a mail.
	if reply == nil {
		return nil, NewErrorf(ErrFrameworkRecoverFromPanic, "Atomos: Message reply is invalid.").AddStack(nil)
	}
	resp = reply.resp
	err = reply.err
	return resp, err
}

// TODO: Think about waitReplyID() is still waiting when cosmos runnable is exiting.
func (m *atomosMail) waitReplyID(a *BaseAtomos, timeout time.Duration) (id ID, err *Error) {
	m.mutex.Lock()
	waitCh := m.waitCh
	m.mutex.Unlock()
	// An empty channel here means the receiver has received. It must be framework problem otherwise it won't happen.
	if waitCh == nil {
		return nil, NewErrorf(ErrFrameworkRecoverFromPanic, "Atomos: Message wait invalid.").AddStack(nil)
	}

	// An empty channel here means the receiver has received. It must be framework problem otherwise it won't happen.
	var reply *mailReply
	if timeout == 0 {
		timeout = messageDeadlineDefault
	}
	select {
	case reply = <-waitCh:
	case <-time.After(timeout):
		if a.mailbox.removeMail(m.mail) {
			return nil, NewErrorf(ErrAtomosPushTimeoutReject, "Atomos: Message is timeout and rejected. id=(%s),name=(%s),timeout=(%v)", a.id.Info(), m.name, timeout).AddStack(nil)
		} else {
			return nil, NewErrorf(ErrAtomosPushTimeoutHandling, "Atomos: Message is handling timeout. id=(%s),name=(%s),timeout=(%v)", a.id.Info(), m.name, timeout).AddStack(nil)
		}
	}
	// Wait channel must be empty before delete a mail.
	id = reply.id
	err = reply.err
	return id, err
}
