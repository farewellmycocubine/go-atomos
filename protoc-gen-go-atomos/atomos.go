package main

import (
	"errors"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	atomosPackage   = protogen.GoImportPath("github.com/hwangtou/go-atomos")
	protobufPackage = protogen.GoImportPath("google.golang.org/protobuf/proto")
)

// generateFile generates a _grpc.pb.go file containing gRPC service definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_atomos.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-atomos. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	generateFileContent(gen, file, g)
	return g
}

// generateFileContent generates the atom definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) == 0 {
		return
	}

	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the atomos package it is being compiled against.")
	g.P()
	g.P("//////")
	g.P("//// INTERFACES")
	g.P("//")
	g.P()
	g.P()
	for _, service := range file.Services {
		elementName := service.GoName
		if elementName == "Main" {
			gen.Error(errors.New("cannot use element name \"Main\""))
			continue
		}
		elementTitle(g, service)
		// ID
		genElementIDInterface(g, service)
		genAtomIDInterface(g, service)
		// Interface
		genElementInterface(g, service)
		genAtomInterface(g, service)
	}
	g.P()
	g.P("//////")
	g.P("//// IMPLEMENTATIONS")
	g.P("//")
	g.P()
	for _, service := range file.Services {
		elementName := service.GoName
		if elementName == "Main" {
			gen.Error(errors.New("cannot use element name \"Main\""))
			continue
		}
		elementTitle(g, service)
		genElementIdInternal(g, service)
		genAtomIdInternal(g, service)
		genImplement(file, g, service)
	}
}

func elementTitle(g *protogen.GeneratedFile, service *protogen.Service) {
	elementName := service.GoName
	head, tail := "////////// Element: ", " //////////"
	nameLen := len(elementName) + len(head) + len(tail)
	c := strings.Repeat("/", nameLen)
	g.P(c)
	g.P(head, elementName, tail)
	g.P(c)
	if len(service.Comments.Leading.String()) > 0 {
		g.P("//")
		g.P(service.Comments.Leading.String(), "//")
	}
	g.P()
}

func genElementIDInterface(g *protogen.GeneratedFile, service *protogen.Service) {
	elementName := service.GoName
	elementIDName := service.GoName + "ElementID"

	g.P("const ", elementName, "Name = \"", elementName, "\"")
	g.P()
	g.P("// ", elementIDName, " is the interface of ", elementName, " element.")
	g.P()

	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	g.Annotate(elementIDName, service.Location)
	g.P("type ", elementIDName, " interface {")
	g.P(atomosPackage.Ident("ID"))
	for _, method := range service.Methods {
		methodName := method.GoName
		if !strings.HasPrefix(methodName, "Element") {
			continue
		}
		if strings.HasPrefix(methodName, "ElementSpawn") {
			continue
		}
		if strings.HasPrefix(methodName, "Spawn") {
			continue
		}
		methodName = strings.Split(methodName, "Element")[1]
		if methodName == "" {
			continue
		}
		g.Annotate(elementIDName+"."+methodName, method.Location)
		if method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
			g.P(deprecationComment)
		}
		g.P()
		commentLen := len(method.Comments.Leading.String())
		if commentLen > 0 {
			g.P(method.Comments.Leading.String()[:commentLen-1])
		}
		methodSign(g, method, methodName)
		g.P()
	}
	g.P("}")
	g.P()

	// NewClient factory.
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P(deprecationComment)
	}
	g.P("func Get", elementIDName, " (c ", atomosPackage.Ident("CosmosNode"), ") (", elementIDName, ", *", atomosPackage.Ident("ErrorInfo"), ") {")
	g.P("ca, err := c.GetElementID(", elementName, "Name)")
	g.P("if err != nil { return nil, err }")
	g.P("return &", noExport(elementIDName), "{ca}, nil")
	g.P("}")
	g.P()
}

