package go_atomos

import (
	"google.golang.org/protobuf/proto"
	"net/http"
	"testing"
)

//func TestAtomosEnvironment(t *testing.T) {
//	testInitWebSocketConnsManager(t)
//	defer testExitWebSocketConnsManager(t)
//
//	e, err := GetTestElementID(self.Cosmos())
//	if err != nil {
//		t.Fatal(err)
//	}
//	o, err := e.Message(self, &String{S: "test"})
//	if err != nil {
//		t.Fatal(err)
//	}
//	if o == nil || len(o.Ss) == 0 || o.Ss[0] != "test" {
//		t.Fatal("test fail")
//	}
//}

const (
	testCosmos = "test_cosmos"
	testNode1  = "test_node_1"
)

var s script
var self SelfID

func testInitWebSocketConnsManager(t *testing.T) {
	testRunnable := CosmosRunnable{}
	testRunnable.
		AddElementImplementation(GetTestImplement(&tDev{})).SetElementSpawn(TestName).
		SetConfig(&Config{Cosmos: testCosmos, Node: testNode1, LogLevel: LogLevel_Debug}).
		SetMainScript(&s)
	MainForTest(testRunnable, t)
}

func testExitWebSocketConnsManager(t *testing.T) {
	if err := SharedCosmosProcess().Stop(); err != nil {
		t.Fatal(err)
	}
}

type script struct{}

func (m *script) OnBoot(local *CosmosProcess) *Error {
	local.Self().Parallel(func() {
		if er := http.ListenAndServe(":8081", nil); er != nil {
			local.Self().Log().Info("OnBoot: ListenAndServe error. err=(%v)", er)
		}
	})
	return nil
}

func (m *script) OnStartUp(local *CosmosProcess) *Error {
	self = local.Self()
	return nil
}

func (m *script) OnShutdown() *Error {
	return nil
}

type tDev struct{}

func (t *tDev) ElementConstructor() Atomos {
	return newTElem()
}

func (t *tDev) AtomConstructor(name string) Atomos {
	return newTAtom()
}

type tElem struct {
	self ElementSelfID
}

func newTElem() TestElement {
	return &tElem{}
}

func (t *tElem) String() string {
	return t.self.String()
}

func (t *tElem) Spawn(self ElementSelfID, data *Nil) *Error {
	t.self = self
	return nil
}

func (t *tElem) Halt(from ID, cancelled []uint64) (save bool, data proto.Message) {
	return false, nil
}

func (t *tElem) Message(from ID, in *String) (out *Strings, err *Error) {
	return &Strings{Ss: []string{in.S}}, nil
}

func (t *tElem) ScaleMessage(from ID, in *String) (*TestAtomID, *Error) {
	return SpawnTestAtom(self, from.Cosmos(), in.S, &Nil{})
}

type tAtom struct {
	self AtomSelfID
}

func newTAtom() TestAtom {
	return &tAtom{}
}

func (t *tAtom) String() string {
	return t.self.String()
}

func (t *tAtom) Spawn(self AtomSelfID, arg *Nil, data *Nil) *Error {
	t.self = self
	return nil
}

func (t *tAtom) Halt(from ID, cancelled []uint64) (save bool, data proto.Message) {
	return false, nil
}

func (t *tAtom) AtomMessage(from ID, in *String) (out *Strings, err *Error) {
	return &Strings{Ss: []string{in.S}}, nil
}

func (t *tAtom) ScaleMessage(from ID, in *String) (out *Strings, err *Error) {
	return &Strings{Ss: []string{in.S}}, nil
}

// Code generated by protoc-gen-go-atomos. DO NOT EDIT.

const TestName = "Test"

////////////////////////////////////
/////////// 需要实现的接口 ///////////
////// Interface to implement //////
////////////////////////////////////

// TestElement is the atomos implements of Test element.

type TestElement interface {
	Atomos
	Spawn(self ElementSelfID, data *Nil) *Error

	Message(from ID, in *String) (out *Strings, err *Error)

	// Scale Methods

	// Scale
	ScaleMessage(from ID, in *String) (*TestAtomID, *Error)
}

// TestAtom is the atomos implements of Test atom.

type TestAtom interface {
	Atomos

	// Atom
	Spawn(self AtomSelfID, arg *Nil, data *Nil) *Error

	AtomMessage(from ID, in *String) (out *Strings, err *Error)

	// Scale Methods

	// Scale
	ScaleMessage(from ID, in *String) (out *Strings, err *Error)
}

////////////////////////////////////
/////////////// 识别符 //////////////
//////////////// ID ////////////////
////////////////////////////////////

