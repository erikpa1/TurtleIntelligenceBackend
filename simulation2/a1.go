package simulation2

import "turtle/simulation2/behaviours"

func InitSimFunctions() bool {
	behaviours.InitBehBuffer()
	behaviours.InitBehProcess()
	behaviours.InitBehDelay()
	behaviours.InitBehSpawn()
	behaviours.InitBehSink()

	//Movable antities
	behaviours.InitBehWorkerPool()

	//Statistiscs
	behaviours.InitEntryStatistics()

	return true
}

var SimInit = InitSimFunctions()