func genAtomIDInterface(g *protogen.GeneratedFile, service *protogen.Service) {
	atomName := service.GoName
	atomNameID := service.GoName + "AtomID"

	g.P("// ", atomNameID, " is the interface of ", atomName, " atom.")
	g.P()

	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	g.Annotate(atomNameID, service.Location)
	g.P("type ", atomNameID, " interface {")
	g.P(atomosPackage.Ident("ID"))
	for _, method := range service.Methods {
		methodName := method.GoName
		if strings.HasPrefix(methodName, "Element") {
			continue
		}
		if strings.HasPrefix(methodName, "Spawn") {
			continue
		}
		if methodName == "" {
			continue
		}
		g.Annotate(atomNameID+"."+methodName, method.Location)
		if method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
			g.P(deprecationComment)
		}
		g.P()
		commentLen := len(method.Comments.Leading.String())
		if commentLen > 0 {
			g.P(method.Comments.Leading.String()[:commentLen-1])
		}
		methodSign(g, method, methodName)
		g.P()
	}
	g.P("}")
	g.P()

	// NewClient factory.
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P(deprecationComment)
	}
	g.P("func Get", atomNameID, " (c ", atomosPackage.Ident("CosmosNode"), ", name string) (", atomNameID, ", *", atomosPackage.Ident("ErrorInfo"), ") {")
	g.P("ca, err := c.GetElementAtomID(", atomName, "Name, name)")
	g.P("if err != nil { return nil, err }")
	g.P("return &", noExport(atomNameID), "{ca}, nil")
	g.P("}")
	g.P()
}

func genElementInterface(g *protogen.GeneratedFile, service *protogen.Service) {
	elementName := service.GoName + "Element"

	// Server struct.
	g.P("// ", elementName, " is the atomos implements of ", service.GoName, " element.")
	g.P()
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	g.Annotate(elementName, service.Location)
	g.P("type ", elementName, " interface {")
	g.P(atomosPackage.Ident("Atomos"))
	hasElementSpawn := false
	for _, method := range service.Methods {
		methodName := method.GoName
		if !strings.HasPrefix(methodName, "Element") {
			continue
		}
		methodName = strings.TrimPrefix(methodName, "Element")
		g.Annotate(elementName+"."+methodName, method.Location)
		if method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
			g.P(deprecationComment)
		}
		commentLen := len(method.Comments.Leading.String())
		if commentLen > 0 {
			g.P(method.Comments.Leading.String()[:commentLen-1])
		}
		if methodName == "Spawn" {
			hasElementSpawn = true
			spawnElementSign(g, method)
		} else {
			methodSign(g, method, methodName)
		}
	}
	if !hasElementSpawn {
		spawnElementDefaultSign(g)
	}
	g.P("}")
	g.P()
}

func genAtomInterface(g *protogen.GeneratedFile, service *protogen.Service) {
	elementName := service.GoName
	atomName := service.GoName + "Atom"

	// Server struct.
	g.P("// ", atomName, " is the atomos implements of ", service.GoName, " atom.")
	g.P()
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	g.Annotate(atomName, service.Location)
	g.P("type ", atomName, " interface {")
	g.P(atomosPackage.Ident("Atomos"))
	var spawnArgTypeName string
	for _, method := range service.Methods {
		methodName := method.GoName
		if strings.HasPrefix(methodName, "Element") {
			continue
		}
		methodName = strings.TrimPrefix(methodName, "Element")
		g.Annotate(atomName+"."+methodName, method.Location)
		if method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
			g.P(deprecationComment)
		}
		commentLen := len(method.Comments.Leading.String())
		if commentLen > 0 {
			g.P(method.Comments.Leading.String()[:commentLen-1])
		}
		if methodName == "Spawn" {
			spawnArgTypeName = g.QualifiedGoIdent(method.Input.GoIdent)
			spawnAtomSign(g, method)
		} else {
			methodSign(g, method, methodName)
		}

	}
	g.P("}")
	g.P()

	idName := atomName + "ID"
	// Spawn
	g.P("func Spawn", atomName, "(c ", atomosPackage.Ident("CosmosNode"),
		", name string, arg *", spawnArgTypeName, ") (",
		idName, ", *", atomosPackage.Ident("ErrorInfo"), ") {")
	g.P("id, err := c.SpawnElementAtom(", elementName, "Name, name, arg)")
	g.P("if err != nil { return nil, err }")
	g.P("return &", noExport(idName), "{id}, nil")
	g.P("}")
}

