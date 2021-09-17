package main

import (
	atomos "github.com/hwangtou/go-atomos"
	"github.com/hwangtou/go-atomos/examples/chat/api"
	"github.com/hwangtou/go-atomos/examples/chat/element"
)

func main() {
	runnable := atomos.CosmosRunnable{}
	runnable.AddElementImplementation(api.GetKvDbAtomImplement(&element.KvDbElement{})).
		AddElementImplementation(api.GetUserManagerImplement(&element.UserManagerElement{})).
		AddElementImplementation(api.GetChatRoomManagerImplement(&element.ChatManagerElement{})).
		AddElementImplementation(api.GetUserImplement(&element.UserElement{})).
		AddElementImplementation(api.GetChatRoomImplement(&element.ChatRoomElement{})).
		SetScript(scriptChat)
	config := &atomos.Config{
		Node:               "Chat",
		LogPath:            "/tmp/cosmos_log/",
		LogLevel:           atomos.LogLevel_Debug,
		EnableServer: &atomos.RemoteServerConfig{
			Port:     10001,
		},
		EnableCert: &atomos.CertConfig{
			CertPath: "server.crt",
			KeyPath:  "server.key",
		},
	}
	// Cycle
	cosmos := atomos.NewCosmosCycle()
	defer cosmos.Close()
	exitCh, err := cosmos.Daemon(config)
	if err != nil {
		return
	}
	cosmos.SendRunnable(runnable)
	<-exitCh
}

func scriptChat(cosmos *atomos.CosmosSelf, mainId atomos.MainId, killNoticeChannel chan bool) {
	// Spawn KV DB
	dbId, err := api.SpawnKvDbAtom(cosmos.Local(), "DB", &api.KvDbSpawnArg{ DbPath: "data" })
	if err != nil {
		mainId.Log().Fatal("KvDb spawn failed, err=%v", err)
		return
	}
	defer func() {
		mainId.Log().Info("KvDb is exiting")
		if err = dbId.Kill(mainId); err != nil {
			mainId.Log().Error("KvDb exited with error, err=%v", err)
		}
	}()

	// Spawn UserManager
	userManagerId, err := api.SpawnUserManager(cosmos.Local(), "UserManager", &api.UserManagerSpawnArg{})
	if err != nil {
		mainId.Log().Fatal("UserManager spawn failed, err=%v", err)
		return
	}
	defer func() {
		mainId.Log().Info("UserManager is exiting")
		if err = userManagerId.Kill(mainId); err != nil {
			mainId.Log().Error("UserManager exited with error, err=%v", err)
		}
	}()

	// Spawn ChatManager
	chatManagerId, err := api.SpawnChatRoomManager(cosmos.Local(), "RoomManager", &api.ChatRoomManagerSpawnArg{})
	if err != nil {
		mainId.Log().Fatal("ChatManager spawn failed, err=%v", err)
	}
	defer func() {
		mainId.Log().Info("ChatManager is exiting")
		if err = chatManagerId.Kill(mainId); err != nil {
			mainId.Log().Error("ChatManager exited with error, err=%v", err)
		}
	}()
	<-killNoticeChannel
}
