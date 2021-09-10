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

//////////////////////////////////////////////
////////// Element: UserManagerAtom //////////
//////////////////////////////////////////////
//
// 用户管理器Atom
//

const UserManagerAtomName = "UserManagerAtom"

// UserManagerAtomId is the interface of UserManagerAtom atomos.

type UserManagerAtomId interface {
	go_atomos.Id

	// 注册用户
	RegisterUser(from go_atomos.Id, in *RegisterUserReq) (*RegisterUserResp, error)

	// 查找用户
	FindUser(from go_atomos.Id, in *FindUserReq) (*FindUserResp, error)
}

func GetUserManagerAtomId(c go_atomos.CosmosNode, name string) (UserManagerAtomId, error) {
	ca, err := c.GetAtomId(UserManagerAtomName, name)
	if err != nil {
		return nil, err
	}
	if c, ok := ca.(UserManagerAtomId); ok {
		return c, nil
	} else {
		return nil, go_atomos.ErrAtomType
	}
}

// UserManagerAtom is the atomos implements of UserManagerAtom atomos.

type UserManagerAtom interface {
	go_atomos.Atom
	Spawn(self go_atomos.AtomSelf, arg *UserManagerSpawnArg, data *UserManager) error
	RegisterUser(from go_atomos.Id, in *RegisterUserReq) (*RegisterUserResp, error)
	FindUser(from go_atomos.Id, in *FindUserReq) (*FindUserResp, error)
}

func SpawnUserManagerAtom(c go_atomos.CosmosNode, name string, arg *UserManagerSpawnArg) (UserManagerAtomId, error) {
	_, err := c.SpawnAtom(UserManagerAtomName, name, arg)
	if err != nil {
		return nil, err
	}
	id, err := c.GetAtomId(UserManagerAtomName, name)
	if err != nil {
		return nil, err
	}
	if i, ok := id.(UserManagerAtomId); ok {
		return i, nil
	}
	return nil, go_atomos.ErrAtomType
}

///////////////////////////////////////
////////// Element: UserAtom //////////
///////////////////////////////////////
//
// 用户Atom
//

const UserAtomName = "UserAtom"

// UserAtomId is the interface of UserAtom atomos.

type UserAtomId interface {
	go_atomos.Id

	// 用户信息
	UserInfo(from go_atomos.Id, in *UserInfoReq) (*UserInfoResp, error)

	// 所有好友
	GetFriends(from go_atomos.Id, in *GetFriendsReq) (*GetFriendsResp, error)

	// 添加好友
	AddFriend(from go_atomos.Id, in *AddFriendReq) (*AddFriendResp, error)

	// 房间信息
	// 房间消息推送
	RoomMessage(from go_atomos.Id, in *RoomMessagePush) (*RoomMessagePushResp, error)
}

func GetUserAtomId(c go_atomos.CosmosNode, name string) (UserAtomId, error) {
	ca, err := c.GetAtomId(UserAtomName, name)
	if err != nil {
		return nil, err
	}
	if c, ok := ca.(UserAtomId); ok {
		return c, nil
	} else {
		return nil, go_atomos.ErrAtomType
	}
}

// UserAtom is the atomos implements of UserAtom atomos.

type UserAtom interface {
	go_atomos.Atom
	Spawn(self go_atomos.AtomSelf, arg *UserSpawnArg, data *User) error
	UserInfo(from go_atomos.Id, in *UserInfoReq) (*UserInfoResp, error)
	GetFriends(from go_atomos.Id, in *GetFriendsReq) (*GetFriendsResp, error)
	AddFriend(from go_atomos.Id, in *AddFriendReq) (*AddFriendResp, error)
	RoomMessage(from go_atomos.Id, in *RoomMessagePush) (*RoomMessagePushResp, error)
}

func SpawnUserAtom(c go_atomos.CosmosNode, name string, arg *UserSpawnArg) (UserAtomId, error) {
	_, err := c.SpawnAtom(UserAtomName, name, arg)
	if err != nil {
		return nil, err
	}
	id, err := c.GetAtomId(UserAtomName, name)
	if err != nil {
		return nil, err
	}
	if i, ok := id.(UserAtomId); ok {
		return i, nil
	}
	return nil, go_atomos.ErrAtomType
}