func genElementIdInternal(g *protogen.GeneratedFile, service *protogen.Service) {
	idName := service.GoName + "ElementID"

	// ID structure.
	g.P("type ", noExport(idName), " struct {")
	g.P(atomosPackage.Ident("ID"))
	g.P("}")
	g.P()

	// Client method implementations.
	for _, method := range service.Methods {
		methodName := method.GoName
		if !strings.HasPrefix(methodName, "Element") {
			continue
		}
		if strings.HasPrefix(methodName, "ElementSpawn") {
			continue
		}
		if strings.HasPrefix(methodName, "Spawn") {
			continue
		}
		methodName = strings.Split(methodName, "Element")[1]
		if methodName == "" {
			continue
		}
		g.P("func (c *", noExport(idName), ") ", methodName+"(from ", atomosPackage.Ident("ID"),
			", in *", g.QualifiedGoIdent(method.Input.GoIdent),
			") (*", g.QualifiedGoIdent(method.Output.GoIdent),
			", *", atomosPackage.Ident("ErrorInfo"), ")", " {")
		g.P("r, err := c.Cosmos().MessageElement(from, c, \"", methodName, "\", in)")
		g.P("if r == nil { return nil, err }")
		g.P("reply, ok := r.(*", method.Output.GoIdent, ")")
		g.P("if !ok { return nil, ", atomosPackage.Ident("NewErrorf("), atomosPackage.Ident("ErrAtomMessageReplyType"), ", \"Reply type=(%T)\", r) }")
		g.P("return reply, err")
		g.P("}")
		g.P()
	}
}

func genAtomIdInternal(g *protogen.GeneratedFile, service *protogen.Service) {
	idName := service.GoName + "AtomID"

	// ID structure.
	g.P("type ", noExport(idName), " struct {")
	g.P(atomosPackage.Ident("ID"))
	g.P("}")
	g.P()

	// Client method implementations.
	for _, method := range service.Methods {
		methodName := method.GoName
		if strings.HasPrefix(methodName, "Element") {
			continue
		}
		if strings.HasPrefix(methodName, "Spawn") {
			continue
		}
		if methodName == "" {
			continue
		}
		g.P("func (c *", noExport(idName), ") ", methodName+"(from ", atomosPackage.Ident("ID"),
			", in *", g.QualifiedGoIdent(method.Input.GoIdent),
			") (*", g.QualifiedGoIdent(method.Output.GoIdent),
			", *", atomosPackage.Ident("ErrorInfo"), ")", " {")
		g.P("r, err := c.Cosmos().MessageAtom(from, c, \"", methodName, "\", in)")
		g.P("if r == nil { return nil, err }")
		g.P("reply, ok := r.(*", method.Output.GoIdent, ")")
		g.P("if !ok { return nil, ", atomosPackage.Ident("NewErrorf("), atomosPackage.Ident("ErrAtomMessageReplyType"), ", \"Reply type=(%T)\", r) }")
		g.P("return reply, err")
		g.P("}")
		g.P()
	}
}

