package core

// CHECKED!

import (
	"sync"

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

const DefaultMailId = 0

// 邮件类型

type MailType int

const (
	// MailHalt
	// 终止邮件，用于停止Atomos的运行。
	// Halt Mail, for stopping an atomos from running.
	MailHalt MailType = 0

	// MailMessage
	// 信息邮件，用于外部给运行中的Atomos传递信息。
	// Message Mail, for messaging to a running atomos from outer.
	MailMessage MailType = 1

	// MailTask
	// 任务邮件，用于内部给运行中的Atomos新增任务。
	// Task Mail, for adding task to a running atomos from inner.
	MailTask MailType = 2

	// MailReload
	// 重载邮件，用于升级Atomos的ElementLocal引用，以实现热更。
	// Reload Mail, for upgrading ElementLocal reference of an atomos, to support hot-reload feature.
	MailReload MailType = 3
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
	// Halt, Message, Task, Reload
	mailType MailType

	// 从哪个Id发来的邮件。
	// Mail send from which ID.
	from ID

	// Message和Task邮件会使用到的，调用的目标对象的名称。
	// Mail target name, used by Message mail and Task mail.
	name string

	// Message和Task邮件的参数。
	// Argument that pass to target, used by Message mail and Task mail.
	arg proto.Message

	// 需要升级的Element。
	// Upgrade Element.
	reload  interface{}
	reloads int

	//wormholeAction int
	//wormhole       WormholeDaemon

	// 用于发邮件时阻塞调用go程，以及返回结果用的channel。
	// A channel used to block messaging goroutine, and return the result.
	mailReply mailReply
	waitCh    chan *mailReply
}

// Atomos邮件内存池
// Atomos Mails Pool
var atomosMailsPool = sync.Pool{
	New: func() interface{} {
		return &atomosMail{
			mailType:  0,
			from:      nil,
			name:      "",
			arg:       nil,
			mailReply: mailReply{},
			waitCh:    nil,
		}
	},
}

// Construct and destruct of Mail may be in different part of code.

func allocAtomosMail() *atomosMail {
	am := atomosMailsPool.Get().(*atomosMail)
	am.mail = newMail(DefaultMailId, am)
	return am
}

func deallocAtomosMail(am *atomosMail) {
	delMail(am.mail)
	atomosMailsPool.Put(am)
}

// 消息邮件
// Message Mail
func initMessageMail(am *atomosMail, from ID, name string, arg proto.Message) {
	am.mail.id = DefaultMailId
	am.mail.action = MailActionRun
	am.mailType = MailMessage
	am.from = from
	am.name = name
	// I think it has to be cloned, because argument is passing between atomos.
	if arg != nil {
		am.arg = proto.Clone(arg)
	} else {
		am.arg = nil
	}
	//am.reload = nil
	//am.upgradeCount = 0
	am.mailReply = mailReply{}
	am.waitCh = make(chan *mailReply, 1)
}

// 任务邮件
// Task Mail
func initTaskMail(am *atomosMail, taskId uint64, name string, arg proto.Message) {
	am.mail.id = taskId
	am.mail.action = MailActionRun
	am.mailType = MailTask
	am.from = nil
	am.name = name
	// I think it doesn't have to clone, because Atomos is thread-safe.
	am.arg = arg
	//am.reload = nil
	//am.upgradeCount = 0
	am.mailReply = mailReply{}
	am.waitCh = make(chan *mailReply, 1)
}

// 重载邮件
// Reload Mail
func initReloadMail(am *atomosMail, elem interface{}, reloads int) {
	am.mail.id = DefaultMailId
	am.mail.action = MailActionRun
	am.mailType = MailReload
	am.from = nil
	am.name = ""
	am.reload = elem
	am.reloads = reloads
	am.mailReply = mailReply{}
	am.waitCh = make(chan *mailReply, 1)
}

//func initWormholeMail(am *atomosMail, action int, wormhole WormholeDaemon) {
//	am.mail.id = DefaultMailId
//	am.mail.action = MailActionRun
//	am.mailType = AtomosMailWormhole
//	am.from = nil
//	am.name = ""
//	am.arg = nil
//	am.wormholeAction = action
//	am.wormhole = wormhole
//	am.reload = nil
//	am.upgradeCount = 0
//	am.mailReply = mailReply{}
//	am.waitCh = make(chan *mailReply, 1)
//}

// 终止邮件
// Halt Mail
func initKillMail(am *atomosMail, from ID) {
	am.mail.id = DefaultMailId
	am.mail.action = MailActionExit
	am.mailType = MailHalt
	am.from = from
	am.name = ""
	//am.reload = nil
	//am.upgradeCount = 0
	am.mailReply = mailReply{}
	am.waitCh = make(chan *mailReply, 1)
}

// Mail返回
// Mail Reply
type mailReply struct {
	resp proto.Message
	err  *ErrorInfo
}

// Method sendReply() will only be called in for-loop of MailBox, it's safe to do so, because while an atomos is
// waiting for replying, the atomos must still be running. Or if the atomos is not waiting for replying, after mailReply
// has been sent to waitCh, there will has no reference to the waitCh, waitCh will be collected.
func (m *atomosMail) sendReply(resp proto.Message, err *ErrorInfo) {
	m.mailReply.resp = resp
	m.mailReply.err = err
	if m.waitCh != nil {
		m.waitCh <- &m.mailReply
		m.waitCh = nil
	} else {
		panic("atomosMail: sendReply waitCh has been replied")
	}
}

// TODO: Think about waitReply() is still waiting when cosmos runnable is exiting.
func (m *atomosMail) waitReply() (resp proto.Message, err *ErrorInfo) {
	// An empty channel here means the receiver has received. It must be framework problem otherwise it won't happen.
	reply := <-m.waitCh
	// Wait channel must be empty before delete a mail.
	resp = reply.resp
	err = reply.err
	return resp, err
}
