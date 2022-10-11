package go_atomos

const (
	OK = iota

	ErrFrameworkPanic

	ErrProcessRunnableInvalid

	ErrMainLoadCertFailed
	ErrMainCheckElementFailed
	ErrMainRunnableNotFound
	ErrMainIsReloading
	ErrMainElementNotFound
	ErrMainElementIsInvalid
	ErrMainReloadFailed
	ErrMainStartRunningPanic
	ErrMainCannotKill
	ErrMainCannotMessage
	ErrMainCannotScale

	ErrCosmosConfigInvalid
	ErrCosmosConfigNotFound
	ErrCosmosNameInvalid
	ErrCosmosNodeNameInvalid
	ErrCosmosNodeNameNotFound
	ErrCosmosConfigLogPathInvalid
	ErrCosmosConfigRunPathInvalid
	ErrCosmosConfigRunPIDIsRunning
	ErrCosmosConfigRunPIDPathInvalid
	ErrCosmosConfigRunPIDInvalid
	ErrCosmosConfigRemovePIDPathFailed
	ErrCosmosConfigBuildPathInvalid
	ErrCosmosConfigBinPathInvalid
	ErrCosmosConfigBuildFailed
	ErrCosmosConfigBinLinkFileFailed
	ErrCosmosReadEtcPath
	ErrCosmosCertConfigInvalid
	ErrCosmosLogOpenFailed
	ErrCosmosWritePIDFileFailed
	ErrCosmosWriteUnixSocketFailed
	ErrCosmosHasAlreadyRun
	ErrCosmosIsBusy
	ErrCosmosIsClosed
	ErrCosmosUnixSocketConnEOF
	ErrCosmosUnixSocketCommandNotSupported
	ErrCosmosUnixSocketConnError
	ErrCosmosUnixSocketConnWriteError
	ErrCosmosDaemonGetExecutableFailed
	ErrCosmosDaemonStartProcessFailed
	ErrCosmosProcessIDFileNotFound
	ErrCosmosNodeRunPathWatchFailed

	ErrPathUserInvalid
	ErrPathGroupInvalid
	ErrPathNotExist
	ErrPathMakeDir
	ErrPathGroupIDInvalid
	ErrPathUserIDInvalid
	ErrPathChangeOwnFailed
	ErrPathChangeModeFailed
	ErrPathIsNotDirectory
	ErrPathIsNotOwner
	ErrPathPermNotMatch
	ErrPathSaveFileFailed
	ErrPathGetGroupIDsFailed

	ErrElementLoaded
	ErrElementReloadInvalid
	ErrElementMessageHandlerNotExists
	ErrElementScaleHandlerNotExists
	ErrElementMessageHandlerPanic
	ErrElementKillHandlerPanic
	ErrElementSpawnArgInvalid
	ErrElementNotImplemented

	ErrAtomosPanic
	ErrAtomosIsStopping
	ErrAtomosIsNotRunning
	ErrAtomosTaskInvalidFn
	ErrAtomosTaskNotExists
	ErrAtomosNotSupportWormhole

	ErrAtomMessageHandlerNotExists
	ErrAtomMessageHandlerPanic
	ErrAtomKillHandlerPanic
	ErrAtomKillElementNoImplement
	ErrAtomKillElementNotImplementAutoDataPersistence
	ErrAtomReloadInvalid
	ErrAtomFromIDInvalid
	ErrAtomToIDInvalid
	ErrAtomSpawnArgInvalid
	ErrAtomCallDeadLock
	ErrAtomNotExists
	ErrAtomExists
	ErrAtomPersistenceRuntime
	ErrAtomCannotScale

	ErrAtomNotImplemented
	ErrAtomMessageAtomType
	ErrAtomMessageArgType
	ErrAtomMessageReplyType

	ErrElementCannotKill
)