func genImplement(file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	elementName := service.GoName
	elementAtomName := elementName + "Atom"
	elementElementName := elementName + "Element"
	interfaceName := service.GoName + "Interface"

	g.P("func Get", elementName, "Interface(dev ", atomosPackage.Ident("ElementDeveloper"), ") *", atomosPackage.Ident("ElementInterface"), "{")
	g.P("elem := ", atomosPackage.Ident("NewInterfaceFromDeveloper("), elementName, "Name, dev)")
	hasElementSpawn := false
	for _, method := range service.Methods {
		if method.GoName != "ElementSpawn" {
			continue
		}
		hasElementSpawn = true
		g.P("elem.ElementSpawner = func(s ", atomosPackage.Ident("ElementSelfID"), ", a ", atomosPackage.Ident("Atomos"), ", data ", protobufPackage.Ident("Message"), ") *", atomosPackage.Ident("ErrorInfo"), " {")
		g.P("dataT, _ := data.(*", method.Output.GoIdent, ")")
		g.P("elem, ok := a.(", elementElementName, ")")
		g.P("if !ok { return ", atomosPackage.Ident("NewErrorf"), "(", atomosPackage.Ident("ErrElementNotImplemented"), ", \"Element not implemented, type=(", elementElementName, ")\") }")
		g.P("return elem.Spawn(s, dataT)")
		g.P("}")
	}
	if !hasElementSpawn {
		g.P("elem.ElementSpawner = func(s ", atomosPackage.Ident("ElementSelfID"), ", a ", atomosPackage.Ident("Atomos"), ", data ", protobufPackage.Ident("Message"), ") *", atomosPackage.Ident("ErrorInfo"), " {")
		//g.P("dataT, _ := data.(*", method.Output.GoIdent, ")")
		g.P("elem, ok := a.(", elementElementName, ")")
		g.P("if !ok { return ", atomosPackage.Ident("NewErrorf"), "(", atomosPackage.Ident("ErrElementNotImplemented"), ", \"Element not implemented, type=(", elementElementName, ")\") }")
		g.P("return elem.Spawn(s, nil)")
		g.P("}")
	}
	for _, method := range service.Methods {
		if method.GoName != "Spawn" {
			continue
		}
		g.P("elem.AtomSpawner = func(s ", atomosPackage.Ident("AtomSelfID"), ", a ", atomosPackage.Ident("Atomos"), ", arg, data ", protobufPackage.Ident("Message"), ") *", atomosPackage.Ident("ErrorInfo"), " {")
		g.P("argT, _ := arg.(*", method.Input.GoIdent, "); dataT, _ := data.(*", method.Output.GoIdent, ")")
		g.P("atom, ok := a.(", elementAtomName, ")")
		g.P("if !ok { return ", atomosPackage.Ident("NewErrorf"), "(", atomosPackage.Ident("ErrAtomNotImplemented"), ", \"Atom not implemented, type=(", elementAtomName, ")\") }")
		g.P("return atom.Spawn(s, argT, dataT)")
		g.P("}")
	}
	//g.P("elem.Config.Messages = map[string]*", atomosPackage.Ident("AtomMessageConfig"), "{")
	//for _, method := range service.Methods {
	//	if method.GoName == "Spawn" || method.GoName == "SpawnWormhole" {
	//		continue
	//	}
	//	g.P("\"", method.GoName, "\": ", atomosPackage.Ident("NewAtomCallConfig"), "(&", method.Input.GoIdent, "{}, &", method.Output.GoIdent, "{}),")
	//}
	//g.P("}")
	//g.P("elem.AtomMessages = map[string]*", atomosPackage.Ident("ElementAtomMessage"), "{")
	//for _, method := range service.Methods {
	//	if method.GoName == "Spawn" || method.GoName == "SpawnWormhole" {
	//		continue
	//	}
	//	g.P("\"", method.GoName, "\": {")
	//	g.P("InDec: func(b []byte) (", protobufPackage.Ident("Message"), ", error) { return ", atomosPackage.Ident("MessageUnmarshal"), "(b, &", method.Input.GoIdent, "{}) },")
	//	g.P("OutDec: func(b []byte) (", protobufPackage.Ident("Message"), ", error) { return ", atomosPackage.Ident("MessageUnmarshal"), "(b, &", method.Output.GoIdent, "{}) },")
	//	g.P("},")
	//}
	//g.P("}")
	g.P("return elem")
	g.P("}")

	g.P()

	g.P("func Get", elementName, "Implement(dev ", atomosPackage.Ident("ElementDeveloper"), ") *", atomosPackage.Ident("ElementImplementation"), "{")
	g.P("elem := ", atomosPackage.Ident("NewImplementationFromDeveloper"), "(dev)")
	g.P("elem.Interface = Get", interfaceName, "(dev)")
	g.P("elem.ElementHandlers = map[string]", atomosPackage.Ident("MessageHandler"), "{")
	for _, method := range service.Methods {
		methodName := method.GoName
		if !strings.HasPrefix(methodName, "Element") {
			continue
		}
		if strings.HasPrefix(methodName, "ElementSpawn") {
			continue
		}
		if strings.HasPrefix(methodName, "Spawn") {
			continue
		}
		methodName = strings.Split(methodName, "Element")[1]
		if methodName == "" {
			continue
		}
		g.P("\"", methodName, "\": func(from ", atomosPackage.Ident("ID"), ", to ", atomosPackage.Ident("Atomos"), ", in ", protobufPackage.Ident("Message"), ") (", protobufPackage.Ident("Message"), ", *", atomosPackage.Ident("ErrorInfo"), ") {")
		g.P("req, ok := in.(*", method.Input.GoIdent, ")")
		g.P("if !ok { return nil, ", atomosPackage.Ident("NewErrorf"), "(", atomosPackage.Ident("ErrAtomMessageArgType"), ", \"Arg type=(%T)\", in) }")
		g.P("a, ok := to.(", elementElementName, ")")
		g.P("if !ok { return nil, ", atomosPackage.Ident("NewErrorf"), "(", atomosPackage.Ident("ErrAtomMessageAtomType"), ", \"Atom type=(%T)\", to) }")
		g.P("return a.", methodName, "(from, req)")
		g.P("},")
	}
	g.P("}")
	g.P("elem.AtomHandlers = map[string]", atomosPackage.Ident("MessageHandler"), "{")
	for _, method := range service.Methods {
		methodName := method.GoName
		if strings.HasPrefix(methodName, "Element") {
			continue
		}
		if strings.HasPrefix(methodName, "Spawn") {
			continue
		}
		if methodName == "" {
			continue
		}
		g.P("\"", methodName, "\": func(from ", atomosPackage.Ident("ID"), ", to ", atomosPackage.Ident("Atomos"), ", in ", protobufPackage.Ident("Message"), ") (", protobufPackage.Ident("Message"), ", *", atomosPackage.Ident("ErrorInfo"), ") {")
		g.P("req, ok := in.(*", method.Input.GoIdent, ")")
		g.P("if !ok { return nil, ", atomosPackage.Ident("NewErrorf"), "(", atomosPackage.Ident("ErrAtomMessageArgType"), ", \"Arg type=(%T)\", in) }")
		g.P("a, ok := to.(", elementAtomName, ")")
		g.P("if !ok { return nil, ", atomosPackage.Ident("NewErrorf"), "(", atomosPackage.Ident("ErrAtomMessageAtomType"), ", \"Atom type=(%T)\", to) }")
		g.P("return a.", methodName, "(from, req)")
		g.P("},")
	}
	g.P("}")
	g.P("return elem")
	g.P("}")
}

