package go_atomos

// CHECKED!

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

// ElementLocal
// 本地Element实现。
// Implementation of local Element.
type ElementLocal struct {
	// CosmosSelf引用。
	// Reference to CosmosSelf.
	main *CosmosMain

	// 基础Atomos，也是实现Atom无锁队列的关键。
	// Base atomos, the key of lockless queue of Atom.
	atomos *BaseAtomos

	// 该Element所有Atom的容器。
	// Container of all atoms.
	// 思考：要考虑在频繁变动的情景下，迭代不全的问题。
	// 两种情景：更新&关闭。
	atoms map[string]*AtomLocal
	// Element的List容器的
	names *list.List
	// Lock.
	lock sync.RWMutex

	//// Available or Reloading
	//avail bool
	// 当前ElementImplementation的引用。
	// Reference to current in use ElementImplementation.
	current *ElementImplementation

	// 调用链
	// 调用链用于检测是否有循环调用，在处理message时把fromID的调用链加上自己之后
	callIDCounter uint64
	curCallChain  string

	messageTracker *MessageTrackerManager
	idTracker      *IDTrackerManager
}

// 生命周期相关
// Life Cycle

// 本地Element创建，用于本地Cosmos的创建过程。
// Create of the Local Element, uses in Local Cosmos creation.
func newElementLocal(main *CosmosMain, runnable *CosmosRunnable, impl *ElementImplementation) *ElementLocal {
	id := &IDInfo{
		Type:    IDType_Element,
		Cosmos:  runnable.config.Node,
		Element: impl.Interface.Config.Name,
		Atom:    "",
	}
	e := &ElementLocal{
		main:           main,
		atomos:         nil,
		atoms:          nil,
		names:          list.New(),
		lock:           sync.RWMutex{},
		current:        impl,
		callIDCounter:  0,
		curCallChain:   "",
		messageTracker: NewMessageTrackerManager(id, len(impl.ElementHandlers)),
		idTracker:      nil,
	}
	e.atomos = NewBaseAtomos(id, impl.Interface.Config.LogLevel, e, impl.Developer.ElementConstructor())
	e.idTracker = NewIDTrackerManager(e)
	if atomsInitNum, ok := impl.Developer.(ElementCustomizeAtomInitNum); ok {
		num := atomsInitNum.GetElementAtomsInitNum()
		e.atoms = make(map[string]*AtomLocal, num)
	} else {
		e.atoms = map[string]*AtomLocal{}
	}
	return e
}

//
// Implementation of ID
//

// ID，相当于Atom的句柄的概念。
// 通过ID，可以访问到Atom所在的Cosmos、Element、Name，以及发送Kill信息，但是否能成功Kill，还需要AtomCanKill函数的认证。
// 直接用AtomLocal继承ID，因此本地的ID直接使用AtomLocal的引用即可。
//
// ID, a concept similar to file descriptor of an atomos.
// With ID, we can access the Cosmos, Element and Name of the Atom. We can also send Kill signal to the Atom,
// then the AtomCanKill method judge kill it or not.
// AtomLocal implements ID interface directly, so local ID is able to use AtomLocal reference directly.

func (e *ElementLocal) GetIDInfo() *IDInfo {
	if e == nil {
		return nil
	}
	return e.atomos.GetIDInfo()
}

func (e *ElementLocal) String() string {
	if e == nil {
		return "nil"
	}
	return e.atomos.String()
}

func (e *ElementLocal) Release(id *IDTracker) {
	e.idTracker.Release(id)
}

func (e *ElementLocal) Cosmos() CosmosNode {
	return e.main
}

func (e *ElementLocal) Element() Element {
	return e
}

func (e *ElementLocal) GetName() string {
	return e.GetIDInfo().Element
}

func (e *ElementLocal) State() AtomosState {
	return e.atomos.GetState()
}

func (e *ElementLocal) IdleTime() time.Duration {
	e.atomos.mailbox.mutex.Lock()
	defer e.atomos.mailbox.mutex.Unlock()
	return e.messageTracker.idleTime()
}

