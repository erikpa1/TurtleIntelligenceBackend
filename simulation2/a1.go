package simulation2

func InitSimFunctions() bool {
	InitBehBuffer()
	InitBehProcess()
	InitBehDelay()
	InitBehSpawn()
	InitBehSink()

	return true
}

var SimInit = InitSimFunctions()
