package simulation2

import (
	"turtle/tools"
)

var SPAWN_FUNCTIONS = SimFunctions{}

func InitSpawnBehaviour() {
	var _init1 FnInit = _SpawnInit1
	SPAWN_FUNCTIONS[FN_INIT1] = _init1

	var _step FnStep = _SpawnStep
	SPAWN_FUNCTIONS[FN_STEP] = _step

}

func NewSpawnBehaviour(entity *SimEntity) {
	entity.Functions = SPAWN_FUNCTIONS

	spawnBehaviour := make(SimBehData)
	spawnBehaviour["SpawnInterval"] = entity.TypeData.GetSeconds("spawn_interval", 1)
	spawnBehaviour["SpawnLimit"] = entity.TypeData.GetInt64("spawn_limit", 1)
	spawnBehaviour["SpawnActor"] = entity.TypeData.GetPrimitiveObjectId("actor")
	spawnBehaviour[spwn_NextSpawnTime] = tools.MaxSeconds()

	spawnBehaviour["Actor"] = nil

	entity.BehaviourEntity = spawnBehaviour

}