//////////////////////////////////////////////////
////////// Element: ChatRoomManagerAtom //////////
//////////////////////////////////////////////////
//
// 房间管理器Atom
//

const ChatRoomManagerAtomName = "ChatRoomManagerAtom"

// ChatRoomManagerAtomId is the interface of ChatRoomManagerAtom atomos.

type ChatRoomManagerAtomId interface {
	go_atomos.Id

	// 创建房间
	CreateRoom(from go_atomos.Id, in *CreateRoomReq) (*CreateRoomResp, error)

	// 查找房间
	FindRoom(from go_atomos.Id, in *FindRoomReq) (*FindRoomResp, error)
}

func GetChatRoomManagerAtomId(c go_atomos.CosmosNode, name string) (ChatRoomManagerAtomId, error) {
	ca, err := c.GetAtomId(ChatRoomManagerAtomName, name)
	if err != nil {
		return nil, err
	}
	if c, ok := ca.(ChatRoomManagerAtomId); ok {
		return c, nil
	} else {
		return nil, go_atomos.ErrAtomType
	}
}

// ChatRoomManagerAtom is the atomos implements of ChatRoomManagerAtom atomos.

type ChatRoomManagerAtom interface {
	go_atomos.Atom
	Spawn(self go_atomos.AtomSelf, arg *ChatRoomManagerSpawnArg, data *ChatRoomManager) error
	CreateRoom(from go_atomos.Id, in *CreateRoomReq) (*CreateRoomResp, error)
	FindRoom(from go_atomos.Id, in *FindRoomReq) (*FindRoomResp, error)
}

func SpawnChatRoomManagerAtom(c go_atomos.CosmosNode, name string, arg *ChatRoomManagerSpawnArg) (ChatRoomManagerAtomId, error) {
	_, err := c.SpawnAtom(ChatRoomManagerAtomName, name, arg)
	if err != nil {
		return nil, err
	}
	id, err := c.GetAtomId(ChatRoomManagerAtomName, name)
	if err != nil {
		return nil, err
	}
	if i, ok := id.(ChatRoomManagerAtomId); ok {
		return i, nil
	}
	return nil, go_atomos.ErrAtomType
}

///////////////////////////////////////////
////////// Element: ChatRoomAtom //////////
///////////////////////////////////////////
//
// 房间Atom
//

const ChatRoomAtomName = "ChatRoomAtom"

// ChatRoomAtomId is the interface of ChatRoomAtom atomos.

type ChatRoomAtomId interface {
	go_atomos.Id

	// 房间信息
	Info(from go_atomos.Id, in *ChatRoomInfoReq) (*ChatRoomInfoResp, error)

	// 添加房间成员
	AddMember(from go_atomos.Id, in *AddMemberReq) (*AddMemberResp, error)

	// 删除房间成员
	DelMember(from go_atomos.Id, in *DelMemberReq) (*DelMemberResp, error)

	// 发送消息
	SendMessage(from go_atomos.Id, in *SendMessageReq) (*SendMessageResp, error)
}

func GetChatRoomAtomId(c go_atomos.CosmosNode, name string) (ChatRoomAtomId, error) {
	ca, err := c.GetAtomId(ChatRoomAtomName, name)
	if err != nil {
		return nil, err
	}
	if c, ok := ca.(ChatRoomAtomId); ok {
		return c, nil
	} else {
		return nil, go_atomos.ErrAtomType
	}
}

// ChatRoomAtom is the atomos implements of ChatRoomAtom atomos.

type ChatRoomAtom interface {
	go_atomos.Atom
	Spawn(self go_atomos.AtomSelf, arg *ChatRoomSpawnArg, data *ChatRoom) error
	Info(from go_atomos.Id, in *ChatRoomInfoReq) (*ChatRoomInfoResp, error)
	AddMember(from go_atomos.Id, in *AddMemberReq) (*AddMemberResp, error)
	DelMember(from go_atomos.Id, in *DelMemberReq) (*DelMemberResp, error)
	SendMessage(from go_atomos.Id, in *SendMessageReq) (*SendMessageResp, error)
}

