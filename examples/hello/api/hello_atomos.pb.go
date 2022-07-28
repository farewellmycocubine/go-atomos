// Code generated by protoc-gen-go-atomos. DO NOT EDIT.

package api

import (
	go_atomos "github.com/hwangtou/go-atomos"
	proto "google.golang.org/protobuf/proto"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the atomos package it is being compiled against.

//////
//// INTERFACES
//

////////////////////////////////////
////////// Element: Hello //////////
////////////////////////////////////
//
// The greeting service definition.
// New line
//

const HelloName = "Hello"

// HelloElementID is the interface of Hello element.

type HelloElementID interface {
	go_atomos.ID

	// Sends a greeting
	SayHello(from go_atomos.ID, in *HelloReq) (*HelloResp, *go_atomos.ErrorInfo)
}

func GetHelloElementID(c go_atomos.CosmosNode) (HelloElementID, *go_atomos.ErrorInfo) {
	ca, err := c.GetElementID(HelloName)
	if err != nil {
		return nil, err
	}
	return &helloElementID{ca}, nil
}

// HelloAtomID is the interface of Hello atom.

type HelloAtomID interface {
	go_atomos.ID

	// Sends a greeting
	SayHello(from go_atomos.ID, in *HelloReq) (*HelloResp, *go_atomos.ErrorInfo)

	// Build net
	BuildNet(from go_atomos.ID, in *BuildNetReq) (*BuildNetResp, *go_atomos.ErrorInfo)
}

func GetHelloAtomID(c go_atomos.CosmosNode, name string) (HelloAtomID, *go_atomos.ErrorInfo) {
	ca, err := c.GetElementAtomID(HelloName, name)
	if err != nil {
		return nil, err
	}
	return &helloAtomID{ca}, nil
}

// HelloElement is the atomos implements of Hello element.

type HelloElement interface {
	go_atomos.Atomos
	// Spawn Element
	Spawn(self go_atomos.ElementSelfID, data *HelloData) *go_atomos.ErrorInfo
	// Sends a greeting
	SayHello(from go_atomos.ID, in *HelloReq) (*HelloResp, *go_atomos.ErrorInfo)
}

// HelloAtom is the atomos implements of Hello atom.

type HelloAtom interface {
	go_atomos.Atomos
	// Spawn
	Spawn(self go_atomos.AtomSelfID, arg *HelloSpawnArg, data *HelloData) *go_atomos.ErrorInfo
	// Sends a greeting
	SayHello(from go_atomos.ID, in *HelloReq) (*HelloResp, *go_atomos.ErrorInfo)
	// Build net
	BuildNet(from go_atomos.ID, in *BuildNetReq) (*BuildNetResp, *go_atomos.ErrorInfo)
}

func SpawnHelloAtom(c go_atomos.CosmosNode, name string, arg *HelloSpawnArg) (HelloAtomID, *go_atomos.ErrorInfo) {
	id, err := c.SpawnElementAtom(HelloName, name, arg)
	if err != nil {
		return nil, err
	}
	return &helloAtomID{id}, nil
}

//////
//// IMPLEMENTATIONS
//

////////////////////////////////////
////////// Element: Hello //////////
////////////////////////////////////
//
// The greeting service definition.
// New line
//

type helloElementID struct {
	go_atomos.ID
}

func (c *helloElementID) SayHello(from go_atomos.ID, in *HelloReq) (*HelloResp, *go_atomos.ErrorInfo) {
	r, err := c.Cosmos().MessageElement(from, c, "SayHello", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*HelloResp)
	if !ok {
		return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageReplyType, "Reply type=(%T)", r)
	}
	return reply, err
}

type helloAtomID struct {
	go_atomos.ID
}

func (c *helloAtomID) SayHello(from go_atomos.ID, in *HelloReq) (*HelloResp, *go_atomos.ErrorInfo) {
	r, err := c.Cosmos().MessageAtom(from, c, "SayHello", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*HelloResp)
	if !ok {
		return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageReplyType, "Reply type=(%T)", r)
	}
	return reply, err
}

func (c *helloAtomID) BuildNet(from go_atomos.ID, in *BuildNetReq) (*BuildNetResp, *go_atomos.ErrorInfo) {
	r, err := c.Cosmos().MessageAtom(from, c, "BuildNet", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*BuildNetResp)
	if !ok {
		return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageReplyType, "Reply type=(%T)", r)
	}
	return reply, err
}