// Element: Test

type TestElementID struct {
	ID
	*IDTracker
}

// 获取某节点中的ElementID
// Get element id of node
func GetTestElementID(c CosmosNode) (*TestElementID, *Error) {
	ca, err := c.CosmosGetElementID(TestName)
	if err != nil {
		return nil, err
	}
	return &TestElementID{ca, nil}, nil
}

// Sync
func (c *TestElementID) Message(callerID SelfID, in *String, ext ...interface{}) (out *Strings, err *Error) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testElementValue.Message }
	return testElementMessengerValue.Message().SyncElement(c, callerID, in, ext...)
}

// Async
func (c *TestElementID) AsyncMessage(callerID SelfID, in *String, callback func(out *Strings, err *Error), ext ...interface{}) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testElementValue.Message }
	testElementMessengerValue.Message().AsyncElement(c, callerID, in, callback, ext...)
}

// GetID
func (c *TestElementID) ScaleMessageGetID(callerID SelfID, in *String, ext ...interface{}) (id *TestAtomID, err *Error) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testElementValue.ScaleMessage }
	i, tracker, err := testElementMessengerValue.ScaleMessage().GetScaleID(c, callerID, TestName, in, ext...)
	if err != nil {
		return nil, err.AddStack(nil)
	}
	return &TestAtomID{i, tracker}, nil
}

// Sync
func (c *TestElementID) ScaleMessage(callerID SelfID, in *String, ext ...interface{}) (out *Strings, err *Error) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testElementValue.ScaleMessage }
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testAtomValue.ScaleMessage }
	id, err := c.ScaleMessageGetID(callerID, in, ext...)
	if err != nil {
		return nil, err.AddStack(nil)
	}
	defer id.Release()
	return testAtomMessengerValue.ScaleMessage().SyncAtom(id, callerID, in)
}

// Async
func (c *TestElementID) ScaleAsyncMessage(callerID SelfID, in *String, callback func(*Strings, *Error), ext ...interface{}) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testElementValue.ScaleMessage }
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testAtomValue.ScaleMessage }
	id, err := c.ScaleMessageGetID(callerID, in, ext...)
	if err != nil {
		callback(nil, err.AddStack(nil))
		return
	}
	defer id.Release()
	testAtomMessengerValue.ScaleMessage().AsyncAtom(id, callerID, in, callback, ext...)
}

// Atom: Test

type TestAtomID struct {
	ID
	*IDTracker
}

// 创建（自旋）某节点中的一个Atom，并返回AtomID
// Create (spin) an atom in a node and return the AtomID
func SpawnTestAtom(caller SelfID, c CosmosNode, name string, arg *Nil) (*TestAtomID, *Error) {
	id, tracker, err := c.CosmosSpawnAtom(caller, TestName, name, arg)
	if id == nil {
		return nil, err.AddStack(nil)
	}
	return &TestAtomID{id, tracker}, err
}

// 获取某节点中的AtomID
// Get atom id of node
func GetTestAtomID(c CosmosNode, name string) (*TestAtomID, *Error) {
	ca, tracker, err := c.CosmosGetAtomID(TestName, name)
	if err != nil {
		return nil, err
	}
	return &TestAtomID{ca, tracker}, nil
}

// Sync
func (c *TestAtomID) ScaleMessage(callerID SelfID, in *String, ext ...interface{}) (out *Strings, err *Error) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testAtomValue.ScaleMessage }
	return testAtomMessengerValue.ScaleMessage().SyncAtom(c, callerID, in, ext...)
}

// Async
func (c *TestAtomID) AsyncScaleMessage(callerID SelfID, in *String, callback func(out *Strings, err *Error), ext ...interface{}) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testAtomValue.ScaleMessage }
	testAtomMessengerValue.ScaleMessage().AsyncAtom(c, callerID, in, callback, ext...)
}

// Sync
func (c *TestAtomID) AtomMessage(callerID SelfID, in *String, ext ...interface{}) (out *Strings, err *Error) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testAtomValue.AtomMessage }
	return testAtomMessengerValue.AtomMessage().SyncAtom(c, callerID, in, ext...)
}

// Async
func (c *TestAtomID) AsyncAtomMessage(callerID SelfID, in *String, callback func(out *Strings, err *Error), ext ...interface{}) {
	/* CODE JUMPER 代码跳转 */ _ = func() { _ = testAtomValue.AtomMessage }
	testAtomMessengerValue.AtomMessage().AsyncAtom(c, callerID, in, callback, ext...)
}

// Atomos Interface