func (e *ElementLocal) MessageByName(from ID, name string, timeout time.Duration, in proto.Message) (proto.Message, *Error) {
	return e.pushMessageMail(from, name, timeout, in)
}

func (e *ElementLocal) DecoderByName(name string) (MessageDecoder, MessageDecoder) {
	decoderFn, has := e.current.Interface.ElementDecoders[name]
	if !has {
		return nil, nil
	}
	return decoderFn.InDec, decoderFn.OutDec
}

func (e *ElementLocal) Kill(from ID, timeout time.Duration) *Error {
	return NewError(ErrElementCannotKill, "Element: Cannot kill element.")
}

func (e *ElementLocal) SendWormhole(from ID, timeout time.Duration, wormhole AtomosWormhole) *Error {
	return e.atomos.PushWormholeMailAndWaitReply(from, timeout, wormhole)
}

func (e *ElementLocal) getElementLocal() *ElementLocal {
	return e
}

func (e *ElementLocal) getAtomLocal() *AtomLocal {
	return nil
}

func (e *ElementLocal) getElementRemote() *ElementRemote {
	return nil
}

func (e *ElementLocal) getAtomRemote() *AtomRemote {
	return nil
}

func (e *ElementLocal) getIDTrackerManager() *IDTrackerManager {
	return e.idTracker
}

func (e *ElementLocal) getCurCallChain() string {
	e.atomos.mailbox.mutex.Lock()
	c := e.curCallChain
	e.atomos.mailbox.mutex.Unlock()
	return c
}

func (e *ElementLocal) First() ID {
	e.atomos.mailbox.mutex.Lock()
	e.callIDCounter += 1
	callID := e.callIDCounter
	e.atomos.mailbox.mutex.Unlock()
	return &FirstID{callID: callID, ID: e}
}

// Implementation of atomos.SelfID
// Implementation of atomos.ParallelSelf
//
// SelfID，是Atom内部可以访问的Atom资源的概念。
// 通过AtomSelf，Atom内部可以访问到自己的Cosmos（CosmosSelf）、可以杀掉自己（KillSelf），以及提供Log和Task的相关功能。
//
// SelfID, a concept that provide Atom resource access to inner Atom.
// With SelfID, Atom can access its self-main with "CosmosSelf", can kill itself use "KillSelf" from inner.
// It also provides Log and Tasks method to inner Atom.

func (e *ElementLocal) CosmosMain() *CosmosMain {
	return e.main
}

// KillSelf
// Atom kill itself from inner
func (e *ElementLocal) KillSelf() {
	if err := e.pushKillMail(e, false, 0); err != nil {
		e.Log().Error("Element: KillSelf failed. err=(%v)", err)
		return
	}
	e.Log().Info("Element: KillSelf")
}

func (e *ElementLocal) Parallel(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := NewErrorf(ErrFrameworkPanic, "Element: Parallel recovers from panic.").AddPanicStack(e, 3, r)
				if ar, ok := e.atomos.instance.(AtomosRecover); ok {
					defer func() {
						recover()
						e.Log().Fatal("Element: Parallel recovers from panic. err=(%v)", err)
					}()
					ar.ParallelRecover(err)
				} else {
					e.Log().Fatal("Element: Parallel recovers from panic. err=(%v)", err)
				}
			}
		}()
		fn()
	}()
}

// Implementation of AtomSelfID

func (e *ElementLocal) Config() map[string][]byte {
	return e.main.runnable.config.Customize
}

func (e *ElementLocal) Persistence() ElementCustomizeAutoDataPersistence {
	p, _ := e.atomos.instance.(ElementCustomizeAutoDataPersistence)
	return p
}

func (e *ElementLocal) GetAtoms() []*AtomLocal {
	e.lock.RLock()
	atoms := make([]*AtomLocal, 0, len(e.atoms))
	for _, atomLocal := range e.atoms {
		if atomLocal.atomos.IsInState(AtomosSpawning, AtomosWaiting, AtomosBusy) {
			atoms = append(atoms, atomLocal)
		}
	}
	e.lock.RUnlock()
	return atoms
}

