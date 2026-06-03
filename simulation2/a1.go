package simulation2

func InitSimFunctions() bool {
	InitBehBuffer()
	InitBehProcess()
	InitBehDelay()
	InitBehSpawn()
	InitBehSink()

	//Statistiscs
	InitEntryStatistics()

	return true
}

var SimInit = InitSimFunctions()
