package simulation2

import "turtle/tools"

var SPAWN_FUNCTIONS = SimFunctions{}

func InitBehSpawn() {
	var _init1 FnInit = _SpawnInit1
	SPAWN_FUNCTIONS[FN_INIT1] = _init1

	var _step FnStep = _SpawnStep
	SPAWN_FUNCTIONS[FN_STEP] = _step
}

func NewSpawnBehaviour(entity *SimEntity) {
	spawn := &BehSpawn{}
	spawn.Entity = entity
	spawn.World = entity.World

	spawn.SpawnInterval = entity.TypeData.GetSeconds("spawn_interval", 1)
	spawn.SpawnLimit = entity.TypeData.GetInt("spawn_limit", 1)

	spawn.SpawnMultiplication = entity.TypeData.GetInt("spawn_multiplication", 1)

	if spawn.SpawnMultiplication <= 0 {
		spawn.SpawnMultiplication = 1
	}

	spawn.SpawnActorUid = entity.TypeData.GetPrimitiveObjectId("actor")
	spawn.NextSpawnTime = tools.MaxSeconds()

	entity.Impl = spawn
	entity.Functions = SPAWN_FUNCTIONS
}

func _SpawnInit1(self *SimEntity) {
	spawn := GetBehSpawn(self)
	spawn._CalculateNextSpawn()
}

func _SpawnStep(self *SimEntity) {
	spawn := GetBehSpawn(self)
	spawn.Step()
}