func (e *ElementLocal) MessageSelfByName(from ID, name string, buf []byte, protoOrJSON bool) ([]byte, *Error) {
	handlerFn, has := e.current.ElementHandlers[name]
	if !has {
		return nil, NewErrorf(ErrElementMessageHandlerNotExists, "Element: Handler not exists. from=(%v),name=(%s)", from, name).AddStack(nil)
	}
	decoderFn, has := e.current.Interface.ElementDecoders[name]
	if !has {
		return nil, NewErrorf(ErrElementMessageDecoderNotExists, "Element: Decoder not exists. from=(%v),name=(%s)", from, name).AddStack(nil)
	}
	in, err := decoderFn.InDec(buf, protoOrJSON)
	if err != nil {
		return nil, err
	}
	var outBuf []byte
	out, err := handlerFn(from, e.atomos.instance, in)
	if out != nil {
		var e error
		outBuf, e = proto.Marshal(out)
		if e != nil {
			return nil, NewErrorf(ErrElementMessageReplyType, "Element: Reply marshal failed. err=(%v)", err)
		}
	}
	return outBuf, err
}

// Implementation of AtomosUtilities

func (e *ElementLocal) Log() Logging {
	return e.atomos.Log()
}

func (e *ElementLocal) Task() Task {
	return e.atomos.Task()
}

// Implementation of Element

func (e *ElementLocal) GetElementName() string {
	return e.GetIDInfo().Element
}

func (e *ElementLocal) GetAtomID(name string, tracker *IDTrackerInfo) (ID, *IDTracker, *Error) {
	return e.elementAtomGet(name, tracker)
}

func (e *ElementLocal) GetAtomsNum() int {
	e.lock.RLock()
	num := len(e.atoms)
	e.lock.RUnlock()
	return num
}

func (e *ElementLocal) GetActiveAtomsNum() int {
	num := 0
	e.lock.RLock()
	for _, atomLocal := range e.atoms {
		if atomLocal.atomos.IsInState(AtomosSpawning, AtomosWaiting, AtomosBusy) {
			num += 1
		}
	}
	e.lock.RUnlock()
	return num
}

func (e *ElementLocal) GetAllInactiveAtomsIDTrackerInfo() map[string]string {
	e.lock.RLock()
	info := make(map[string]string, len(e.atoms))
	atoms := make([]*AtomLocal, 0, len(e.atoms))
	for _, atomLocal := range e.atoms {
		if atomLocal.atomos.IsInState(AtomosHalt) {
			atoms = append(atoms, atomLocal)
		}
	}
	e.lock.RUnlock()
	for _, atomLocal := range atoms {
		info[atomLocal.String()] = fmt.Sprintf(" -> %s\n", atomLocal.idTracker)
	}
	return info
}

func (e *ElementLocal) SpawnAtom(name string, arg proto.Message, tracker *IDTrackerInfo) (*AtomLocal, *IDTracker, *Error) {
	e.lock.RLock()
	current := e.current
	e.lock.RUnlock()
	// Auto data persistence.
	persistence, _ := current.Developer.(ElementCustomizeAutoDataPersistence)
	return e.elementAtomSpawn(name, arg, current, persistence, tracker)
}

func (e *ElementLocal) MessageElement(fromID, toID ID, name string, timeout time.Duration, args proto.Message) (reply proto.Message, err *Error) {
	if fromID == nil {
		return reply, NewErrorf(ErrAtomFromIDInvalid, "Element: MessageElement, FromID invalid. from=(%s),to=(%s),name=(%s),args=(%v)",
			fromID, toID, name, args).AddStack(e)
	}
	elem := toID.getElementLocal()
	if elem == nil {
		return reply, NewErrorf(ErrAtomToIDInvalid, "Element: MessageElement, ToID invalid. from=(%s),to=(%s),name=(%s),args=(%v)",
			fromID, toID, name, args).AddStack(e)
	}
	return elem.pushMessageMail(fromID, name, timeout, args)
}

