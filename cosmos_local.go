package go_atomos

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"sync"
)

const (
	lCLICreateElementFailed     = "CosmosLocal.watcherInit: Create element failed, nodeName=%s,err=%s"
	lCLILoadElementFailed       = "CosmosLocal.watcherInit: Load element failed, nodeName=%s,err=%s"
	lCLILoadElement             = "CosmosLocal.watcherInit: Load element succeed, nodeName=%s"
	lCLILoadElementFailedUnload = "CosmosLocal.watcherInit: Unload loaded element, nodeName=%s,err=%v"
	lCLIHasInitialized          = "CosmosLocal.watcherInit: Has initialized"
	lCLIClusterInitFailed       = "CosmosLocal.watcherInit: Cluster watcherInit failed"
	lCLRCosmosRunning           = "CosmosLocal.runRunnable: Cosmos running"
	lCLRCosmosInvalid           = "CosmosLocal.runRunnable: Cosmos is invalid"
	lCLCInfo                    = "CosmosLocal.close: Unload element, nodeName=%s,err=%v"
	lCLEInfo                    = "CosmosLocal.exitRunnable: Exiting"
	lCLGENotFound               = "CosmosLocal.getElement: Element not found, nodeName=%s"
	lCLAElementExists           = "CosmosLocal.addElement: Element exists, nodeName=%s"
	lCLAElementLoadFailed       = "CosmosLocal.addElement: Element load failed, nodeName=%s,err=%v"
	lCLDElementNotFound         = "CosmosLocal.delElement: Element not found, nodeName=%s"
	lCLDElementUnloadFailed     = "CosmosLocal.delElement: Element unload failed, nodeName=%s,err=%v"
)

type CosmosLocal struct {
	mutex             sync.RWMutex
	config            *Config
	cosmos            *CosmosSelf
	mainId            *idMain
	elements          map[string]*ElementLocal
	killNoticeCh      chan bool
	elementInterfaces map[string]*ElementInterface
}

// Life cycle

func newCosmosLocal() *CosmosLocal {
	return &CosmosLocal{}
}

// 初始化Runnable。
// Initial Runnable.
func (c *CosmosLocal) initRunnable(self *CosmosSelf, runnable CosmosRunnable) error {
	if err := runnable.Check(); err != nil {
		return err
	}
	// 初始化Runnable中所有的本地Elements。
	// Initial all locals elements in the Runnable. TODO: Parallels loading.
	loaded := map[string]*ElementLocal{}
	elements := make(map[string]*ElementLocal, len(runnable.implementations))
	for name, define := range runnable.implementations {
		// Create the element.
		elem, err := newElementLocal(self, define)
		if err != nil {
			err = fmt.Errorf(lCLICreateElementFailed, name, err)
			self.logFatal("%s", err.Error())
			return err
		}
		// Load the element.
		if err = elem.load(); err != nil {
			err = fmt.Errorf(lCLILoadElementFailed, name, err)
			self.logFatal("%s", err.Error())
			for loadedName, loadedElem := range loaded {
				err = loadedElem.unload()
				err = fmt.Errorf(lCLILoadElementFailedUnload, loadedName, err)
			}
			return err
		}
		// Add the element.
		elements[name] = elem
		self.logInfo(lCLILoadElement, name)
	}
	elementInterfaces := make(map[string]*ElementInterface)
	for name, i := range runnable.interfaces {
		elementInterfaces[name] = i
	}

	// Lock, set elements, and unlock.
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.config != nil || c.elements != nil {
		err := fmt.Errorf(lCLIHasInitialized)
		self.logFatal("%s", err.Error())
		return err
	}
	c.config = self.config
	c.cosmos = self
	c.elements = elements
	c.elementInterfaces = elementInterfaces
	c.mainId = &idMain{
		cosmosLocal: c,
	}
	c.killNoticeCh = make(chan bool)
	if err := self.remotes.init(); err != nil {
		err = fmt.Errorf(lCLIClusterInitFailed)
		self.logFatal("%s", err.Error())
		return err
	}

	return nil
}

