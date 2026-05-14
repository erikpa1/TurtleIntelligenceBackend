package simulation2

func InitSimFunctions() bool {
	InitBehBuffer()
	InitBehProcess()
	InitSpawnBehaviour()

	return true
}

var SimInit = InitSimFunctions()
