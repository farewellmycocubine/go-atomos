package go_atomos

// CHECKED!

import (
	"google.golang.org/protobuf/proto"

	"github.com/hwangtou/go-atomos/core"
)

//
// Atom
//

const RunnableName = "AtomosRunnable"

//// 可以被重载的Atom类型
//// Reloadable Atom type.
//type AtomReloadable interface {
//	// 旧的Atom实例在Reload前被调用，通知Atom准备重载。
//	// Calls the old Atom before it is going to be reloaded.
//	WillReload()
//
//	// 新的Atom实例在Reload后被调用，通知Atom已经重载。
//	// Calls the new Atom after it has already been reloaded.
//	DoReload()
//}
//
//// 有状态的Atom
//// 充分发挥protobuf的优势，在Atom非运行时，可以能够保存和恢复所有数据，所以千万不要把数据放在闭包中。
////
//// Stateful Atom
//// Take advantage of protobuf, Atom can save and recovery all data of itself, so never to take data
//// reference into closure.
//type AtomStateful interface {
//	Atom
//	proto.Message
//}
//
//// 无状态的Atom
//// 这种Atom没有需要保存和恢复的东西。
////
//// Stateless Atom
//// Such a type of Atom, is no need to save or recovery anything.
//type AtomStateless interface {
//	Atom
//}

// 暴露给Atom开发者使用的Atom接口。
// Some methods of Atom interface that expose Atom developers to use.

//
// Id
//

// ID 是Atom的类似句柄的对象。
// ID, an instance that similar to file descriptor of the Atom.
type ID interface {
	core.ID

	getCallChain() []ID

	// Release
	// 释放Id的引用计数
	// Release reference count of ID.
	// TODO:思考是否真的需要Release
	Release()

	// Cosmos
	// Atom所在Cosmos节点。
	// Cosmos Node of the Atom.
	Cosmos() CosmosNode

	// Element
	// Atom所属的Element类型。
	// Element type of the Atom.
	Element() Element

	// GetName
	// Atom的名称。
	// Name of the Atom.
	GetName() string

	// GetVersion
	// ElementInterface的版本。
	// Version of ElementInterface.
	GetVersion() uint64

	// Kill
	// 从其它Atom或者main发送Kill消息。
	// write Kill signal from other Atom or main.
	Kill(from ID) *core.ErrorInfo

	//// 内部使用，如果是本地Atom，会返回本地Atom的引用。
	//// Inner use only, if Atom is local, it returns the local AtomCore reference.
	//getLocalAtom() *AtomCore

	String() string
}

type CallProtoBuffer interface {
	// CallNameWithProtoBuffer
	// 直接接收调用
	CallNameWithProtoBuffer(name string, buf []byte) ([]byte, *core.ErrorInfo)
}

type CallJson interface {
	// CallNameWithJson
	// 直接接收调用
	CallNameWithJson(name string, buf []byte) ([]byte, error)
}

//
// AtomSelf
//

// AtomSelf
// 是Atom内部可以访问的Atom资源的概念。
// 通过AtomSelf，Atom内部可以访问到自己的Cosmos（CosmosProcess）、可以杀掉自己（KillSelf），以及提供Log和Task的相关功能。
//
// AtomSelf, a concept that provide Atom resource access to inner Atom.
// With AtomSelf, Atom can access its self-cosmos with "CosmosProcess", can kill itself use "KillSelf" from inner.
// It also provides Log and Tasks method to inner Atom.
type AtomSelf interface {
	ID

	// CosmosProcess
	// 获取Atom的CosmosProcess。
	// Access to the CosmosProcess of the Atom.
	CosmosSelf() *CosmosProcess

	// TODO
	ElementSelf() Element

	// KillSelf
	// Atom从内部杀死自己。
	// Atom kills itself from inner.
	KillSelf()

	// Log
	// Atom日志。
	// Atom Logs.
	Log() core.Logging

	// Task
	// Atom任务
	// Atom Tasks.
	Task() core.Task
}

type ParallelSelf interface {
	ID
	CosmosSelf() *CosmosProcess
	KillSelf()
	Log() core.Logging
}

type ParallelFn func(self ParallelSelf, message proto.Message, id ...ID)

////
//// Wormhole
////
//
//// WormholeAtom
//// 支持WormholeAtom的Atom，可以得到Wormhole的支持。
//// Implement WormholeAtom interface to gain wormhole support.
//type WormholeAtom interface {
//	Atomos
//	AcceptWorm(control WormholeControl) error
//	CloseWorm(control WormholeControl)
//}

//// WormholeId
//// 是Id接口的延伸，提供向WormholeAtom发送Wormhole的可能。
//// Extend of Id, it lets send wormhole to WormholeAtom become possible.
//type WormholeId interface {
//	ID
//	Accept(daemon WormholeDaemon) error
//}
//
//// WormholeDaemon
//// 通常包装着wormhole（真实网络连接）。负责接受信息并处理，并提供操作接口。
//// WormholeDaemon generally used to wrap the real connection. It handles message processing,
//// and provides operating methods.
//type WormholeDaemon interface {
//	// Daemon
//	// 加载&卸载
//	// Loaded & Unloaded
//	Daemon(AtomSelf) error
//	WormholeControl
//}
//
//// WormholeControl
//// 向WormholeAtom提供发送和关闭接口。
//// WormholeControl provides Send and Close to WormholeAtom.
//type WormholeControl interface {
//	Send([]byte) error
//	Close(isKickByNew bool) error
//}