func (e *ElementLocal) MessageAtom(fromID, toID ID, name string, timeout time.Duration, args proto.Message) (reply proto.Message, err *Error) {
	if fromID == nil {
		return reply, NewErrorf(ErrAtomFromIDInvalid, "Element: MessageAtom, FromID invalid. from=(%s),to=(%s),name=(%s),args=(%v)",
			fromID, toID, name, args).AddStack(e)
	}
	a := toID.getAtomLocal()
	if a == nil {
		return reply, NewErrorf(ErrAtomToIDInvalid, "Element: MessageAtom, ToID invalid. from=(%s),to=(%s),name=(%s),args=(%v)",
			fromID, toID, name, args).AddStack(e)
	}
	return a.pushMessageMail(fromID, name, timeout, args)
}

func (e *ElementLocal) ScaleGetAtomID(fromID ID, name string, timeout time.Duration, args proto.Message, tracker *IDTrackerInfo) (ID, *IDTracker, *Error) {
	if fromID == nil {
		return nil, nil, NewErrorf(ErrAtomFromIDInvalid, "Element: ScaleGetAtomID, FromID invalid. from=(%s),name=(%s),args=(%v)",
			fromID, name, args).AddStack(e)
	}
	return e.pushScaleMail(fromID, name, timeout, args, tracker)
}

func (e *ElementLocal) KillAtom(fromID, toID ID, timeout time.Duration) *Error {
	if fromID == nil {
		return NewErrorf(ErrAtomFromIDInvalid, "Element: KillAtom, FromID invalid. from=(%s),to=(%s)",
			fromID, toID).AddStack(e)
	}
	a := toID.getElementLocal()
	if a == nil {
		return NewErrorf(ErrAtomToIDInvalid, "Element: KillAtom, ToID invalid. from=(%s),to=(%s)",
			fromID, toID).AddStack(e)
	}
	return a.pushKillMail(fromID, true, timeout)
}

// Check chain.

func (e *ElementLocal) isInChain(callChainID string) bool {
	e.atomos.mailbox.mutex.Lock()
	defer e.atomos.mailbox.mutex.Unlock()

	return e.curCallChain == callChainID
}

func (e *ElementLocal) setMessageAndCallChain(callChain string) *Error {
	e.atomos.mailbox.mutex.Lock()
	defer e.atomos.mailbox.mutex.Unlock()
	e.curCallChain = callChain
	return nil
}

func (e *ElementLocal) unsetMessageAndCallChain() {
	e.atomos.mailbox.mutex.Lock()
	defer e.atomos.mailbox.mutex.Unlock()
	e.curCallChain = ""
}

// 内部实现
// INTERNAL

// 邮箱控制器相关
// Mailbox Handler

func (e *ElementLocal) pushMessageMail(from ID, name string, timeout time.Duration, arg proto.Message) (reply proto.Message, err *Error) {
	// Dead Lock Checker.
	// OnMessaging处理消息的时候，才做addChain操作
	if from == nil {
		return nil, NewError(ErrElementNoFromID, "Element: No fromID.").AddStack(e)
	}

	_, ok := from.(*FirstID)
	if !ok && e.curCallChain != "" && e.isInChain(from.getCurCallChain()) {
		return reply, NewErrorf(ErrAtomosCallDeadLock, "Element: Call Dead Lock. to=(%v),name=(%s),arg=(%v)", e, name, arg).AddStack(e)
	}
	return e.atomos.PushMessageMailAndWaitReply(from, name, timeout, arg)
}