func SpawnChatRoomAtom(c go_atomos.CosmosNode, name string, arg *ChatRoomSpawnArg) (ChatRoomAtomId, error) {
	_, err := c.SpawnAtom(ChatRoomAtomName, name, arg)
	if err != nil {
		return nil, err
	}
	id, err := c.GetAtomId(ChatRoomAtomName, name)
	if err != nil {
		return nil, err
	}
	if i, ok := id.(ChatRoomAtomId); ok {
		return i, nil
	}
	return nil, go_atomos.ErrAtomType
}

///////////////////////////////////////
////////// Element: KvDbAtom //////////
///////////////////////////////////////
//
// KvDb Atom
//

const KvDbAtomName = "KvDbAtom"

// KvDbAtomId is the interface of KvDbAtom atomos.

type KvDbAtomId interface {
	go_atomos.Id

	Get(from go_atomos.Id, in *DbGetReq) (*DbGetResp, error)

	Set(from go_atomos.Id, in *DbSetReq) (*DbSetResp, error)

	Del(from go_atomos.Id, in *DbDelReq) (*DbDelResp, error)
}

func GetKvDbAtomId(c go_atomos.CosmosNode, name string) (KvDbAtomId, error) {
	ca, err := c.GetAtomId(KvDbAtomName, name)
	if err != nil {
		return nil, err
	}
	if c, ok := ca.(KvDbAtomId); ok {
		return c, nil
	} else {
		return nil, go_atomos.ErrAtomType
	}
}

// KvDbAtom is the atomos implements of KvDbAtom atomos.

type KvDbAtom interface {
	go_atomos.Atom
	Spawn(self go_atomos.AtomSelf, arg *KvDbSpawnArg, data *KvDb) error
	Get(from go_atomos.Id, in *DbGetReq) (*DbGetResp, error)
	Set(from go_atomos.Id, in *DbSetReq) (*DbSetResp, error)
	Del(from go_atomos.Id, in *DbDelReq) (*DbDelResp, error)
}

func SpawnKvDbAtom(c go_atomos.CosmosNode, name string, arg *KvDbSpawnArg) (KvDbAtomId, error) {
	_, err := c.SpawnAtom(KvDbAtomName, name, arg)
	if err != nil {
		return nil, err
	}
	id, err := c.GetAtomId(KvDbAtomName, name)
	if err != nil {
		return nil, err
	}
	if i, ok := id.(KvDbAtomId); ok {
		return i, nil
	}
	return nil, go_atomos.ErrAtomType
}

//////
//// IMPLEMENTATIONS
//

//////////////////////////////////////////////
////////// Element: UserManagerAtom //////////
//////////////////////////////////////////////
//
// 用户管理器Atom
//

type userManagerAtomId struct {
	go_atomos.Id
}

func (c *userManagerAtomId) RegisterUser(from go_atomos.Id, in *RegisterUserReq) (*RegisterUserResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "RegisterUser", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*RegisterUserResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *userManagerAtomId) FindUser(from go_atomos.Id, in *FindUserReq) (*FindUserResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "FindUser", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*FindUserResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func GetUserManagerAtomInterface(dev go_atomos.ElementDeveloper) *go_atomos.ElementInterface {
	elem := go_atomos.NewInterfaceFromDeveloper(UserManagerAtomName, dev)
	elem.AtomIdConstructor = func(id go_atomos.Id) go_atomos.Id { return &userManagerAtomId{id} }
	elem.AtomSpawner = func(s go_atomos.AtomSelf, a go_atomos.Atom, arg, data proto.Message) error {
		argT, _ := arg.(*UserManagerSpawnArg)
		dataT, _ := data.(*UserManager)
		return a.(UserManagerAtom).Spawn(s, argT, dataT)
	}
	elem.Config.Messages = map[string]*go_atomos.AtomMessageConfig{
		"RegisterUser": go_atomos.NewAtomCallConfig(&RegisterUserReq{}, &RegisterUserResp{}),
		"FindUser":     go_atomos.NewAtomCallConfig(&FindUserReq{}, &FindUserResp{}),
	}
	elem.AtomMessages = map[string]*go_atomos.ElementAtomMessage{
		"RegisterUser": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &RegisterUserReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &RegisterUserResp{}) },
		},
		"FindUser": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &FindUserReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &FindUserResp{}) },
		},
	}
	return elem
}

