package simulation2

import "turtle/core/lgr"

var SPAWN_FUNCTIONS = SimFunctions{}

func InitSpawnBehaviour() {
	var _step FnStep = _SpawnStep
	SPAWN_FUNCTIONS[FN_STEP] = _step
}

func _SpawnStep(self *SimEntity) {
	lgr.Error("Spawn step")

}

func NewSpawnBehaviour(entity *SimEntity) {
	entity.Functions = BUFFER_FUNCTIONS
}
