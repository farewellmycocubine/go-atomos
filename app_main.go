package go_atomos

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var app *App
var nodeOrSupervisor = true
var supervisorCommandHandlers map[string]AppUDSCommandFn

func SharedApp() *App {
	return app
}

func SupervisorSet(handlers map[string]AppUDSCommandFn) {
	nodeOrSupervisor = false
	supervisorCommandHandlers = handlers
}

func Main(runnable CosmosRunnable) {
	if nodeOrSupervisor {
		log.Printf("Welcome to Atomos! pid=(%d)", os.Getpid())
	} else {
		log.Printf("Welcome to Atomos Supervisor! pid=(%d)", os.Getpid())
	}

	var (
		configPath = flag.String("config", "", "config path")
		standalone = flag.Bool("standalone", false, "standalone")
	)
	flag.Parse()

	var err *Error

	// Load Config.
	if configPath == nil {
		log.Println("App: No config path specified.", os.Getpid())
		os.Exit(1)
	}

	// Check.
	if nodeOrSupervisor {
		app, err = NewCosmosNodeApp(*configPath)
	} else {
		app, err = NewCosmosSupervisorApp(*configPath)
	}
	if err != nil {
		log.Printf("App: Config is invalid. pid=(%d),err=(%v)", os.Getpid(), err)
		os.Exit(1)
	}

	// Init.
	InitCosmosProcess(app.logging.WriteAccessLog, app.logging.WriteErrorLog)

	isRunning, processID, err := app.Check()
	if err != nil && !isRunning {
		msg := fmt.Sprintf("App: Check failed. err=(%v)", err)
		SharedLogging().PushProcessLog(LogLevel_Fatal, msg)
		log.Printf(msg)
		os.Exit(1)
	}
	if isRunning {
		msg := fmt.Sprintf("App: App is already running. pid=(%d)", processID)
		SharedLogging().PushProcessLog(LogLevel_Fatal, msg)
		log.Printf(msg)
		os.Exit(1)
	}

	sa := false
	if standalone != nil {
		sa = *standalone
	}
	if IsParentProcess() && !sa {
		if err = app.ForkAppProcess(); err != nil {
			msg := fmt.Sprintf("App: Fork app failed. err=(%v)", err)
			SharedLogging().PushProcessLog(LogLevel_Fatal, msg)
			log.Printf(msg)
			os.Exit(1)
		}
		msg := fmt.Sprintf("App: Fork app succeed. Loader will exit.")
		SharedLogging().PushProcessLog(LogLevel_Info, msg)
		log.Printf(msg)
		log.Printf("App: Access Log File=(%s)", app.logging.curAccessLogName)
		log.Printf("App: Error Log File=(%s)", app.logging.curErrorLogName)
		app.logging.Close()
		return
	} else {
		if err = app.LaunchApp(); err != nil {
			msg := fmt.Sprintf("App: Launch app failed. err=(%v)", err)
			SharedLogging().PushProcessLog(LogLevel_Fatal, msg)
			log.Printf(msg)
			os.Exit(1)
		}
		defer func() {
			SharedLogging().PushProcessLog(LogLevel_Info, "App: Exiting.")
			app.close()
		}()
		runnable.SetConfig(app.config)
		if err = SharedCosmosProcess().Start(&runnable); err != nil {
			SharedLogging().PushProcessLog(LogLevel_Err, "App: Runnable starts failed. err=(%v)", err.AddStack(nil))
			return
		}
		SharedLogging().PushProcessLog(LogLevel_Info, "App: Started.")
		<-app.WaitExitApp()
		if err = SharedCosmosProcess().Stop(); err != nil {
			SharedLogging().PushProcessLog(LogLevel_Err, "App: Runnable stops with error. err=(%v)", err.AddStack(nil))
		}
		return
	}
}