func GetUserManagerAtomImplement(dev go_atomos.ElementDeveloper) *go_atomos.ElementImplementation {
	elem := go_atomos.NewImplementationFromDeveloper(dev)
	elem.Interface = GetUserManagerAtomInterface(dev)
	elem.AtomHandlers = map[string]go_atomos.MessageHandler{
		"RegisterUser": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*RegisterUserReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(UserManagerAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.RegisterUser(from, req)
		},
		"FindUser": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*FindUserReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(UserManagerAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.FindUser(from, req)
		},
	}
	return elem
}

///////////////////////////////////////
////////// Element: UserAtom //////////
///////////////////////////////////////
//
// 用户Atom
//

type userAtomId struct {
	go_atomos.Id
}

func (c *userAtomId) UserInfo(from go_atomos.Id, in *UserInfoReq) (*UserInfoResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "UserInfo", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*UserInfoResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *userAtomId) GetFriends(from go_atomos.Id, in *GetFriendsReq) (*GetFriendsResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "GetFriends", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*GetFriendsResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *userAtomId) AddFriend(from go_atomos.Id, in *AddFriendReq) (*AddFriendResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "AddFriend", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*AddFriendResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *userAtomId) RoomMessage(from go_atomos.Id, in *RoomMessagePush) (*RoomMessagePushResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "RoomMessage", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*RoomMessagePushResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func GetUserAtomInterface(dev go_atomos.ElementDeveloper) *go_atomos.ElementInterface {
	elem := go_atomos.NewInterfaceFromDeveloper(UserAtomName, dev)
	elem.AtomIdConstructor = func(id go_atomos.Id) go_atomos.Id { return &userAtomId{id} }
	elem.AtomSpawner = func(s go_atomos.AtomSelf, a go_atomos.Atom, arg, data proto.Message) error {
		argT, _ := arg.(*UserSpawnArg)
		dataT, _ := data.(*User)
		return a.(UserAtom).Spawn(s, argT, dataT)
	}
	elem.Config.Messages = map[string]*go_atomos.AtomMessageConfig{
		"UserInfo":    go_atomos.NewAtomCallConfig(&UserInfoReq{}, &UserInfoResp{}),
		"GetFriends":  go_atomos.NewAtomCallConfig(&GetFriendsReq{}, &GetFriendsResp{}),
		"AddFriend":   go_atomos.NewAtomCallConfig(&AddFriendReq{}, &AddFriendResp{}),
		"RoomMessage": go_atomos.NewAtomCallConfig(&RoomMessagePush{}, &RoomMessagePushResp{}),
	}
	elem.AtomMessages = map[string]*go_atomos.ElementAtomMessage{
		"UserInfo": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &UserInfoReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &UserInfoResp{}) },
		},
		"GetFriends": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &GetFriendsReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &GetFriendsResp{}) },
		},
		"AddFriend": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &AddFriendReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &AddFriendResp{}) },
		},
		"RoomMessage": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &RoomMessagePush{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &RoomMessagePushResp{}) },
		},
	}
	return elem
}

func GetUserAtomImplement(dev go_atomos.ElementDeveloper) *go_atomos.ElementImplementation {
	elem := go_atomos.NewImplementationFromDeveloper(dev)
	elem.Interface = GetUserAtomInterface(dev)
	elem.AtomHandlers = map[string]go_atomos.MessageHandler{
		"UserInfo": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*UserInfoReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(UserAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.UserInfo(from, req)
		},
		"GetFriends": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*GetFriendsReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(UserAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.GetFriends(from, req)
		},
		"AddFriend": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*AddFriendReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(UserAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.AddFriend(from, req)
		},
		"RoomMessage": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*RoomMessagePush)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(UserAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.RoomMessage(from, req)
		},
	}
	return elem
}

//////////////////////////////////////////////////
////////// Element: ChatRoomManagerAtom //////////
//////////////////////////////////////////////////
//
// 房间管理器Atom
//

type chatRoomManagerAtomId struct {
	go_atomos.Id
}