// 执行Runnable。
// Run runnable.
func (c *CosmosLocal) runRunnable(runnable CosmosRunnable) error {
	c.mutex.Lock()
	cosmos := c.cosmos
	c.mutex.Unlock()
	if cosmos == nil {
		err := fmt.Errorf(lCLRCosmosInvalid)
		c.cosmos.logFatal("%s", err.Error())
		return err
	}
	c.cosmos.logInfo(lCLRCosmosRunning)
	runnable.script(cosmos, c.mainId, c.killNoticeCh)
	return nil
}

// 退出Runnable。
// Exit runnable.
func (c *CosmosLocal) exitRunnable() {
	c.cosmos.logInfo(lCLEInfo)
	// todo: remove all actors of elements, then remove all element. Save actor data and warning if actor still run.
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cosmos.remotes.close()
	c.close()
}

func (c *CosmosLocal) close() {
	// TODO: Parallels unloading
	for elemName, elem := range c.elements {
		err := elem.unload()
		c.cosmos.logInfo(lCLCInfo, elemName, err)
		delete(c.elements, elemName)
	}
	c.config = nil
	c.cosmos = nil
	c.mainId = nil
	c.elements = nil
	c.elementInterfaces = nil
}

// Element

func (c *CosmosLocal) getElement(name string) (*ElementLocal, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	elem, has := c.elements[name]
	if !has {
		err := fmt.Errorf(lCLGENotFound, name)
		c.cosmos.logFatal("%s", err.Error())
		return nil, err
	}
	return elem, nil
}

func (c *CosmosLocal) addElement(name string, elem *ElementLocal) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check exists
	if _, has := c.elements[name]; has {
		err := fmt.Errorf(lCLAElementExists, name)
		c.cosmos.logFatal("%s", err.Error())
		return err
	}
	// Try load and set
	if err := elem.load(); err != nil {
		err := fmt.Errorf(lCLAElementLoadFailed, name, err)
		c.cosmos.logFatal("%s", err.Error())
		return err
	}
	c.elements[name] = elem
	inf := fmt.Sprintf(lCLILoadElement, name)
	c.cosmos.logInfo("%s", inf)
	return nil
}

func (c *CosmosLocal) delElement(name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check exists
	elem, has := c.elements[name]
	if !has {
		err := fmt.Errorf(lCLDElementNotFound, name)
		c.cosmos.logFatal("%s", err.Error())
		return err
	}
	// Try unload and unset
	if err := elem.unload(); err != nil {
		err := fmt.Errorf(lCLDElementUnloadFailed, name, err)
		c.cosmos.logFatal("%s", err.Error())
		return err
	}
	delete(c.elements, name)
	return nil
}

// Atom

func (c *CosmosLocal) GetNodeName() string {
	return c.config.Node
}

func (c *CosmosLocal) IsLocal() bool {
	return true
}

func (c *CosmosLocal) GetAtomId(elemName, atomName string) (Id, error) {
	// Get element.
	e, err := c.getElement(elemName)
	if err != nil {
		return nil, err
	}
	// Get atom.
	return e.GetAtomId(atomName)
}

func (c *CosmosLocal) SpawnAtom(elemName, atomName string, arg proto.Message) (Id, error) {
	// Get element.
	e, err := c.getElement(elemName)
	if err != nil {
		return nil, err
	}
	// Try spawning.
	i, err := e.SpawnAtom(atomName, arg)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (c *CosmosLocal) MessageAtom(fromId, toId Id, message string, args proto.Message) (reply proto.Message, err error) {
	return toId.Element().MessagingAtom(fromId, toId, message, args)
}

func (c *CosmosLocal) KillAtom(fromId, toId Id) error {
	return toId.Element().KillAtom(fromId, toId)
}

// idMain, uses in Script entrance.

type idMain struct {
	cosmosLocal *CosmosLocal
}

func (c *CosmosLocal) initIdMain() MainId {
	id := &idMain{
		cosmosLocal: c,
	}
	return id
}

func (i *idMain) Cosmos() CosmosNode {
	return i.cosmosLocal
}

func (i *idMain) Element() Element {
	return nil
}

func (i *idMain) Name() string {
	return "main"
}

func (i *idMain) Version() uint64 {
	return 0
}

func (i *idMain) Kill(from Id) error {
	return ErrAtomCannotKill
}

func (i *idMain) getLocalAtom() *AtomCore {
	return nil
}