func GetHelloInterface(dev go_atomos.ElementDeveloper) *go_atomos.ElementInterface {
	elem := go_atomos.NewInterfaceFromDeveloper(HelloName, dev)
	elem.ElementSpawner = func(s go_atomos.ElementSelfID, a go_atomos.Atomos, data proto.Message) *go_atomos.ErrorInfo {
		dataT, _ := data.(*HelloData)
		elem, ok := a.(HelloElement)
		if !ok {
			return go_atomos.NewErrorf(go_atomos.ErrElementNotImplemented, "Element not implemented, type=(HelloElement)")
		}
		return elem.Spawn(s, dataT)
	}
	elem.AtomSpawner = func(s go_atomos.AtomSelfID, a go_atomos.Atomos, arg, data proto.Message) *go_atomos.ErrorInfo {
		argT, _ := arg.(*HelloSpawnArg)
		dataT, _ := data.(*HelloData)
		atom, ok := a.(HelloAtom)
		if !ok {
			return go_atomos.NewErrorf(go_atomos.ErrAtomNotImplemented, "Atom not implemented, type=(HelloAtom)")
		}
		return atom.Spawn(s, argT, dataT)
	}
	return elem
}

func GetHelloImplement(dev go_atomos.ElementDeveloper) *go_atomos.ElementImplementation {
	elem := go_atomos.NewImplementationFromDeveloper(dev)
	elem.Interface = GetHelloInterface(dev)
	elem.ElementHandlers = map[string]go_atomos.MessageHandler{
		"SayHello": func(from go_atomos.ID, to go_atomos.Atomos, in proto.Message) (proto.Message, *go_atomos.ErrorInfo) {
			req, ok := in.(*HelloReq)
			if !ok {
				return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageArgType, "Arg type=(%T)", in)
			}
			a, ok := to.(HelloElement)
			if !ok {
				return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageAtomType, "Atom type=(%T)", to)
			}
			return a.SayHello(from, req)
		},
	}
	elem.AtomHandlers = map[string]go_atomos.MessageHandler{
		"SayHello": func(from go_atomos.ID, to go_atomos.Atomos, in proto.Message) (proto.Message, *go_atomos.ErrorInfo) {
			req, ok := in.(*HelloReq)
			if !ok {
				return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageArgType, "Arg type=(%T)", in)
			}
			a, ok := to.(HelloAtom)
			if !ok {
				return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageAtomType, "Atom type=(%T)", to)
			}
			return a.SayHello(from, req)
		},
		"BuildNet": func(from go_atomos.ID, to go_atomos.Atomos, in proto.Message) (proto.Message, *go_atomos.ErrorInfo) {
			req, ok := in.(*BuildNetReq)
			if !ok {
				return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageArgType, "Arg type=(%T)", in)
			}
			a, ok := to.(HelloAtom)
			if !ok {
				return nil, go_atomos.NewErrorf(go_atomos.ErrAtomMessageAtomType, "Atom type=(%T)", to)
			}
			return a.BuildNet(from, req)
		},
	}
	elem.ElementDecoders = map[string]*go_atomos.IOMessageDecoder{
		"SayHello": {
			InDec: func(b []byte, p bool) (proto.Message, *go_atomos.ErrorInfo) {
				return go_atomos.MessageUnmarshal(b, &HelloReq{}, p)
			},
			OutDec: func(b []byte, p bool) (proto.Message, *go_atomos.ErrorInfo) {
				return go_atomos.MessageUnmarshal(b, &HelloResp{}, p)
			},
		},
	}
	elem.AtomDecoders = map[string]*go_atomos.IOMessageDecoder{
		"SayHello": {
			InDec: func(b []byte, p bool) (proto.Message, *go_atomos.ErrorInfo) {
				return go_atomos.MessageUnmarshal(b, &HelloReq{}, p)
			},
			OutDec: func(b []byte, p bool) (proto.Message, *go_atomos.ErrorInfo) {
				return go_atomos.MessageUnmarshal(b, &HelloResp{}, p)
			},
		},
		"BuildNet": {
			InDec: func(b []byte, p bool) (proto.Message, *go_atomos.ErrorInfo) {
				return go_atomos.MessageUnmarshal(b, &BuildNetReq{}, p)
			},
			OutDec: func(b []byte, p bool) (proto.Message, *go_atomos.ErrorInfo) {
				return go_atomos.MessageUnmarshal(b, &BuildNetResp{}, p)
			},
		},
	}
	return elem
}