func (e *ElementLocal) OnMessaging(from ID, name string, arg proto.Message) (reply proto.Message, err *Error) {
	if from == nil {
		return nil, NewError(ErrAtomNoFromID, "Atom: No fromID.").AddStack(e)
	}
	var fromChain string
	if f, ok := from.(*FirstID); ok {
		fromChain = fmt.Sprintf("%s:%d", f.ID, f.callID)
	} else {
		fromChain = from.getCurCallChain()
	}
	if err = e.setMessageAndCallChain(fromChain); err != nil {
		return nil, err.AddStack(e)
	}
	defer e.unsetMessageAndCallChain()
	handler := e.current.ElementHandlers[name]
	if handler == nil {
		return nil, NewErrorf(ErrElementMessageHandlerNotExists,
			"Element: Message handler not found. from=(%s),name=(%s),args=(%v)", from, name, arg)
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if err == nil {
					err = NewErrorf(ErrFrameworkPanic, "Element: Messaging recovers from panic.").AddPanicStack(e, 3, r)
					if ar, ok := e.atomos.instance.(AtomosRecover); ok {
						defer func() {
							recover()
							e.Log().Fatal("Element: Messaging recovers from panic. err=(%v)", err)
						}()
						ar.MessageRecover(name, arg, err)
					} else {
						e.Log().Fatal("Element: Messaging recovers from panic. err=(%v)", err)
					}
				}
			}
		}()
		fromID, _ := from.(ID)
		reply, err = handler(fromID, e.atomos.GetInstance(), arg)
	}()
	return
}

func (e *ElementLocal) pushScaleMail(fromID ID, name string, timeout time.Duration, arg proto.Message, t *IDTrackerInfo) (ID, *IDTracker, *Error) {
	// Dead Lock Checker.
	if fromID != nil {
		cc := fromID.getCurCallChain()
		if cc != "" && e.isInChain(cc) {
			return nil, nil, NewErrorf(ErrAtomosCallDeadLock, "Element: Call Dead Lock. to=(%v),name=(%s),arg=(%v)",
				e, name, arg).AddStack(e)
		}
	}
	tracker := newIDTracker(t)
	id, err := e.atomos.PushScaleMailAndWaitReply(fromID, name, timeout, arg, tracker)
	if err != nil {
		return nil, nil, err.AddStack(e, &String{S: name}, arg)
	}
	return id, tracker, nil
}

func (e *ElementLocal) OnScaling(from ID, name string, arg proto.Message, tracker *IDTracker) (id ID, err *Error) {
	if from == nil {
		return nil, NewError(ErrAtomNoFromID, "Atom: No fromID.").AddStack(e)
	}
	var fromChain string
	if f, ok := from.(*FirstID); ok {
		fromChain = fmt.Sprintf("%s:%d", f.ID, f.callID)
	} else {
		fromChain = from.getCurCallChain()
	}
	if err = e.setMessageAndCallChain(fromChain); err != nil {
		return nil, err.AddStack(e)
	}
	defer e.unsetMessageAndCallChain()
	handler := e.current.ScaleHandlers[name]
	if handler == nil {
		return nil, NewErrorf(ErrElementScaleHandlerNotExists,
			"Element: Scale handler not found. from=(%s),name=(%s),arg=(%v)", from, name, arg)
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if err == nil {
					err = NewErrorf(ErrFrameworkPanic, "Element: Scaling recovers from panic.").AddPanicStack(e, 3, r)
					if ar, ok := e.atomos.instance.(AtomosRecover); ok {
						defer func() {
							recover()
							e.Log().Fatal("Element: Scaling recovers from panic. err=(%v)", err)
						}()
						ar.ScaleRecover(name, arg, err)
					} else {
						e.Log().Fatal("Element: Scaling recovers from panic. err=(%v)", err)
					}
				}
			}
		}()
		id, err = handler(from, e.atomos.instance, name, arg)
		// Retain New.
		id.getIDTrackerManager().NewScaleIDTracker(tracker)
		// Release Old.
		releasable, ok := id.(ReleasableID)
		if ok {
			releasable.Release()
		}
	}()
	return
}