func (c *chatRoomManagerAtomId) CreateRoom(from go_atomos.Id, in *CreateRoomReq) (*CreateRoomResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "CreateRoom", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*CreateRoomResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *chatRoomManagerAtomId) FindRoom(from go_atomos.Id, in *FindRoomReq) (*FindRoomResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "FindRoom", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*FindRoomResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func GetChatRoomManagerAtomInterface(dev go_atomos.ElementDeveloper) *go_atomos.ElementInterface {
	elem := go_atomos.NewInterfaceFromDeveloper(ChatRoomManagerAtomName, dev)
	elem.AtomIdConstructor = func(id go_atomos.Id) go_atomos.Id { return &chatRoomManagerAtomId{id} }
	elem.AtomSpawner = func(s go_atomos.AtomSelf, a go_atomos.Atom, arg, data proto.Message) error {
		argT, _ := arg.(*ChatRoomManagerSpawnArg)
		dataT, _ := data.(*ChatRoomManager)
		return a.(ChatRoomManagerAtom).Spawn(s, argT, dataT)
	}
	elem.Config.Messages = map[string]*go_atomos.AtomMessageConfig{
		"CreateRoom": go_atomos.NewAtomCallConfig(&CreateRoomReq{}, &CreateRoomResp{}),
		"FindRoom":   go_atomos.NewAtomCallConfig(&FindRoomReq{}, &FindRoomResp{}),
	}
	elem.AtomMessages = map[string]*go_atomos.ElementAtomMessage{
		"CreateRoom": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &CreateRoomReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &CreateRoomResp{}) },
		},
		"FindRoom": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &FindRoomReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &FindRoomResp{}) },
		},
	}
	return elem
}

func GetChatRoomManagerAtomImplement(dev go_atomos.ElementDeveloper) *go_atomos.ElementImplementation {
	elem := go_atomos.NewImplementationFromDeveloper(dev)
	elem.Interface = GetChatRoomManagerAtomInterface(dev)
	elem.AtomHandlers = map[string]go_atomos.MessageHandler{
		"CreateRoom": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*CreateRoomReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(ChatRoomManagerAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.CreateRoom(from, req)
		},
		"FindRoom": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*FindRoomReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(ChatRoomManagerAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.FindRoom(from, req)
		},
	}
	return elem
}

///////////////////////////////////////////
////////// Element: ChatRoomAtom //////////
///////////////////////////////////////////
//
// 房间Atom
//

type chatRoomAtomId struct {
	go_atomos.Id
}

func (c *chatRoomAtomId) Info(from go_atomos.Id, in *ChatRoomInfoReq) (*ChatRoomInfoResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "Info", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*ChatRoomInfoResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *chatRoomAtomId) AddMember(from go_atomos.Id, in *AddMemberReq) (*AddMemberResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "AddMember", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*AddMemberResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *chatRoomAtomId) DelMember(from go_atomos.Id, in *DelMemberReq) (*DelMemberResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "DelMember", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*DelMemberResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *chatRoomAtomId) SendMessage(from go_atomos.Id, in *SendMessageReq) (*SendMessageResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "SendMessage", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*SendMessageResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func GetChatRoomAtomInterface(dev go_atomos.ElementDeveloper) *go_atomos.ElementInterface {
	elem := go_atomos.NewInterfaceFromDeveloper(ChatRoomAtomName, dev)
	elem.AtomIdConstructor = func(id go_atomos.Id) go_atomos.Id { return &chatRoomAtomId{id} }
	elem.AtomSpawner = func(s go_atomos.AtomSelf, a go_atomos.Atom, arg, data proto.Message) error {
		argT, _ := arg.(*ChatRoomSpawnArg)
		dataT, _ := data.(*ChatRoom)
		return a.(ChatRoomAtom).Spawn(s, argT, dataT)
	}
	elem.Config.Messages = map[string]*go_atomos.AtomMessageConfig{
		"Info":        go_atomos.NewAtomCallConfig(&ChatRoomInfoReq{}, &ChatRoomInfoResp{}),
		"AddMember":   go_atomos.NewAtomCallConfig(&AddMemberReq{}, &AddMemberResp{}),
		"DelMember":   go_atomos.NewAtomCallConfig(&DelMemberReq{}, &DelMemberResp{}),
		"SendMessage": go_atomos.NewAtomCallConfig(&SendMessageReq{}, &SendMessageResp{}),
	}
	elem.AtomMessages = map[string]*go_atomos.ElementAtomMessage{
		"Info": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &ChatRoomInfoReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &ChatRoomInfoResp{}) },
		},
		"AddMember": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &AddMemberReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &AddMemberResp{}) },
		},
		"DelMember": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &DelMemberReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &DelMemberResp{}) },
		},
		"SendMessage": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &SendMessageReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &SendMessageResp{}) },
		},
	}
	return elem
}

