package atomos

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var ErrAtomTypeNotExists = errors.New("atomos: to type not exists")
var ErrAtomNameExists = errors.New("atomos: to name exists")
var ErrAtomNameNotExists = errors.New("atomos: to name not exists")
var ErrAtomCallNotExists = errors.New("atomos: to call not exists")
var ErrAtomHalt = errors.New("atomos: to halt")
var ErrCustomizeAtomType = errors.New("atomos: customize to type illegal")

// Atom
// Atomos means atom. Atomos refer in particular to framework name, and Atom is a concrete type here.
// Atomos就是Atom。这里的Atomos更多是指这个框架的名字，而Atom指的是具体的类型。
// Developer implements the interface "Atom" and makes it spawn .
// 开发者实现Atom接口，然后用Cosmos来让它运转。

type Atom interface {
	// 在Atom启动的过程中被调用，返回error可让Atom启动失败。
	//   self: Atom内部使用的对象，提供任务、日志和关闭Atom的接口。
	//   buf: Atom上一次运行结束时保存的字节流，可供恢复Atom状态之用。
	Spawn(self AtomSelf, buf []byte) error

	// 在Atom关闭的过程中被调用，返回Atom业务想要保存的状态的字节流。
	//   task: 还未执行的任务id和任务参数。
	Close(tasks map[uint64]proto.Message) []byte
}

type AtomSelf interface {
	Id
	SelfCosmos() *Cosmos
	Halt()
	Log() *AtomLog
	Task() *AtomTask
}

// Id is an interface generated by protoc-gen-go-atomos. Developer can call to with it.
// Id是一个由protoc-gen-go-atomos工具生成的接口。开发者可以使用它来调用atom。

type Id interface {
	Cosmos() CosmosNode
	Type() string
	Name() string
	// Atom可调用接口的通用关闭函数，返回错误则目标Atom拒绝关闭。
	//   from: 传入Spawn调用时的self参数，用于给目标Atom验证。
	Kill(from Id) error
}

// Cosmos节点需要支持的接口内容
// 仅供生成器内部使用

type CosmosNode interface {
	// 获得某个Atom类型的Atom的引用
	GetAtomId(desc *AtomTypeDesc, name string) (Id, error)

	// 调用某个Atom类型的Atom的引用
	CallAtom(from Id, aType, aName, cName string, args proto.Message) (reply proto.Message, err error)

	// 关闭某个Atom类型的Atom
	CloseAtom(from Id, aType, aName string) error
}

//type BoostrapFn func() error todo

// Cosmos
// Cosmos is universe, universe is the container of atoms
// 简单地说，Cosmos就是宇宙，宇宙就是原子的容器

type Cosmos struct {
	mu   sync.Mutex
	conf CosmosConfig
	ds   *dataStorage
	ats  atomTypes
	n    *nodes
	url  string
}

func NewCosmos(conf CosmosConfig, atomTypeDesc ...*AtomTypeDesc) (*Cosmos, error) {
	c := &Cosmos{
		mu:   sync.Mutex{},
		conf: conf,
		ats:  map[string]*atomType{},
		ds:   &dataStorage{},
		n:    &nodes{},
	}
	if err := c.ds.open(conf.DataFilePath); err != nil {
		return nil, err
	}
	if err := c.n.init(c, atomTypeDesc...); err != nil {
		return nil, err
	}
	return c, nil
}

// Local

func (c *Cosmos) Register(desc *AtomTypeDesc) (err error) {
	c.ats.addType(desc)
	c.n.toUpdate = true
	return nil
}

func (c *Cosmos) Run() {
	// Signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// Command
	cmdCh := make(chan string)
	// Loop
	exit := false
	for !exit {
		select {
		case s := <-sigCh:
			log.Println("Exit:", s)
			exit = true
		case c := <-cmdCh:
			log.Println("Cmd:", c)
		}
	}
}

// todo
//func (c *Cosmos) loadPlugin(path string) {
//	p, err := plugin.Open(path)
//	if err != nil {
//		log.Println("Cannot open plugin")
//		return
//	}
//	boostrap, err := p.Lookup("Boostrap")
//	if err != nil {
//		log.Println("Cannot find boostrap function in plugin")
//		return
//	}
//	fn, ok := boostrap.(BoostrapFn)
//	if !ok {
//		log.Println("Illegal boostrap in plugin")
//		return
//	}
//	if err = fn(); err != nil {
//		log.Println("Execute boostrap error,", err)
//		return
//	}
//}

func (c *Cosmos) Close() {
	log.Println("closing cosmos")
	c.n.etcdUnregisterNode(c.url)
	// todo stop context
}

func (c *Cosmos) GetAtomId(desc *AtomTypeDesc, name string) (Id, error) {
	t := c.ats.getType(desc.Name)
	if t == nil {
		return nil, ErrAtomTypeNotExists
	}
	if !t.has(name) {
		return nil, ErrAtomNameNotExists
	}
	return desc.NewId(c, name), nil
}

// Consider about hot reload
func (c *Cosmos) SpawnAtom(desc *AtomTypeDesc, name string, atom Atom) (Id, error) {
	log.Printf("SpawnAtom(%q)", desc.Name)
	c.mu.Lock()
	defer c.mu.Unlock()

	t := c.ats.getType(desc.Name)
	if t == nil {
		return nil, ErrAtomTypeNotExists
	}
	if t.has(name) {
		return nil, ErrAtomNameExists
	}

	aw, err := t.init(c, name, atom)
	if err != nil {
		return nil, err
	}
	go aw.loop()
	return desc.NewId(c, name), nil
}

func (c *Cosmos) CallAtom(from Id, aType, aName, cName string, args proto.Message) (reply proto.Message, err error) {
	typeInfo, has := c.ats[aType]
	if !has {
		return nil, ErrAtomTypeNotExists
	}
	return typeInfo.call(from, aName, cName, args)
}

// todo delete from cosmos
func (c *Cosmos) CloseAtom(from Id, aType, aName string) error {
	typeInfo, has := c.ats[aType]
	if !has {
		return ErrAtomTypeNotExists
	}
	return typeInfo.close(from, aName)
}

// Remote Cosmos

func (c *Cosmos) GetRemoteCosmos(node string) *remoteNode {
	return c.n.remote[node]
}

func (r *remoteNode) GetAtomId(desc *AtomTypeDesc, name string) (Id, error) {
	log.Println("remoteNode.GetAtomId")
	_, has := r.Types[desc.Name]
	if !has {
		return nil, ErrAtomTypeNotExists
	}
	// todo find remote actor
	exist, err := r.sendGetAtomReq(r.CosmosData.Node, desc.Name, name)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrAtomNameNotExists
	}
	return desc.NewId(r, name), nil
}

func (r *remoteNode) CallAtom(from Id, aType, aName, cName string, args proto.Message) (reply proto.Message, err error) {
	// todo send remote actor
	return nil, nil
}

// Inner

func (c *Cosmos) getCosmosData() *CosmosData {
	d := &CosmosData{
		CosmosId: c.conf.CosmosId,
		Node:     c.conf.NodeId,
		Network:  c.conf.Cluster.ListenNetwork,
		Addr:     c.conf.Cluster.ListenAddress,
		Types:    map[string]*AtomType{},
	}
	for typeName, typeInfo := range c.ats {
		var f []string
		for funcName, _ := range typeInfo.calls {
			f = append(f, funcName)
		}
		d.Types[typeName] = &AtomType{
			Name: typeName,
			Func: f,
		}
	}
	return d
}