func (e *ElementLocal) pushKillMail(from ID, wait bool, timeout time.Duration) *Error {
	// Dead Lock Checker.
	if from != nil && wait {
		cc := from.getCurCallChain()
		if cc != "" && e.isInChain(cc) {
			return NewErrorf(ErrAtomosCallDeadLock, "Element: Kill Deadlock. from=(%v),to(%s)", from, e).AddStack(e)
		}
	}
	return e.atomos.PushKillMailAndWaitReply(from, wait, true, timeout)
}

func (e *ElementLocal) OnStopping(from ID, cancelled map[uint64]CancelledTask) (err *Error) {
	impl := e.current
	if impl == nil {
		return NewErrorf(ErrAtomKillElementNoImplement,
			"Element: OnStopping, no element implement. id=(%s),element=(%+v)", e.GetIDInfo(), e.GetElementName()).AddStack(e)
	}
	// Atoms
	// Send Kill to all atoms.
	var stopTimeout, stopGap time.Duration
	elemExit, ok := impl.Developer.(ElementCustomizeExit)
	if ok && elemExit != nil {
		stopTimeout = elemExit.StopTimeout()
		stopGap = elemExit.StopGap()
	}
	exitWG := sync.WaitGroup{}
	for nameElem := e.names.Back(); nameElem != nil; nameElem = nameElem.Prev() {
		name := nameElem.Value.(string)
		atom, has := e.atoms[name]
		if !has {
			continue
		}
		e.Log().Info("Element: Kill atom. name=(%s)", name)
		exitWG.Add(1)
		go func(a *AtomLocal, n string) {
			if err := a.pushKillMail(e, true, stopTimeout); err != nil {
				e.Log().Error("Element: Kill atom failed. name=(%s),err=(%v)", n, err)
			}
			exitWG.Done()
		}(atom, name)
		if stopGap > 0 {
			<-time.After(stopGap)
		}
	}
	exitWG.Wait()
	e.Log().Info("Element: All atoms killed. element=(%s)", e.GetName())

	if StoppingPrintStatic {
		e.Log().Warn("Static >> AtomStopping IDTracker=(%v)", e.idTracker)
		e.Log().Warn("Static >> AtomStopping MessageTracker=(%v)", e.GetMessagingInfo())
	}

	var persistence ElementCustomizeAutoDataPersistence
	var elemPersistence ElementAutoDataPersistence

	// Element
	save, data := e.atomos.GetInstance().Halt(from, cancelled)
	if !save {
		goto autoLoad
	}

	// Save data.
	// Auto Save
	persistence, ok = impl.Developer.(ElementCustomizeAutoDataPersistence)
	if !ok || persistence == nil {
		err = NewErrorf(ErrAtomKillElementNotImplementAutoDataPersistence,
			"Element: Save data error, no auto data persistence. id=(%s),element=(%+v)", e.GetIDInfo(), e.GetElementName()).AddStack(e)
		e.Log().Fatal(err.Error())
		goto autoLoad
	}
	elemPersistence = persistence.ElementAutoDataPersistence()
	if elemPersistence == nil {
		err = NewErrorf(ErrAtomKillElementNotImplementAutoDataPersistence,
			"Element: Save data error, no element auto data persistence. id=(%s),element=(%+v)", e.GetIDInfo(), e.GetElementName()).AddStack(e)
		e.Log().Fatal(err.Error())
		return err
	}
	if err = elemPersistence.SetElementData(data); err != nil {
		e.Log().Error("Element: Save data failed, set atom data error. id=(%s),instance=(%+v),err=(%s)",
			e.GetIDInfo(), e.atomos.String(), err)
		goto autoLoad
	}
autoLoad:
	// Auto Load
	pa, ok := impl.Developer.(ElementCustomizeAutoLoadPersistence)
	if !ok || pa == nil {
		return nil
	}
	if err = pa.Unload(); err != nil {
		e.Log().Error("Element: Unload failed. id=(%s),instance=(%+v),err=(%s)",
			e.GetIDInfo(), e.atomos.String(), err)
		return err
	}
	return err
}