func GetChatRoomAtomImplement(dev go_atomos.ElementDeveloper) *go_atomos.ElementImplementation {
	elem := go_atomos.NewImplementationFromDeveloper(dev)
	elem.Interface = GetChatRoomAtomInterface(dev)
	elem.AtomHandlers = map[string]go_atomos.MessageHandler{
		"Info": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*ChatRoomInfoReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(ChatRoomAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.Info(from, req)
		},
		"AddMember": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*AddMemberReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(ChatRoomAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.AddMember(from, req)
		},
		"DelMember": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*DelMemberReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(ChatRoomAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.DelMember(from, req)
		},
		"SendMessage": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*SendMessageReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(ChatRoomAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.SendMessage(from, req)
		},
	}
	return elem
}

///////////////////////////////////////
////////// Element: KvDbAtom //////////
///////////////////////////////////////
//
// KvDb Atom
//

type kvDbAtomId struct {
	go_atomos.Id
}

func (c *kvDbAtomId) Get(from go_atomos.Id, in *DbGetReq) (*DbGetResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "Get", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*DbGetResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *kvDbAtomId) Set(from go_atomos.Id, in *DbSetReq) (*DbSetResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "Set", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*DbSetResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func (c *kvDbAtomId) Del(from go_atomos.Id, in *DbDelReq) (*DbDelResp, error) {
	r, err := c.Cosmos().MessageAtom(from, c, "Del", in)
	if r == nil {
		return nil, err
	}
	reply, ok := r.(*DbDelResp)
	if !ok {
		return nil, go_atomos.ErrAtomMessageReplyType
	}
	return reply, nil
}

func GetKvDbAtomInterface(dev go_atomos.ElementDeveloper) *go_atomos.ElementInterface {
	elem := go_atomos.NewInterfaceFromDeveloper(KvDbAtomName, dev)
	elem.AtomIdConstructor = func(id go_atomos.Id) go_atomos.Id { return &kvDbAtomId{id} }
	elem.AtomSpawner = func(s go_atomos.AtomSelf, a go_atomos.Atom, arg, data proto.Message) error {
		argT, _ := arg.(*KvDbSpawnArg)
		dataT, _ := data.(*KvDb)
		return a.(KvDbAtom).Spawn(s, argT, dataT)
	}
	elem.Config.Messages = map[string]*go_atomos.AtomMessageConfig{
		"Get": go_atomos.NewAtomCallConfig(&DbGetReq{}, &DbGetResp{}),
		"Set": go_atomos.NewAtomCallConfig(&DbSetReq{}, &DbSetResp{}),
		"Del": go_atomos.NewAtomCallConfig(&DbDelReq{}, &DbDelResp{}),
	}
	elem.AtomMessages = map[string]*go_atomos.ElementAtomMessage{
		"Get": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &DbGetReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &DbGetResp{}) },
		},
		"Set": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &DbSetReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &DbSetResp{}) },
		},
		"Del": {
			InDec:  func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &DbDelReq{}) },
			OutDec: func(b []byte) (proto.Message, error) { return go_atomos.MessageUnmarshal(b, &DbDelResp{}) },
		},
	}
	return elem
}

func GetKvDbAtomImplement(dev go_atomos.ElementDeveloper) *go_atomos.ElementImplementation {
	elem := go_atomos.NewImplementationFromDeveloper(dev)
	elem.Interface = GetKvDbAtomInterface(dev)
	elem.AtomHandlers = map[string]go_atomos.MessageHandler{
		"Get": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*DbGetReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(KvDbAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.Get(from, req)
		},
		"Set": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*DbSetReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(KvDbAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.Set(from, req)
		},
		"Del": func(from go_atomos.Id, to go_atomos.Atom, in proto.Message) (proto.Message, error) {
			req, ok := in.(*DbDelReq)
			if !ok {
				return nil, go_atomos.ErrAtomMessageArgType
			}
			a, ok := to.(KvDbAtom)
			if !ok {
				return nil, go_atomos.ErrAtomMessageAtomType
			}
			return a.Del(from, req)
		},
	}
	return elem
}