const deprecationComment = "// Deprecated: Do not use."

func spawnElementSign(g *protogen.GeneratedFile, method *protogen.Method) {
	g.P("Spawn(self ", atomosPackage.Ident("ElementSelfID"),
		", data *", g.QualifiedGoIdent(method.Output.GoIdent),
		") *", atomosPackage.Ident("ErrorInfo"))
}

func spawnElementDefaultSign(g *protogen.GeneratedFile) {
	g.P("Spawn(self ", atomosPackage.Ident("ElementSelfID"),
		", data *", atomosPackage.Ident("Nil"),
		") *", atomosPackage.Ident("ErrorInfo"))
}

func spawnAtomSign(g *protogen.GeneratedFile, method *protogen.Method) {
	g.P("Spawn(self ", atomosPackage.Ident("AtomSelfID"),
		", arg *", g.QualifiedGoIdent(method.Input.GoIdent),
		", data *", g.QualifiedGoIdent(method.Output.GoIdent),
		") *", atomosPackage.Ident("ErrorInfo"))
}

func methodSign(g *protogen.GeneratedFile, method *protogen.Method, methodName string) {
	g.P(methodName+"(from ", atomosPackage.Ident("ID"),
		", in *", g.QualifiedGoIdent(method.Input.GoIdent),
		") (*", g.QualifiedGoIdent(method.Output.GoIdent),
		", *", atomosPackage.Ident("ErrorInfo"),
		")")
}

func noExport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