func (e *ElementLocal) OnWormhole(from ID, wormhole AtomosWormhole) *Error {
	holder, ok := e.atomos.instance.(AtomosAcceptWormhole)
	if !ok || holder == nil {
		err := NewErrorf(ErrAtomosNotSupportWormhole, "ElementLocal: Not supported wormhole, type=(%T)", e.atomos.instance)
		e.Log().Error(err.Message)
		return err
	}
	return holder.AcceptWormhole(from, wormhole)
}

// Set & Unset

func (e *ElementLocal) Spawn() {
	e.messageTracker.Start()
}

func (e *ElementLocal) Set(message string) {
	e.messageTracker.Set(message)
}

func (e *ElementLocal) Unset(message string) {
	e.messageTracker.Unset(message)
}

func (e *ElementLocal) Stopping() {
	e.messageTracker.Stopping()
}

func (e *ElementLocal) Halted() {
	e.messageTracker.Halt()
}

func (e *ElementLocal) GetMessagingInfo() string {
	e.atomos.mailbox.mutex.Lock()
	defer e.atomos.mailbox.mutex.Unlock()
	return e.messageTracker.dump()
}

// Internal

func (e *ElementLocal) elementAtomGet(name string, t *IDTrackerInfo) (*AtomLocal, *IDTracker, *Error) {
	e.lock.RLock()
	current := e.current
	atom, hasAtom := e.atoms[name]
	e.lock.RUnlock()
	if hasAtom && atom.atomos.isNotHalt() {
		return atom, atom.idTracker.NewIDTracker(t), nil
	}
	// Auto data persistence.
	persistence, ok := current.Developer.(ElementCustomizeAutoDataPersistence)
	if !ok || persistence == nil {
		return nil, nil, NewErrorf(ErrAtomNotExists, "Atom: Atom not exists. name=(%s)", name).AddStack(e)
	}
	return e.elementAtomSpawn(name, nil, current, persistence, t)
}

func (e *ElementLocal) elementAtomSpawn(name string, arg proto.Message, current *ElementImplementation, persistence ElementCustomizeAutoDataPersistence, t *IDTrackerInfo) (*AtomLocal, *IDTracker, *Error) {
	// Element的容器逻辑。
	// Alloc an atomos and try setting.
	atom := newAtomLocal(name, e, current, current.Interface.Config.LogLevel)
	// If not exist, lock and set a new one.
	e.lock.Lock()
	oldAtom, has := e.atoms[name]
	if !has {
		e.atoms[name] = atom
		atom.nameElement = e.names.PushBack(name)
	}
	e.lock.Unlock()
	// If exists and running, release new and return error.
	// 不用担心两个Atom同时创建的问题，因为Atom创建的时候就是AtomSpawning了，除非其中一个在极端短的时间内AtomHalt了
	if has {
		// 如果旧的存在且不再运行，则用新的atom覆盖。
		oldAtom.atomos.mailbox.mutex.Lock()
		if oldAtom.atomos.state > AtomosHalt {
			oldAtom.atomos.mailbox.mutex.Unlock()
			if oldAtom.atomos.state < AtomosStopping {
				tracker := oldAtom.idTracker.NewIDTracker(t)
				return oldAtom, tracker, NewErrorf(ErrAtomExists, "Atom: Atom exists. name=(%s)", name).AddStack(oldAtom, arg)
			} else {
				return nil, nil, NewErrorf(ErrAtomIsStopping, "Atom: Atom is stopping. name=(%s)", name).AddStack(oldAtom, arg)
			}
		}
		// TODO: 验证这种情况下，IDTrackerManager下面还有引用，引用Release的情况。
		oldLock := &oldAtom.atomos.mailbox.mutex
		atom.nameElement = oldAtom.nameElement
		atom.idTracker = oldAtom.idTracker
		*oldAtom = *atom
		oldLock.Unlock()
		atom = oldAtom
	}

	// Atom的Spawn逻辑。
	if err := atom.atomos.start(func() *Error {
		if e.main.runnable.hookAtomSpawning != nil {
			e.main.runnable.hookAtomSpawning(e.atomos.id.Element, name)
		}
		if err := atom.elementAtomSpawn(current, persistence, arg); err != nil {
			return err.AddStack(nil)
		}
		return nil
	}); err != nil {
		e.elementAtomRelease(atom, nil)
		return nil, nil, err.AddStack(nil)
	}
	tracker := atom.idTracker.NewIDTracker(t)
	return atom, tracker, nil
}

