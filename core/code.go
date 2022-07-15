package core

const (
	OK = iota

	ErrFrameworkPanic

	ErrMainFnLoadCertFailed
	ErrMainFnCheckElementFailed

	ErrCosmosConfigInvalid
	ErrCosmosConfigNodeNameInvalid
	ErrCosmosConfigLogPathInvalid
	ErrCosmosCertConfigInvalid
	ErrCosmosHasAlreadyRun
	ErrCosmosIsBusy
	ErrCosmosIsClosed

	ErrAtomosPanic
	ErrAtomosIsNotRunning
	ErrAtomosTaskInvalidFn
	ErrAtomosTaskNotExists

	ErrAtomMessageHandlerNotExists
	ErrAtomMessageHandlerPanic
	ErrAtomKillHandlerPanic
	ErrAtomKillElementNoImplement
	ErrAtomKillElementNotImplementAutoDataPersistence
	ErrAtomReloadInvalid
	ErrAtomFromIDInvalid
	ErrAtomToIDInvalid
	ErrAtomCallDeadLock
	ErrAtomNotExists
	ErrAtomExists
	ErrAtomPersistenceRuntime
	ErrAtomSpawnArgInvalid

	ErrElementCannotKill
)
