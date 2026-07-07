package simulation2

import (
	"turtle/simulation2/behaviours/behBuffer"
	"turtle/simulation2/behaviours/behDelay"
	"turtle/simulation2/behaviours/behLogisticsControl"
	"turtle/simulation2/behaviours/behProcess"
	"turtle/simulation2/behaviours/behSink"
	"turtle/simulation2/behaviours/behSpawn"
	"turtle/simulation2/behaviours/behWorkerPool"
	"turtle/simulation2/behaviours/entryStats"
)

func InitSimFunctions() bool {
	behBuffer.InitBehBuffer()
	behProcess.InitBehProcess()
	behDelay.InitBehDelay()
	behSpawn.InitBehSpawn()
	behSink.InitBehSink()

	//Movable antities
	behWorkerPool.InitBehWorkerPool()

	//Statistiscs
	entryStats.InitEntryStatistics()

	//Contorls
	behLogisticsControl.InitBehWorkerPool()

	return true
}

var SimInit = InitSimFunctions()