func GetTestImplement(dev ElementDeveloper) *ElementImplementation {
	elem := NewImplementationFromDeveloper(dev)
	elem.Interface = GetTestInterface(dev)
	elem.ElementHandlers = map[string]MessageHandler{
		"Message": func(from ID, to Atomos, in proto.Message) (proto.Message, *Error) {
			a, i, err := testElementMessengerValue.Message().ExecuteAtom(to, in)
			if err != nil {
				return nil, err.AddStack(nil)
			}
			return a.Message(from, i)
		},
	}
	elem.AtomHandlers = map[string]MessageHandler{
		"ScaleMessage": func(from ID, to Atomos, in proto.Message) (proto.Message, *Error) {
			a, i, err := testAtomMessengerValue.ScaleMessage().ExecuteAtom(to, in)
			if err != nil {
				return nil, err.AddStack(nil)
			}
			return a.ScaleMessage(from, i)
		},
		"AtomMessage": func(from ID, to Atomos, in proto.Message) (proto.Message, *Error) {
			a, i, err := testAtomMessengerValue.AtomMessage().ExecuteAtom(to, in)
			if err != nil {
				return nil, err.AddStack(nil)
			}
			return a.AtomMessage(from, i)
		},
	}
	elem.ScaleHandlers = map[string]ScaleHandler{
		"ScaleMessage": func(from ID, e Atomos, message string, in proto.Message) (id ID, err *Error) {
			a, i, err := testElementMessengerValue.ScaleMessage().ExecuteScale(e, in)
			if err != nil {
				return nil, err.AddStack(nil)
			}
			return a.ScaleMessage(from, i)
		},
	}
	return elem
}
func GetTestInterface(dev ElementDeveloper) *ElementInterface {
	elem := NewInterfaceFromDeveloper(TestName, dev)
	elem.ElementSpawner = func(s ElementSelfID, a Atomos, data proto.Message) *Error {
		dataT, _ := data.(*Nil)
		elem, ok := a.(TestElement)
		if !ok {
			return NewErrorf(ErrElementNotImplemented, "Element not implemented, type=(TestElement)")
		}
		return elem.Spawn(s, dataT)
	}
	elem.AtomSpawner = func(s AtomSelfID, a Atomos, arg, data proto.Message) *Error {
		argT, _ := arg.(*Nil)
		dataT, _ := data.(*Nil)
		atom, ok := a.(TestAtom)
		if !ok {
			return NewErrorf(ErrAtomNotImplemented, "Atom not implemented, type=(TestAtom)")
		}
		return atom.Spawn(s, argT, dataT)
	}
	elem.ElementDecoders = map[string]*IOMessageDecoder{
		"Message":      testElementMessengerValue.Message().Decoder(&String{}, &Strings{}),
		"ScaleMessage": testElementMessengerValue.ScaleMessage().Decoder(&String{}, &Strings{}),
	}
	elem.AtomDecoders = map[string]*IOMessageDecoder{
		"ScaleMessage": testAtomMessengerValue.ScaleMessage().Decoder(&String{}, &Strings{}),
		"AtomMessage":  testAtomMessengerValue.AtomMessage().Decoder(&String{}, &Strings{}),
	}
	return elem
}

// Atomos Internal

// Element Define

type testElementMessenger struct{}

func (m testElementMessenger) Message() Messenger[*TestElementID, *TestAtomID, TestElement, *String, *Strings] {
	return Messenger[*TestElementID, *TestAtomID, TestElement, *String, *Strings]{nil, nil, "Message"}
}
func (m testElementMessenger) ScaleMessage() Messenger[*TestElementID, *TestAtomID, TestElement, *String, *Strings] {
	return Messenger[*TestElementID, *TestAtomID, TestElement, *String, *Strings]{nil, nil, "ScaleMessage"}
}

var testElementMessengerValue testElementMessenger
var testElementValue TestElement

// Atom Define

type testAtomMessenger struct{}

func (m testAtomMessenger) ScaleMessage() Messenger[*TestElementID, *TestAtomID, TestAtom, *String, *Strings] {
	return Messenger[*TestElementID, *TestAtomID, TestAtom, *String, *Strings]{nil, nil, "ScaleMessage"}
}
func (m testAtomMessenger) AtomMessage() Messenger[*TestElementID, *TestAtomID, TestAtom, *String, *Strings] {
	return Messenger[*TestElementID, *TestAtomID, TestAtom, *String, *Strings]{nil, nil, "AtomMessage"}
}

var testAtomMessengerValue testAtomMessenger
var testAtomValue TestAtom