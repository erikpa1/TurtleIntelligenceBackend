package simulation2

func InitSimFunctions() bool {
	InitBehBuffer()
	InitBehProcess()
	InitBehDelay()
	InitBehSpawn()
	InitBehSink()

	//Movable antities
	InitBehWorkerPool()

	//Statistiscs
	InitEntryStatistics()

	return true
}

var SimInit = InitSimFunctions()