func (e *ElementLocal) elementAtomRelease(atom *AtomLocal, tracker *IDTracker) {
	atom.idTracker.Release(tracker)
	if atom.atomos.isNotHalt() {
		return
	}
	if atom.idTracker.RefCount() > 0 {
		return
	}
	e.lock.Lock()
	name := atom.GetName()
	_, has := e.atoms[name]
	if has {
		delete(e.atoms, name)
	} else {
		e.lock.Unlock()
		return
	}
	if atom.nameElement != nil {
		e.names.Remove(atom.nameElement)
		atom.nameElement = nil
	}
	e.lock.Unlock()
	// assert
	if atom.atomos.mailbox.isRunning() {
		sharedLogging.pushFrameworkErrorLog("elementAtomRelease: Mailbox is still running. name=(%s)", name)
	}
}

func (e *ElementLocal) elementAtomStopping(atom *AtomLocal) {
	if atom.idTracker.RefCount() > 0 {
		return
	}
	e.lock.Lock()
	name := atom.GetName()
	_, has := e.atoms[name]
	if has {
		delete(e.atoms, name)
	}
	if atom.nameElement != nil {
		e.names.Remove(atom.nameElement)
		atom.nameElement = nil
	}
	e.lock.Unlock()
	// assert
	if atom.atomos.mailbox.isRunning() {
		sharedLogging.pushFrameworkErrorLog("elementAtomStopping: Mailbox is still running. name=(%s)", name)
	}

	if StoppingPrintStatic {
		e.Log().Warn("Static >> AtomStopping IDTracker=(%v)", atom.idTracker)
		e.Log().Warn("Static >> AtomStopping MessageTracker=(%v)", atom.GetMessagingInfo())
	}
}

func (e *ElementLocal) cosmosElementSpawn(runnable *CosmosRunnable, current *ElementImplementation) (err *Error) {
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				err = NewErrorf(ErrFrameworkPanic, "Element: Spawn recovers from panic.").AddPanicStack(e, 3, r)
				if ar, ok := e.atomos.instance.(AtomosRecover); ok {
					defer func() {
						recover()
						e.Log().Fatal("Element: Spawn recovers from panic. err=(%v)", err)
					}()
					ar.SpawnRecover(nil, err)
				} else {
					e.Log().Fatal("Element: Spawn recovers from panic. err=(%v)", err)
				}
			}
		}
	}()
	// Get data and Spawning.
	var data proto.Message
	// 尝试进行自动数据持久化逻辑，如果支持的话，就会被执行。
	// 会从对象中GetAtomData，如果返回错误，证明服务不可用，那将会拒绝Atom的Spawn。
	// 如果GetAtomData拿不出数据，且Spawn没有传入参数，则认为是没有对第一次Spawn的Atom传入参数，属于错误。
	pa, ok := current.Developer.(ElementCustomizeAutoLoadPersistence)
	if ok && pa != nil {
		if err = pa.Load(e, runnable.config.Customize); err != nil {
			return err.AddStack(e)
		}
	}
	persistence, ok := current.Developer.(ElementCustomizeAutoDataPersistence)
	if ok && persistence != nil {
		elemPersistence := persistence.ElementAutoDataPersistence()
		if elemPersistence != nil {
			data, err = elemPersistence.GetElementData()
			if err != nil {
				return err
			}
		}
	}
	if err := current.Interface.ElementSpawner(e, e.atomos.instance, data); err != nil {
		return err
	}
	return nil
}
