package simulation2

func InitSimFunctions() bool {
	InitBehBuffer()
	InitBehProcess()
	InitBehSpawn()
	InitBehSink()

	return true
}

var SimInit = InitSimFunctions()
